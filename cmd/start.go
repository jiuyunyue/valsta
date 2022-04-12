package cmd

import (
	"github.com/jiuyunyue/valsta/src/client"
	"github.com/jiuyunyue/valsta/src/database"
	"github.com/jiuyunyue/valsta/src/types"
)

var GrpcUrl string
var RpcUrl string

func ValSta(startHeight, endHeight int64) ([]types.ValidatorInfo, error) {
	grpcClient, err := client.NewGRPCClient(GrpcUrl, RpcUrl)
	if err != nil {
		return nil, err
	}
	uptime := grpcClient.QueryUptime(startHeight, endHeight)
	jailed := grpcClient.QueryJailed(startHeight, endHeight)

	for k, v := range jailed {
		val := uptime[k]
		if v == true && val.Jailed == false && len(val.AccAddress) != 0 {
			val.Jailed = true
			uptime[k] = val
		}
	}

	var noJail []string
	for k, v := range uptime {
		if v.Jailed == false {
			noJail = append(noJail, k)
		}
	}

	// set with db
	db := database.GetDB()
	defer db.Close()
	err = database.ContextInsertIntoValInfos(uptime)
	if err != nil {
		return nil, err
	}

	validatorInfos, err := database.GetValidatorInfos()
	if err != nil {
		return nil, err
	}
	return validatorInfos, nil
}
