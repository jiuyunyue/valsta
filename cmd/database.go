package cmd

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/jiuyunyue/valsta/src/database"
	"github.com/jiuyunyue/valsta/src/types"
)

func Init() error {
	db, err := sqlx.Connect(database.Mysql, database.LocalUrl+database.Database)
	if err != nil {
		return err
	}
	defer db.Close()
	db.SetConnMaxLifetime(100)
	db.SetMaxIdleConns(10)
	err = db.Ping()
	if err != nil {
		return err
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS " +
		" validator_infos " +
		"( " +
		" address varchar(40) PRIMARY KEY ," +
		" acc_address varchar(45) NOT NULL," +
		" jailed boolean NOT NULL ," +
		" times integer NOT NULL ," +
		" sur_rate varchar(10) NOT NULL )")

	if err != nil {
		return err
	}
	return nil
}

func CleanDatabase() {
	db, err := sqlx.Connect(database.Mysql, database.LocalUrl+database.Database)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db.SetConnMaxLifetime(100)
	db.SetMaxIdleConns(10)
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(fmt.Sprintf("DROP table %s", database.TableValInfo))
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(fmt.Sprintf("DROP database %s", database.Database))
	if err != nil {
		panic(err)
	}
}

func GetValInfos() ([]types.ValidatorInfo, error) {
	db := database.GetDB()
	defer db.Close()

	infos, err := database.GetValidatorInfos()
	if err != nil {
		return nil, err
	}
	return infos, nil
}
