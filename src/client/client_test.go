package client

import (
	"bytes"
	"context"
	"encoding/hex"
	"strconv"
	"strings"
	"testing"

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
	UptickGrpcUrl = "peer1.testnet.uptick.network:9090"
	LocalGrpc     = "localhost:9090"

	UptickRpc = "http://peer1.testnet.uptick.network:26657"
	LocalRpc  = "http://localhost:26657"
)

func TestQueryBalance(t *testing.T) {
	c, err := NewGRPCClient(UptickGrpcUrl, UptickGrpcUrl)
	require.NoError(t, err)
	res, _ := c.BankQuery.Balance(context.Background(), &types.QueryBalanceRequest{Address: "uptick10t4kkjetahnjh5d8h2d6dqnp7cvaestxxyvjgw", Denom: "auptick"})
	t.Log(res.String())
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

func TestAddressConvert(t *testing.T) {
	consAddress, err := sdk.ConsAddressFromBech32("uptickvalcons1ksw5fem64junr3330ssy65l9uj0xkg9k3nerlw")
	require.NoError(t, err)
	AccAddress, err := sdk.AccAddressFromHex(hex.EncodeToString(consAddress.Bytes()))
	require.NoError(t, err)
	ValAddress, err := sdk.ValAddressFromHex(hex.EncodeToString(consAddress.Bytes()))
	require.NoError(t, err)
	ConsAddress, err := sdk.ConsAddressFromHex(hex.EncodeToString(consAddress.Bytes()))
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
		Limit:      200,
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

		t.Logf(" pubkey %s", pk.Address())
		t.Logf(" pubkey %s", pubKey.Address())
	}
}

func TestQueryJailed(t *testing.T) {
	c, err := NewGRPCClient(UptickGrpcUrl, UptickRpc)
	require.NoError(t, err)
	jailed := c.QueryJailed(300069, 300069)
	for k, v := range jailed {
		t.Log("keys: ", k, "value:", v)
	}
}

func TestQueryValidatorInfos(t *testing.T) {
	c, err := NewGRPCClient(UptickGrpcUrl, UptickRpc)
	require.NoError(t, err)
	validators := c.QueryValidators()
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
