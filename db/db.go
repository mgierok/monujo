package db

import (
	"fmt"

	"github.com/mgierok/monujo/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var db *sqlx.DB

type Db struct {
	connection *sqlx.DB
}

func New(c config.Db) (*Db, error) {
	dbinfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.Host,
		c.Port,
		c.User,
		c.Password,
		c.Dbname,
	)

	db, err := sqlx.Connect("postgres", dbinfo)

	return &Db{
		connection: db,
	}, err
}

func (d *Db) Connection() *sqlx.DB {
	return d.connection
}
