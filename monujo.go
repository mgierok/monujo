package main

import (
	"flag"

	"github.com/mgierok/monujo/app"
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

	config.MustInitialize(env)
	db.MustInitialize()

	defer db.Connection().Close()

	if len(dump) > 0 {
		db.Dump(dump, file)
	} else {
		c, _ := console.New()
		a, _ := app.New(c, c)
		a.Run()
	}
}
