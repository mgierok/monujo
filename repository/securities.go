package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/mgierok/monujo/repository/entity"
)

func (r *Repository) SecurityExists(ticker string) (bool, error) {
	var exists bool
	err := r.db.Get(
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

func (r *Repository) Securities(tickers []string) (entity.Securities, error) {
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
		query = r.db.Rebind(query)
		err = r.db.Select(&s, query, args...)
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
		err = r.db.Select(&s, query)
	}

	return s, err
}

func (r *Repository) StoreSecurity(s entity.Security) (string, error) {
	stmt, err := r.db.PrepareNamed(`
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
