package client

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/jiuyunyue/valsta/src/types"
	"github.com/jiuyunyue/valsta/utils"

	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	slakingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func (c *GClient) QueryUptime(start, end int64) (sig types.Uptime) {
	sig = make(map[string]types.ValidatorInfo)
	validators := c.QueryValidators()

	for from := start; from <= end; from++ {
		block, err := c.TMServiceQuery.GetBlockByHeight(context.Background(), &tmservice.GetBlockByHeightRequest{Height: from})
		if err != nil {
			panic(err)
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

			if val.Times == 0 {
				val.Times = 1
			} else {
				val.Times++
			}

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
	return sig
}

func (c *GClient) QueryJailed(start, end int64) (list types.Jailed) {
	list = make(map[string]bool)

	for from := start; from <= end; from++ {
		res, err := c.SignClient.BlockResults(context.Background(), &from)
		if err != nil {
			panic(err)
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
							panic(err)
						}
						address := strings.ToUpper(hex.EncodeToString(consAddress.Bytes()))
						list[address] = true
					}
				}
			}
		}
	}
	return list
}

func (c *GClient) QueryValidators() map[string]sdk.AccAddress {
	validatorInfos := make(map[string]sdk.AccAddress)
	Pagination := &query.PageRequest{
		Key:        []byte(""),
		Limit:      1000,
		Offset:     0,
		CountTotal: false,
		Reverse:    false,
	}
	validators, err := c.StakingQuery.Validators(
		context.Background(),
		&stakingtypes.QueryValidatorsRequest{Pagination: Pagination},
	)
	if err != nil {
		panic(err)
	}

	for _, val := range validators.Validators {
		var pk cryptotypes.PubKey
		err = cdc.UnpackAny(val.ConsensusPubkey, &pk)
		valAddress, err := sdk.ValAddressFromBech32(val.OperatorAddress)
		if err != nil {
			panic(err)
		}
		accAddress, err := sdk.AccAddressFromHex(hex.EncodeToString(valAddress.Bytes()))
		if err != nil {
			panic(err)
		}
		validatorInfos[pk.Address().String()] = accAddress
	}
	return validatorInfos
}
