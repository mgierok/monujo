package entity

import "time"

type Quote struct {
	Ticker  string    `db:"ticker"`
	Date    time.Time `db:"date"`
	Open    float64   `db:"open"`
	High    float64   `db:"high"`
	Low     float64   `db:"low"`
	Close   float64   `db:"close"`
	Volume  float64   `db:"volume"`
	OpenInt float64   `db:"openint"`
}

type Quotes []Quote
