package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/mgierok/monujo/db"
	"github.com/mgierok/monujo/repository/entity"
)

func SecurityExists(ticker string) (bool, error) {
	var exists bool
	err := db.Connection().Get(
		&exists,
		`SELECT
			COUNT(1)
		FROM securities
		WHERE ticker = $1
		`,
		ticker,
	)
	return exists, err
}

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
				quotes_source,
				ticker_bankier
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
				quotes_source,
				ticker_bankier
			FROM securities`
		err = db.Connection().Select(&s, query)
	}

	return s, err
}

func StoreSecurity(s entity.Security) (string, error) {
	stmt, err := db.Connection().PrepareNamed(`
		INSERT INTO securities (ticker, short_name, full_name, market, leverage, quotes_source, ticker_bankier)
		VALUES (:ticker, :short_name, :full_name, :market, :leverage, :quotes_source, :ticker_bankier)
		RETURNING ticker
	`)

	var t string
	if nil == err {
		err = stmt.Get(&t, s)
	}

	return t, err
}
