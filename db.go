package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	_ "github.com/lib/pq"
	"github.com/cep21/xdgbasedir"
	"github.com/jmoiron/sqlx"
)

type DbConfig struct {
	Host     string
	User     string
	Password string
	Dbname   string
}

func GetDbConnection() *sqlx.DB {
	dbConfigFilePath, err := xdgbasedir.GetConfigFileLocation("monujo/db.json")
	dbConfigFile, err := ioutil.ReadFile(dbConfigFilePath)
	dbConfig := DbConfig{}
	json.Unmarshal(dbConfigFile, &dbConfig)

	dbinfo := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s sslmode=disable",
		dbConfig.Host,
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Dbname,
	)
	db, err := sqlx.Connect("postgres", dbinfo)

	if err != nil {
		log.Fatalln(err)
	}

	return db
}
