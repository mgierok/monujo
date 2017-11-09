package entity

import "database/sql"

type Security struct {
	Ticker        string         `db:"ticker"`
	ShortName     string         `db:"short_name"`
	FullName      string         `db:"full_name"`
	Market        string         `db:"market"`
	Leverage      float64        `db:"leverage"`
	QuotesSource  string         `db:"quotes_source"`
	TickerBankier sql.NullString `db:"ticker_bankier"`
}

type Securities []Security
