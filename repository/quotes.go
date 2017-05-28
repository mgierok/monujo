package repository

import (
	"github.com/mgierok/monujo/db"
	"github.com/mgierok/monujo/repository/entity"
)

func StoreLatestQuote(quote entity.Quote) (string, error) {
	stmt, err := db.Connection().PrepareNamed(`
		INSERT INTO latest_quotes (ticker, date, open, high, low, close, volume, openint)
		VALUES (:ticker, :date, :open, :high, :low, :close, :volume, :openint)
		RETURNING ticker
	`)

	var t string
	if nil == err {
		err = stmt.Get(&t, quote)
	}

	return t, err
}
