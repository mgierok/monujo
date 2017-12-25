package main

import (
	"flag"

	"github.com/mgierok/monujo/config"
	"github.com/mgierok/monujo/console"
	"github.com/mgierok/monujo/db"
)

func main() {
	var env string
	var dump string
	var file string
	flag.StringVar(&env, "env", "", "force environment")
	flag.StringVar(&dump, "dump", "", "dump 'data' or 'schema'")
	flag.StringVar(&file, "file", "", "where to store the dump")
	flag.Parse()

	conf, err := config.New(env)
	if err != nil {
		panic(err)
	}

	database, err := db.New(conf.Db())
	if err != nil {
		panic(err)
	}

	defer database.Connection().Close()

	if len(dump) > 0 {
		db.Dump(conf.Db(), conf.Sys(), dump, file)
	} else {
		console, _ := console.New()
		repository, _ := NewRepository(database.Connection())
		a, _ := NewApp(conf.App(), repository, console, console)
		a.Run()
	}
}
