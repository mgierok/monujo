package monujo

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

func Connect(c dbConf) (*sqlx.DB, error) {
	dbinfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.Host,
		c.Port,
		c.User,
		c.Password,
		c.Dbname,
	)

	return sqlx.Connect("postgres", dbinfo)
}
