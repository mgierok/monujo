package entities

import "strconv"

type Transaction struct {
	PortfolioId              int64   `db:"portfolio_id"`
	Date                     string  `db:"date"`
	Ticker                   string  `db:"ticker"`
	Price                    float64 `db:"price"`
	TransactionOperationType string  `db:"transaction_operation_type"`
	Currency                 string  `db:"currency"`
	Shares                   float64 `db:"shares"`
	Commision                float64 `db:"commision"`
	ExchangeRate             float64 `db:"exchange_rate"`
	Tax                      float64 `db:"tax"`
}

type Transactions []Transaction

func (t *Transaction) DisplayableArray() [][]string {
	return [][]string{
		[]string{"Portfolio ID", strconv.FormatInt(t.PortfolioId, 10)},
		[]string{"Date", t.Date},
		[]string{"Ticker", t.Ticker},
		[]string{"Price", strconv.FormatFloat(t.Price, 'f', -1, 64)},
		[]string{"Type", t.TransactionOperationType},
		[]string{"Currency", t.Currency},
		[]string{"Shares", strconv.FormatFloat(t.Shares, 'f', -1, 64)},
		[]string{"Commision", strconv.FormatFloat(t.Commision, 'f', -1, 64)},
		[]string{"Exchange Rate", strconv.FormatFloat(t.ExchangeRate, 'f', -1, 64)},
		[]string{"Tax", strconv.FormatFloat(t.Tax, 'f', -1, 64)},
	}
}
