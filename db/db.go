package db

import (
	"fmt"

	"github.com/mgierok/monujo/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var db *sqlx.DB

func MustInitialize() {
	dbinfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Db().Host,
		config.Db().Port,
		config.Db().User,
		config.Db().Password,
		config.Db().Dbname,
	)

	var err error
	db, err = sqlx.Connect("postgres", dbinfo)

	if err != nil {
		panic(err)
	}
}

func Connection() *sqlx.DB {
	return db
}
