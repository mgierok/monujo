package repository

import (
	"github.com/mgierok/monujo/repository/entities"
)

func StoreTransaction(transaction entities.Transaction) (int64, error) {
	stmt, err := Db().PrepareNamed(`
		INSERT INTO transactions (portfolio_id, date, ticker, price, type, currency, shares, commision, exchange_rate, tax)
		VALUES (:portfolio_id, :date, :ticker, :price, :transaction_operation_type, :currency, :shares, :commision, :exchange_rate, :tax)
		RETURNING transaction_id
	`)

	var transactionId int64
	if nil == err {
		stmt.Get(&transactionId, transaction)
	}

	return transactionId, err
}
