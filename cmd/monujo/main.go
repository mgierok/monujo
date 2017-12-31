package main

import (
	"flag"

	"github.com/mgierok/monujo"
	"github.com/mgierok/monujo/console"

	_ "github.com/lib/pq"
)

func main() {
	var env string
	var dump string
	var file string
	flag.StringVar(&env, "env", "", "force environment")
	flag.StringVar(&dump, "dump", "", "dump 'data' or 'schema'")
	flag.StringVar(&file, "file", "", "where to store the dump")
	flag.Parse()

	conf, err := monujo.NewConfig(env)
	if err != nil {
		panic(err)
	}

	db, err := monujo.Connect(conf.Db())
	if err != nil {
		panic(err)
	}

	defer db.Close()

	console, _ := console.New()
	repository, _ := monujo.NewRepository(db, conf)
	a, _ := monujo.New(conf, repository, console, console)

	if len(dump) > 0 {
		a.Dump(dump, file)
	} else {
		a.Run()
	}
}
