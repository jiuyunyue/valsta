package client

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	typestx "github.com/cosmos/cosmos-sdk/types/tx"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/jiuyunyue/valsta/src/types"
	"github.com/jiuyunyue/valsta/utils"

	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	slakingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func (c *GClient) QueryUptime(start, end int64) (sig types.Uptime, err error) {
	sig = make(map[string]types.ValidatorInfo)
	validators, err := c.QueryValidators()
	if err != nil {
		return nil, err
	}

	for from := start; from <= end; from++ {
		block, err := c.TMServiceQuery.GetBlockByHeight(context.Background(), &tmservice.GetBlockByHeightRequest{Height: from})
		if err != nil {
			return nil, err
		}

		utils.CallClear()
		now := float64(from-start) / float64(end-start) * 100
		if now == 1 {
			fmt.Printf("QueryUptime info : now at end : height %d \n", from)
		} else {
			fmt.Printf("QueryUptime info : now at %d : end %d  : total %d : %2f%% \n", from, end, end-start+1, now)
		}
		for _, signature := range block.Block.LastCommit.Signatures {
			if len(signature.ValidatorAddress) == 0 {
				continue
			}
			key := strings.ToUpper(hex.EncodeToString(signature.ValidatorAddress))

			val := sig[key]
			val.Address = key
			val.AccAddress = validators[key].String()
			val.Times++
			sig[key] = val
		}
	}

	for k, v := range sig {
		val := v
		num := float64(v.Times) / float64(end-start+1) * 100
		str := strconv.FormatFloat(num, 'f', 2, 64)
		val.SurRate = str
		sig[k] = val
	}
	return sig, nil
}

func (c *GClient) QueryJailed(start, end int64) (list types.Jailed, err error) {
	list = make(map[string]bool)

	for from := start; from <= end; from++ {
		res, err := c.SignClient.BlockResults(context.Background(), &from)
		if err != nil {
			return nil, err
		}

		utils.CallClear()
		now := float64(from-start) / float64(end-start) * 100
		if now == 100 {
			fmt.Printf("QueryJailed info : now at end : height %d \n", from)
		} else {
			fmt.Printf("QueryJailed info : now at %d : end %d  : %v%% \n", from, end, now)
		}

		for _, event := range res.BeginBlockEvents {
			if event.Type == slakingtypes.EventTypeSlash {
				for _, attr := range event.Attributes {
					if bytes.Equal(attr.Key, []byte(slakingtypes.AttributeKeyJailed)) {
						consAddress, err := sdk.ConsAddressFromBech32(string(attr.Value))
						if err != nil {
							return nil, err
						}
						address := strings.ToUpper(hex.EncodeToString(consAddress.Bytes()))
						list[address] = true
					}
				}
			}
		}
	}
	return list, nil
}

func (c *GClient) QueryValidators() (map[string]sdk.AccAddress, error) {
	validatorInfos := make(map[string]sdk.AccAddress)
	Pagination := &query.PageRequest{
		Key:        []byte(""),
		Limit:      10000,
		Offset:     0,
		CountTotal: false,
		Reverse:    false,
	}
	validators, err := c.StakingQuery.Validators(
		context.Background(),
		&stakingtypes.QueryValidatorsRequest{Pagination: Pagination},
	)
	if err != nil {
		return nil, err
	}
	for _, val := range validators.Validators {
		var pk cryptotypes.PubKey
		err = cdc.UnpackAny(val.ConsensusPubkey, &pk)
		valAddress, err := sdk.ValAddressFromBech32(val.OperatorAddress)
		if err != nil {
			return nil, err
		}
		accAddress, err := sdk.AccAddressFromHex(hex.EncodeToString(valAddress.Bytes()))
		if err != nil {
			return nil, err
		}
		validatorInfos[pk.Address().String()] = accAddress
	}
	return validatorInfos, nil
}

func (c *GClient) QueryVoters() (voters map[string]types.VoterInfo, err error) {
	voters = make(map[string]types.VoterInfo)

	proposals, err := c.GovQuery.Proposals(context.Background(), &govtypes.QueryProposalsRequest{})
	if err != nil {
		return nil, err
	}
	for i := range proposals.Proposals {
		queryEvent := fmt.Sprintf("proposal_vote.proposal_id=%v", i)
		Pagination := &query.PageRequest{
			Key:        []byte(""),
			Limit:      20000,
			Offset:     0,
			CountTotal: false,
			Reverse:    false,
		}
		res, err := c.TxClient.GetTxsEvent(context.Background(), &typestx.GetTxsEventRequest{Events: []string{queryEvent}, Pagination: Pagination, OrderBy: 1})
		if err != nil {
			return nil, err
		}
		for _, tx := range res.Txs {
			for _, msg := range tx.Body.Messages {
				var mv sdk.Msg
				err = cdc.UnpackAny(msg, &mv)
				if err != nil {
					return nil, err
				}
				v := mv.(*govtypes.MsgVote)
				voter := voters[v.Voter]
				voter.Address = v.Voter
				voter.VoteTimes++
				voter.VoteProposals = append(voter.VoteProposals, uint64(i))

				voters[v.Voter] = voter
			}
		}
	}

	return voters, nil
}

func (c *GClient) SignTimes(addr string) (heights []int, err error) {
	start := 841500
	end := 1412246

	validators, err := c.QueryValidators()
	if err != nil {
		return heights, err
	}

	for from := start; from <= end; from++ {
		block, err := c.TMServiceQuery.GetBlockByHeight(context.Background(), &tmservice.GetBlockByHeightRequest{Height: int64(from)})
		if err != nil {
			return heights, err
		}

		utils.CallClear()
		now := float64(from-start) / float64(end-start) * 100
		if now == 1 {
			fmt.Printf("QueryUptime info : now at end : height %d \n", from)
		} else {
			fmt.Printf("QueryUptime info : now at %d : end %d  : total %d : %2f%% \n", from, end, end-start+1, now)
		}
		for _, signature := range block.Block.LastCommit.Signatures {
			if len(signature.ValidatorAddress) == 0 {
				continue
			}
			key := strings.ToUpper(hex.EncodeToString(signature.ValidatorAddress))

			if validators[key].String() == addr {
				heights = append(heights, from)
			}
		}
	}

	return heights, nil
}

func (c *GClient) SignHeight(addr string) (bool, int, error) {
	start := 841500
	end := 1412246
	height := 0
	have := false

	validators, err := c.QueryValidators()
	if err != nil {
		return have, height, err
	}

	for from := start; from <= end; from++ {
		block, err := c.TMServiceQuery.GetBlockByHeight(context.Background(), &tmservice.GetBlockByHeightRequest{Height: int64(from)})
		if err != nil {
			return have, height, err
		}

		utils.CallClear()
		now := float64(from-start) / float64(end-start) * 100
		if now == 1 {
			fmt.Printf("QueryUptime info : now at end : height %d \n", from)
		} else {
			fmt.Printf("QueryUptime info : now at %d : end %d  : total %d : %2f%% \n", from, end, end-start+1, now)
		}
		for _, signature := range block.Block.LastCommit.Signatures {
			if len(signature.ValidatorAddress) == 0 {
				continue
			}
			key := strings.ToUpper(hex.EncodeToString(signature.ValidatorAddress))
			vaaAddr := validators[key].String()
			if vaaAddr == addr {
				have = true
				height = from
				return have, height, err
			}
		}
	}

	return have, height, err
}
