package entities

import "time"

type Transaction struct {
	PortfolioId  int64     `db:"portfolio_id"`
	Date         time.Time `db:"date"`
	Ticker       string    `db:"ticker"`
	Price        float64   `db:"price"`
	Currency     string    `db:"currency"`
	Shares       float64   `db:"shares"`
	Commision    float64   `db:"commision"`
	ExchangeRate float64   `db:"exchange_rate"`
	Tax          float64   `db:"tax"`
}

type Transactions []Transaction
