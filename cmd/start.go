package cmd

import (
	"fmt"
	"io/ioutil"

	"github.com/tendermint/tendermint/libs/json"

	"github.com/jiuyunyue/valsta/src/client"
	"github.com/jiuyunyue/valsta/src/database"
	"github.com/jiuyunyue/valsta/src/types"
)

var GrpcUrl string
var RpcUrl string

const CacheNum = 10000

func ValSta(startHeight, endHeight int64) ([]types.ValidatorInfo, error) {
	grpcClient, err := client.NewGRPCClient(GrpcUrl, RpcUrl)
	if err != nil {
		return nil, err
	}

	all := make(types.Uptime)
	times := (endHeight-startHeight)/CacheNum + 1
	for tmp := int64(0); tmp < times; tmp++ {
		run := CacheNum*tmp + startHeight
		end := CacheNum + run
		if end > endHeight {
			end = endHeight
		}
		if run > endHeight {
			run = endHeight
		}

		uptime, err := grpcClient.QueryUptime(run, end)
		if err != nil {
			return nil, err
		}
		jailed, err := grpcClient.QueryJailed(run, end)
		if err != nil {
			return nil, err
		}

		for k, v := range jailed {
			val := uptime[k]
			if v == true && val.Jailed == false && len(val.AccAddress) != 0 {
				val.Jailed = true
				uptime[k] = val
			}
		}

		//var noJail []string
		//for k, v := range uptime {
		//	if v.Jailed == false {
		//		noJail = append(noJail, k)
		//	}
		//}

		content, err := json.Marshal(uptime)
		if err != nil {
			return nil, err
		}
		err = ioutil.WriteFile(fmt.Sprintf("%v_%v.txt", run, end), content, 0777)
		if err != nil {
			return nil, err
		}

		// overwrite
		for k, v := range uptime {
			all[k] = v
		}
	}

	// set with db
	db := database.GetDB()
	defer db.Close()
	err = database.ContextInsertIntoValInfos(all)
	if err != nil {
		return nil, err
	}

	validatorInfos, err := database.GetValidatorInfos()
	if err != nil {
		return nil, err
	}
	return validatorInfos, nil
}

func GetVoterInfos() (voters map[string]types.VoterInfo, err error) {
	grpcClient, err := client.NewGRPCClient(GrpcUrl, RpcUrl)
	if err != nil {
		return nil, err
	}
	voters, err = grpcClient.QueryVoters()
	if err != nil {
		return nil, err
	}
	return voters, nil
}
