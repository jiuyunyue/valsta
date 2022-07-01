package database

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/jiuyunyue/valsta/src/client"
)

/*
	docker pull mysql:oracle
	docker run -itd --name docker-mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=123456 -v ~/Docker/db/mysql/8.0.9:/var/lib/mysql mysql:oracle
	docker run --name docker-mysql --privileged=true -d -v ~/Docker/mysql/data/:/var/lib/mysql -v ~/Docker/mysql/conf:/etc/mysql/conf.d -v ~/Docker/mysql/log:/var/log/mysql -p 3306:3306  -e MYSQL_ROOT_PASSWORD=123456  mysql:oracle
*/

func Test_connect(t *testing.T) {
	defer db.Close()
	db.SetConnMaxLifetime(100)
	db.SetMaxIdleConns(10)
	err := db.Ping()
	require.NoError(t, err)
}

func TestDropTable(t *testing.T) {
	defer db.Close()
	DropTable()
}

func TestCreateValTable(t *testing.T) {
	defer db.Close()
	CreateTable()
}

func TestInsertIntoValInfos(t *testing.T) {
	defer db.Close()
	address := "test"
	valAddress := "valAddress"
	surRate := "90.12"
	times := 99
	err := InsertIntoValInfos(address, valAddress, surRate, times, false)
	require.NoError(t, err)
}

func TestGetValidatorInfos(t *testing.T) {
	defer db.Close()
	validatorInfos, err := GetValidatorInfos()
	require.NoError(t, err)
	for _, validatorInfo := range validatorInfos {
		t.Log(validatorInfo.String())
	}
}

func TestQueryOne(t *testing.T) {
	defer db.Close()
	one, err := QueryOne("test")
	require.NoError(t, err)
	t.Log(one.String())
}

func TestCleanIntoValInfosData(t *testing.T) {
	defer db.Close()
	err := CleanIntoValInfosData()
	require.NoError(t, err)
}

func TestContextInsertIntoValInfos(t *testing.T) {
	defer db.Close()
	err := CleanIntoValInfosData()
	require.NoError(t, err)

	grpcClient, err := client.NewGRPCClient("peer1.testnet.uptick.network:9090", "http://peer1.testnet.uptick.network:26657")
	require.NoError(t, err)

	start := int64(300069)
	end := int64(300069)
	uptime, _ := grpcClient.QueryUptime(start, end)
	jailed, _ := grpcClient.QueryJailed(start, end)
	for k, v := range jailed {
		val := uptime[k]
		if v == true && val.Jailed == false && len(val.AccAddress) != 0 {
			val.Jailed = true
			uptime[k] = val
		}
	}

	err = ContextInsertIntoValInfos(uptime)
	require.NoError(t, err)
}
