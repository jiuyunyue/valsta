package database

import (
	"database/sql"
	"fmt"

	"github.com/jiuyunyue/valsta/src/types"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const (
	Mysql        = "mysql"
	Database     = "valsta"
	TableValInfo = "validator_infos"
	LocalUrl     = "root:123456@tcp(127.0.0.1:3306)/"
)

var db *sqlx.DB

func GetDB() *sqlx.DB {
	return db
}

func init() {
	CreateDatabase()
	database, err := sqlx.Connect(Mysql, LocalUrl+Database)
	if err != nil {
		panic(err)
	}
	db = database
}

func CreateDatabase() {
	DB, err := sql.Open(Mysql, LocalUrl)
	defer DB.Close()
	if err != nil {
		panic(err)
	}
	_, err = DB.Exec("CREATE DATABASE IF NOT EXISTS valsta")
	if err != nil {
		panic(err)
	}
}

func CreateTable() {
	// Create tables
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS " +
		" validator_infos " +
		"( " +
		" address varchar(40) PRIMARY KEY ," +
		" acc_address varchar(45) NOT NULL," +
		" jailed boolean NOT NULL ," +
		" times integer NOT NULL ," +
		" sur_rate varchar(10) NOT NULL )")

	if err != nil {
		panic(err)
	}
}

func DropTable() {
	_, err := db.Exec("DROP table validator_infos")
	if err != nil {
		panic(err)
	}
}

func InsertIntoValInfos(address, valAddress, surRate string, times int, jailed bool) error {
	_, err := db.Exec(
		fmt.Sprintf("INSERT into validator_infos (address,times,acc_address,jailed,sur_rate) VALUES ('%s',%d,'%s',%t,'%s')",
			address,
			times,
			valAddress,
			jailed,
			surRate,
		))
	if err != nil {
		return err
	}
	return nil
}

func ContextInsertIntoValInfos(validatorInfos map[string]types.ValidatorInfo) error {
	conn, err := db.Begin()
	if err != nil {
		return err
	}
	for _, v := range validatorInfos {
		res, err := db.Exec(
			fmt.Sprintf("INSERT into validator_infos (address,times,acc_address,jailed,sur_rate) VALUES ('%s',%d,'%s',%t,'%s')",
				v.Address,
				v.Times,
				v.AccAddress,
				v.Jailed,
				v.SurRate,
			))
		if err != nil {
			conn.Rollback()
			return err
		}
		_, err = res.LastInsertId()
		if err != nil {
			conn.Rollback()
			return err
		}
	}
	conn.Commit()

	return nil
}

func CleanIntoValInfosData() error {
	// truncate table validator_infos
	_, err := db.Exec("truncate table validator_infos")
	if err != nil {
		return err
	}
	return nil
}

func GetValidatorInfos() ([]types.ValidatorInfo, error) {
	var validatorInfos []types.ValidatorInfo
	rows, err := db.Query("select * from validator_infos")
	if err != nil {
		return nil, err
	}
	validatorInfo := new(types.ValidatorInfo)
	for rows.Next() {
		err = rows.Scan(&validatorInfo.Address, &validatorInfo.AccAddress, &validatorInfo.Jailed, &validatorInfo.Times, &validatorInfo.SurRate)
		if err != nil {
			return nil, err
		}
		validatorInfos = append(validatorInfos, *validatorInfo)
	}
	return validatorInfos, nil
}

func QueryOne(address string) (*types.ValidatorInfo, error) {
	valInfo := new(types.ValidatorInfo)
	row := db.QueryRow("select * from validator_infos where address=?", address)
	if err := row.Scan(&valInfo.Address, &valInfo.AccAddress, &valInfo.Jailed, &valInfo.Times, &valInfo.SurRate); err != nil {
		return nil, err
	}
	return valInfo, nil
}
