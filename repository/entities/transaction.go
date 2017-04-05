package entities

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
