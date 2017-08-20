package entity

import (
	"time"
)

type Operation struct {
	OperationId int64     `db:"operation_id"`
	PortfolioId int64     `db:"portfolio_id"`
	Date        time.Time `db:"date"`
	Type        string    `db:"type"`
	Value       float64   `db:"value"`
	Description string    `db:"description"`
	Commision   float64   `db:"commision"`
}

type Operations []Operation
