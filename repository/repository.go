package repository

import "github.com/jmoiron/sqlx"

var db *sqlx.DB

func SetDb(dbConnection *sqlx.DB) {
	db = dbConnection
}

func Db() *sqlx.DB {
	return db
}
