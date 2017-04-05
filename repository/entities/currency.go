package entities

type Currency struct {
	Symbol string `db:"currency"`
}

type Currencies []Currency
