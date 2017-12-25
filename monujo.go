package main

import (
	"flag"
	"fmt"

	"github.com/jmoiron/sqlx"
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

	conf, err := NewConfig(env)
	if err != nil {
		panic(err)
	}

	dbinfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		conf.Db.Host,
		conf.Db.Port,
		conf.Db.User,
		conf.Db.Password,
		conf.Db.Dbname,
	)

	db, err := sqlx.Connect("postgres", dbinfo)

	if err != nil {
		panic(err)
	}

	defer db.Close()

	console, _ := console.New()
	repository, _ := NewRepository(db)
	a, _ := NewApp(conf, repository, console, console)

	if len(dump) > 0 {
		a.Dump(dump, file)
	} else {
		a.Run()
	}
}
