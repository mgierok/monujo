package repository

import (
	"github.com/mgierok/monujo/repository/entity"
)

func (r *Repository) StoreTransaction(transaction entity.Transaction) (int64, error) {
	stmt, err := r.db.PrepareNamed(`
		INSERT INTO transactions (portfolio_id, date, ticker, price, currency, shares, commision, exchange_rate, tax)
		VALUES (:portfolio_id, :date, :ticker, :price, :currency, :shares, :commision, :exchange_rate, :tax)
		RETURNING transaction_id
	`)

	var transactionId int64
	if nil == err {
		err = stmt.Get(&transactionId, transaction)
	}

	return transactionId, err
}

func (r *Repository) PortfolioTransactions(portfolio entity.Portfolio) (entity.Transactions, error) {
	transactions := entity.Transactions{}
	err := r.db.Select(&transactions,
		`SELECT
		transaction_id,
		portfolio_id,
		date,
		ticker,
		price,
		currency,
		shares,
		commision,
		exchange_rate,
		tax
	FROM transactions
	WHERE portfolio_id = $1
	ORDER BY
		date ASC,
		transaction_id ASC
	`,
		portfolio.PortfolioId)
	return transactions, err
}

func (r *Repository) DeleteTransaction(transaction entity.Transaction) error {
	_, err := r.db.Exec("DELETE FROM transactions WHERE portfolio_id = $1 AND transaction_id = $2", transaction.PortfolioId, transaction.TransactionId)
	return err
}
