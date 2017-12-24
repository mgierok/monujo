package config

import (
	"io/ioutil"

	"github.com/cep21/xdgbasedir"
	"github.com/json-iterator/go"
)

type Config struct {
	db  Db
	sys Sys
	app App
}

type Db struct {
	Host     string
	Port     string
	User     string
	Password string
	Dbname   string
}

type Sys struct {
	Pgdump string
}

type App struct {
	Alphavantagekey string
}

func New(env string) (*Config, error) {
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

	err = jsoniter.Unmarshal(dbConfigFile, &c.db)
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

	err = jsoniter.Unmarshal(sysConfigFile, &c.sys)
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

	err = jsoniter.Unmarshal(appConfigFile, &c.app)
	if err != nil {
		return c, err
	}

	return c, nil
}

func (c *Config) Db() Db {
	return c.db
}

func (c *Config) Sys() Sys {
	return c.sys
}

func (c *Config) App() App {
	return c.app
}
