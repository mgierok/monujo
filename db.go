package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/cep21/xdgbasedir"
)

type DbConfig struct {
	Host     string
	User     string
	Password string
	Dbname   string
}

func GetDb() *sql.DB {
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
	db, err := sql.Open("postgres", dbinfo)

	if err != nil {
		log.Fatal("Error: The data source arguments are not valid")
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Error: Could not establish a connection with the database")
	}

	return db
}
