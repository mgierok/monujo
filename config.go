package main

import (
	"io/ioutil"

	"github.com/cep21/xdgbasedir"
	"github.com/json-iterator/go"
)

type Config struct {
	Db  DbConf
	Sys SysConf
	App AppConf
}

type DbConf struct {
	Host     string
	Port     string
	User     string
	Password string
	Dbname   string
}

type SysConf struct {
	Pgdump string
}

type AppConf struct {
	Alphavantagekey string
}

func NewConfig(env string) (*Config, error) {
	c := new(Config)
	var suffix string
	if len(env) > 0 {
		suffix = "." + env
	}

	dbConfigFilePath, err := xdgbasedir.GetConfigFileLocation("monujo/db.json" + suffix)
	if err != nil {
		return c, err
	}

	dbConfigFile, err := ioutil.ReadFile(dbConfigFilePath)
	if err != nil {
		return c, err
	}

	err = jsoniter.Unmarshal(dbConfigFile, &c.Db)
	if err != nil {
		return c, err
	}

	sysConfigFilePath, err := xdgbasedir.GetConfigFileLocation("monujo/sys.json" + suffix)
	if err != nil {
		return c, err
	}

	sysConfigFile, err := ioutil.ReadFile(sysConfigFilePath)
	if err != nil {
		return c, err
	}

	err = jsoniter.Unmarshal(sysConfigFile, &c.Sys)
	if err != nil {
		return c, err
	}

	appConfigFilePath, err := xdgbasedir.GetConfigFileLocation("monujo/app.json" + suffix)
	if err != nil {
		return c, err
	}

	appConfigFile, err := ioutil.ReadFile(appConfigFilePath)
	if err != nil {
		return c, err
	}

	err = jsoniter.Unmarshal(appConfigFile, &c.App)
	if err != nil {
		return c, err
	}

	return c, nil
}
