package entities

import "database/sql"

type Currency struct {
	Symbol sql.NullString `db:"currency"`
}

type Currencies []Currency
