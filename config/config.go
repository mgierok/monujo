package config

import (
	"encoding/json"
	"io/ioutil"

	"github.com/cep21/xdgbasedir"
)

type dbConf struct {
	Host     string
	User     string
	Password string
	Dbname   string
}

var db dbConf

func MustInitialize() {
	dbConfigFilePath, err := xdgbasedir.GetConfigFileLocation("monujo/db.json")
	if err != nil {
		panic(err)
	}

	dbConfigFile, err := ioutil.ReadFile(dbConfigFilePath)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(dbConfigFile, &db)
	if err != nil {
		panic(err)
	}
}

func Db() dbConf {
	return db
}
