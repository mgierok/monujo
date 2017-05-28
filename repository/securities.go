package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/mgierok/monujo/db"
	"github.com/mgierok/monujo/repository/entity"
)

func Securities(tickers []string) (entity.Securities, error) {
	s := entity.Securities{}
	var query string
	var err error

	if len(tickers) > 0 {
		var args []interface{}
		query, args, err = sqlx.In(
			`SELECT
				ticker,
				short_name,
				full_name,
				market,
				leverage,
				quotes_source
			FROM securities
			WHERE ticker IN (?)
			`,
			tickers,
		)
		query = db.Connection().Rebind(query)
		err = db.Connection().Select(&s, query, args...)
	} else {
		query =
			`SELECT
				ticker,
				short_name,
				full_name,
				market,
				leverage,
				quotes_source
			FROM securities`
		err = db.Connection().Select(&s, query)
	}

	return s, err
}
