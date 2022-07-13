package client

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"testing"

	abci "github.com/tendermint/tendermint/abci/types"

	typestx "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/tendermint/tendermint/crypto/tmhash"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/stretchr/testify/require"

	tmcrypto "github.com/tendermint/tendermint/crypto"

	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
	slakingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

const (
	UptickGrpcUrl = "peer0.testnet.uptick.network:9090"
	LocalGrpc     = "localhost:9090"

	UptickRpc = "http://peer0.testnet.uptick.network:26657"
	LocalRpc  = "http://localhost:26657"
)

func TestQueryBalance(t *testing.T) {
	c, err := NewGRPCClient(UptickGrpcUrl, UptickGrpcUrl)
	require.NoError(t, err)

	Pagination := &query.PageRequest{
		Key:        []byte(""),
		Limit:      1000,
		Offset:     0,
		CountTotal: false,
		Reverse:    false,
	}
	res, err := c.GovQuery.Votes(context.Background(), &govtypes.QueryVotesRequest{ProposalId: 1, Pagination: Pagination})
	require.NoError(t, err)
	fmt.Println(len(res.Votes))

	res1, err := c.BankQuery.Balance(context.Background(), &types.QueryBalanceRequest{Address: "uptick1xt42uffeg655ew0tr2ldvuy4lqhy89tmnkavmu", Denom: "auptick"})
	require.NoError(t, err)
	t.Log(res1.String())
}

func TestTmQuery(t *testing.T) {
	c, err := NewGRPCClient(UptickGrpcUrl, UptickGrpcUrl)
	require.NoError(t, err)
	for i := 841500; i <= 1412246; i += 1000 {
		valSet, err := c.TMServiceQuery.GetValidatorSetByHeight(context.Background(), &tmservice.GetValidatorSetByHeightRequest{Height: int64(i)})
		require.NoError(t, err)
		require.NotNil(t, valSet.Validators)
		block, err := c.TMServiceQuery.GetBlockByHeight(context.Background(), &tmservice.GetBlockByHeightRequest{Height: int64(i)})
		require.NoError(t, err)
		require.NotNil(t, block.Block.LastCommit.Signatures)
		for _, v := range valSet.Validators {
			if v.Address == "uptickvalcons1c8y75a5nypmhngz5dktq9mjvp9d6auz982vnsq" {
				fmt.Println("find")
			}
		}
		for _, v := range block.Block.LastCommit.Signatures {
			if len(v.Signature) == 0 {
				continue
			}
			hexStr := strings.ToUpper(hex.EncodeToString(v.ValidatorAddress))
			if hexStr == "F6E863E750CAAC822D1E388C938CDCEA2260E256" {
				fmt.Println("find")
			}
		}
		fmt.Printf("now at %v  no found \n", i)
	}
}

func TestTmQueryBlock(t *testing.T) {
	c, err := NewGRPCClient(UptickGrpcUrl, UptickGrpcUrl)
	require.NoError(t, err)
	valSet, err := c.TMServiceQuery.GetValidatorSetByHeight(context.Background(), &tmservice.GetValidatorSetByHeightRequest{Height: 300069})
	require.NoError(t, err)
	require.NotNil(t, valSet.Validators)
	block, err := c.TMServiceQuery.GetBlockByHeight(context.Background(), &tmservice.GetBlockByHeightRequest{Height: 300069})
	require.NoError(t, err)
	require.NotNil(t, block.Block.LastCommit.Signatures)

	val, err := sdk.ValAddressFromHex(hex.EncodeToString(block.Block.LastCommit.Signatures[0].ValidatorAddress))
	require.NoError(t, err)
	//val.String()
	err = sdk.VerifyAddressFormat(block.Block.LastCommit.Signatures[0].ValidatorAddress)
	require.NoError(t, err)

	var pubKey tmcrypto.PubKey
	var pk cryptotypes.PubKey
	err = cdc.UnpackAny(valSet.Validators[0].PubKey, &pk)
	require.NoError(t, err)
	pubKey, err = cryptocodec.ToTmPubKeyInterface(pk)
	require.NoError(t, err)

	require.True(t, bytes.Equal(pubKey.Bytes(), pk.Bytes()))
	require.True(t, bytes.Equal(pubKey.Address().Bytes(), val.Bytes()))
}

func TestStakingQueryBlock(t *testing.T) {
	c, err := NewGRPCClient(UptickGrpcUrl, UptickRpc)
	require.NoError(t, err)
	info, err := c.StakingQuery.HistoricalInfo(context.Background(), &stakingtypes.QueryHistoricalInfoRequest{Height: 10000})
	require.NoError(t, err)
	for _, val := range info.Hist.Valset {
		t.Log(val.OperatorAddress)
		t.Log(val.Jailed)
	}
}

func TestSlakingQueryBlock(t *testing.T) {
	c, err := NewGRPCClient(UptickGrpcUrl, UptickRpc)
	require.NoError(t, err)
	Pagination := &query.PageRequest{
		Key:        []byte(""),
		Limit:      200,
		Offset:     0,
		CountTotal: false,
		Reverse:    false,
	}
	info, err := c.SlakingQuery.SigningInfos(context.Background(), &slakingtypes.QuerySigningInfosRequest{Pagination: Pagination})
	require.NoError(t, err)
	for _, val := range info.Info {
		t.Log(val.String())
	}
}

func TestQueryHttpEvent(t *testing.T) {
	c, err := NewGRPCClient(UptickGrpcUrl, UptickRpc)
	require.NoError(t, err)
	height := int64(300069)
	res, err := c.SignClient.BlockResults(context.Background(), &height)
	require.NoError(t, err)
	for _, event := range res.BeginBlockEvents {
		if event.Type == slakingtypes.EventTypeSlash {
			t.Log(event.Type)
			for _, attr := range event.Attributes {
				if bytes.Equal(attr.Key, []byte(slakingtypes.AttributeKeyJailed)) {
					consAddress, err := sdk.ConsAddressFromBech32(string(attr.Value))
					require.NoError(t, err)
					address := strings.ToUpper(hex.EncodeToString(consAddress.Bytes()))
					t.Log(address)
				}
			}
		}
	}
}

func TestQueryProposals(t *testing.T) {
	c, err := NewGRPCClient(UptickGrpcUrl, UptickRpc)
	require.NoError(t, err)

	proposals, err := c.GovQuery.Proposals(context.Background(), &govtypes.QueryProposalsRequest{})
	require.NoError(t, err)
	t.Logf("has proposal :%v", len(proposals.Proposals))
}

func TestQueryTxByEventType(t *testing.T) {
	c, err := NewGRPCClient(UptickGrpcUrl, UptickRpc)
	require.NoError(t, err)
	Pagination := &query.PageRequest{
		Key:        []byte(""),
		Limit:      200,
		Offset:     0,
		CountTotal: false,
		Reverse:    false,
	}
	res, err := c.TxClient.GetTxsEvent(context.Background(), &typestx.GetTxsEventRequest{Events: []string{"proposal_vote.proposal_id=5"}, Pagination: Pagination, OrderBy: 1})
	require.NoError(t, err)
	for _, tx := range res.Txs {
		for _, msg := range tx.Body.Messages {
			var mv sdk.Msg
			err = cdc.UnpackAny(msg, &mv)
			require.NoError(t, err)

			v := mv.(*govtypes.MsgVote)
			fmt.Println(v.Voter)
		}
	}
}

func TestQueryVoters(t *testing.T) {
	c, err := NewGRPCClient(UptickGrpcUrl, UptickRpc)
	require.NoError(t, err)
	voters, err := c.QueryVoters()
	require.NoError(t, err)
	for k, v := range voters {
		t.Logf("voter %v, %s", k, v)
	}
}

func TestQueryTxEvent(t *testing.T) {
	c, err := NewGRPCClient(UptickGrpcUrl, UptickRpc)
	require.NoError(t, err)
	height := int64(390733)

	for ; height < 4000000; height++ {
		res, err := c.TMServiceQuery.GetBlockByHeight(context.Background(), &tmservice.GetBlockByHeightRequest{Height: height})
		require.NoError(t, err)

		for _, tx := range res.Block.GetData().Txs {
			hash := hex.EncodeToString(tmhash.Sum(tx))
			res, err := c.TxClient.GetTx(context.Background(), &typestx.GetTxRequest{
				Hash: hash,
			})
			fmt.Println(res.TxResponse.Timestamp)
			require.NoError(t, err)
			if len(res.TxResponse.Logs) == 0 {
				continue
			}
			stringEvents := res.TxResponse.Logs[0].Events
			//var events []proto.Message
			for _, e := range stringEvents {
				abciEvent := abci.Event{}
				if e.Type == govtypes.EventTypeProposalVote {
					abciEvent.Type = e.Type
					for _, attr := range e.Attributes {
						abciEvent.Attributes = append(abciEvent.Attributes, abci.EventAttribute{
							Key:   []byte(attr.Key),
							Value: []byte(attr.Value),
						})
					}
					//protoEvent, err := sdk.ParseTypedEvent(abciEvent)
					require.NoError(t, err)
					fmt.Println(e.String())
					//events = append(events, protoEvent)
				}
			}

			if len(stringEvents) != 2 {
				continue
			}
			if stringEvents[0].Type == sdk.EventTypeMessage && stringEvents[1].Type == govtypes.EventTypeProposalVote {
				fmt.Printf("height:%d \n", height)
				fmt.Println(stringEvents[0].String())
				fmt.Println(stringEvents[1].String())
			}
		}

	}
}

func TestAddressConvert(t *testing.T) {
	valAddress, err := sdk.ValAddressFromBech32("uptickvaloper1c8y75a5nypmhngz5dktq9mjvp9d6auz9nel0up")
	require.NoError(t, err)
	AccAddress, err := sdk.AccAddressFromHex(hex.EncodeToString(valAddress.Bytes()))
	require.NoError(t, err)
	ValAddress, err := sdk.ValAddressFromHex(hex.EncodeToString(valAddress.Bytes()))
	require.NoError(t, err)
	ConsAddress, err := sdk.ConsAddressFromHex(hex.EncodeToString(valAddress.Bytes()))
	require.NoError(t, err)

	t.Logf("AccAddress:%s \n", AccAddress.String())
	t.Logf("ValAddress:%s \n", ValAddress.String())
	t.Logf("ConsAddress:%s \n", ConsAddress.String())

	t.Logf("Hex : %s \n", strings.ToUpper(hex.EncodeToString(AccAddress.Bytes())))
}

func TestQueryValVaa(t *testing.T) {
	c, err := NewGRPCClient(UptickGrpcUrl, UptickRpc)
	require.NoError(t, err)
	Pagination := &query.PageRequest{
		Key:        []byte(""),
		Limit:      20000,
		Offset:     0,
		CountTotal: false,
		Reverse:    false,
	}
	validators, err := c.StakingQuery.Validators(
		context.Background(),
		&stakingtypes.QueryValidatorsRequest{Pagination: Pagination},
	)
	require.NoError(t, err)
	for _, val := range validators.Validators {
		t.Logf(" address %s", val.OperatorAddress)

		var pubKey tmcrypto.PubKey
		var pk cryptotypes.PubKey

		err = cdc.UnpackAny(val.ConsensusPubkey, &pk)
		pubKey, err = cryptocodec.ToTmPubKeyInterface(pk)
		require.NoError(t, err)

		t.Logf(" pubkey %s", pk.Address())
		t.Logf(" pubkey %s", pubKey.Address())
	}
}

func TestQueryJailed(t *testing.T) {
	c, err := NewGRPCClient(UptickGrpcUrl, UptickRpc)
	require.NoError(t, err)
	jailed, err := c.QueryJailed(300069, 300069)
	require.NoError(t, err)
	for k, v := range jailed {
		t.Log("keys: ", k, "value:", v)
	}
}

func TestQueryValidatorInfos(t *testing.T) {
	c, err := NewGRPCClient(UptickGrpcUrl, UptickRpc)
	require.NoError(t, err)
	validators, err := c.QueryValidators()
	require.NoError(t, err)
	t.Log(len(validators))
	for k, v := range validators {
		t.Log("keys: ", k, "value:", v.String())
	}
}

func TestStringFloat(t *testing.T) {
	num := 123.12312412
	str := strconv.FormatFloat(num, 'f', 2, 64)
	f, err := strconv.ParseFloat(str, 64)
	require.NoError(t, err)
	require.Equal(t, 123.12, f)
}
