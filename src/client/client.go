package client

import (
	"google.golang.org/grpc"

	tmclient "github.com/tendermint/tendermint/rpc/client"
	"github.com/tendermint/tendermint/rpc/client/http"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	grpc1 "github.com/gogo/protobuf/grpc"

	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/types/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	slakingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

var cdc *codec.ProtoCodec

func init() {
	cfg := sdk.GetConfig()
	cfg.SetBech32PrefixForAccount("uptick", "uptickpub")
	cfg.SetBech32PrefixForValidator("uptickvaloper", "uptickpub")
	cfg.SetBech32PrefixForConsensusNode("uptickvalcons", "uptickpub")

	registry := codectypes.NewInterfaceRegistry()
	cryptocodec.RegisterInterfaces(registry)
	cdc = codec.NewProtoCodec(registry)
}

type GClient struct {
	clientConn     grpc1.ClientConn
	AuthQuery      authtypes.QueryClient
	BankQuery      banktypes.QueryClient
	GovQuery       govtypes.QueryClient
	StakingQuery   stakingtypes.QueryClient
	SlakingQuery   slakingtypes.QueryClient
	TMServiceQuery tmservice.ServiceClient
	SignClient     tmclient.SignClient
	TxClient       tx.ServiceClient
}

func NewGRPCClient(url string, rpc string) (GClient, error) {
	dialOpts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	clientConn, err := grpc.Dial(url, dialOpts...)
	if err != nil {
		return GClient{}, err
	}
	httpClient, err := http.New(rpc, "/websocket")
	if err != nil {
		return GClient{}, err
	}
	return GClient{
		clientConn:     clientConn,
		StakingQuery:   stakingtypes.NewQueryClient(clientConn),
		BankQuery:      banktypes.NewQueryClient(clientConn),
		AuthQuery:      authtypes.NewQueryClient(clientConn),
		TMServiceQuery: tmservice.NewServiceClient(clientConn),
		SlakingQuery:   slakingtypes.NewQueryClient(clientConn),
		SignClient:     httpClient,
		TxClient:       tx.NewServiceClient(clientConn),
	}, nil
}
