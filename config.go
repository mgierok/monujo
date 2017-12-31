package monujo

import (
	"io/ioutil"

	"github.com/cep21/xdgbasedir"
	"github.com/json-iterator/go"
)

type config struct {
	db  dbConf
	sys sysConf
	app appConf
}

type Config interface {
	Db() dbConf
	Sys() sysConf
	App() appConf
}

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

type appConf struct {
	Alphavantagekey string
}

func NewConfig(env string) (Config, error) {
	c := new(config)
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

func (c *config) Db() dbConf {
	return c.db
}

func (c *config) Sys() sysConf {
	return c.sys
}

func (c *config) App() appConf {
	return c.app
}
