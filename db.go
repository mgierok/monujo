package main

import (
	"fmt"
	"log"

	"github.com/mgierok/monujo/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func GetDbConnection() *sqlx.DB {
	dbinfo := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Db().Host,
		config.Db().User,
		config.Db().Password,
		config.Db().Dbname,
	)
	db, err := sqlx.Connect("postgres", dbinfo)

	if err != nil {
		log.Fatalln(err)
	}

	return db
}
