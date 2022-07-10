package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"

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

		uptime, err := grpcClient.QueryUptime(run, end-1)
		if err != nil {
			return nil, err
		}
		jailed, err := grpcClient.QueryJailed(run, end-1)
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

		content, err := json.Marshal(uptime)
		if err != nil {
			return nil, err
		}
		err = ioutil.WriteFile(fmt.Sprintf("cache/%v_%v.json", run, end), content, 0777)
		if err != nil {
			return nil, err
		}

		// overwrite
		for k, v := range uptime {
			uptimeTmp := all[k]
			uptimeTmp.Times += v.Times
			if !uptimeTmp.Jailed {
				uptimeTmp.Jailed = v.Jailed
			}
		}
	}

	// recalculate
	for k, v := range all {
		tmp := all[k]
		num := float64(v.Times) / float64(endHeight-startHeight+1) * 100
		tmp.SurRate = strconv.FormatFloat(num, 'f', 2, 64)
		all[k] = tmp
	}

	content, err := json.Marshal(all)
	if err != nil {
		return nil, err
	}
	err = ioutil.WriteFile("cache/all.json", content, 0777)
	if err != nil {
		return nil, err
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
