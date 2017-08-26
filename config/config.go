package config

import (
	"io/ioutil"

	"github.com/cep21/xdgbasedir"
	"github.com/json-iterator/go"
)

type dbConf struct {
	Host     string
	Port     string
	User     string
	Password string
	Dbname   string
}

type sysConf struct {
	Pgdump string
}

var db dbConf
var sys sysConf

func MustInitialize(env string) {
	var suffix string
	if len(env) > 0 {
		suffix = "." + env
	}

	dbConfigFilePath, err := xdgbasedir.GetConfigFileLocation("monujo/db.json" + suffix)
	if err != nil {
		panic(err)
	}

	dbConfigFile, err := ioutil.ReadFile(dbConfigFilePath)
	if err != nil {
		panic(err)
	}

	err = jsoniter.Unmarshal(dbConfigFile, &db)
	if err != nil {
		panic(err)
	}

	sysConfigFilePath, err := xdgbasedir.GetConfigFileLocation("monujo/sys.json" + suffix)
	if err != nil {
		panic(err)
	}

	sysConfigFile, err := ioutil.ReadFile(sysConfigFilePath)
	if err != nil {
		panic(err)
	}

	err = jsoniter.Unmarshal(sysConfigFile, &sys)
	if err != nil {
		panic(err)
	}
}

func Db() dbConf {
	return db
}

func Sys() sysConf {
	return sys
}
