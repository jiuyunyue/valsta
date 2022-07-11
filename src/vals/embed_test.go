package vals_test

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"testing"

	"github.com/jiuyunyue/valsta/cmd"

	"github.com/stretchr/testify/require"

	"github.com/jiuyunyue/valsta/src/vals"
)

func TestGetJson(t *testing.T) {
	up := vals.Uptime
	require.NotNil(t, up)
}

func Test_score(t *testing.T) {
	up := vals.Uptime
	require.NotNil(t, up)

	cmd.GrpcUrl = "peer0.testnet.uptick.network:9090"
	cmd.RpcUrl = "http://peer0.testnet.uptick.network:26657"

	voters, err := cmd.GetVoterInfos()
	require.NoError(t, err)

	score := make(map[string]uint64)
	for _, v := range voters {
		score[v.Address] = 20
	}
	for _, v := range up {
		sru := v.SurRate
		sruRate, err := strconv.ParseFloat(sru, 64)
		require.NoError(t, err)
		if !v.Jailed && sruRate > float64(80) {
			score[v.AccAddress] += 100
		} else if sruRate > float64(80) {
			score[v.AccAddress] += 80
		}
	}

	type UserScore struct {
		Address string
		Score   uint64
	}
	var scores []UserScore
	for k, v := range score {
		scores = append(scores, UserScore{k, v})
	}

	sort.Slice(scores, func(i, j int) bool {
		return scores[i].Score > scores[j].Score
	})
	marshal, err := json.Marshal(scores)
	require.NoError(t, err)
	t.Logf(string(marshal))
	t.Logf("total: %v ", len(score))

	addressList := make(map[string]bool)
	for _, v := range scores {
		addressList[v.Address] = true
		//fmt.Printf("%v : %v \n", v.Address, v.Score)
	}

	//f, err := os.Open("marketplace.txt")
	//require.NoError(t, err)
	//buf := bufio.NewReader(f)

	//for {
	//	line, err := buf.ReadString('\n')
	//	line = strings.TrimSpace(line)
	//	addressList[line] = true
	//	if err != nil {
	//		if err == io.EOF {
	//			break
	//		}
	//		require.NoError(t, err)
	//	}
	//}

	for k, _ := range addressList {
		fmt.Printf("%v \n", k)
	}
	fmt.Printf("total : %v \n", len(addressList))

}
