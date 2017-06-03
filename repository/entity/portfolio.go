package entity

import (
	"database/sql"
)

type Portfolio struct {
	PortfolioId int64  `db:"portfolio_id"`
	Name        string `db:"name"`
	Currency    string `db:"currency"`
}

type Portfolios []Portfolio

type PortfolioExt struct {
	Portfolio
	CacheValue            sql.NullFloat64 `db:"cache_value"`
	GainOfSoldShares      sql.NullFloat64 `db:"gain_of_sold_shares"`
	Commision             sql.NullFloat64 `db:"commision"`
	Tax                   sql.NullFloat64 `db:"tax"`
	GainOfOwnedShares     sql.NullFloat64 `db:"gain_of_owned_shares"`
	EstimatedGain         sql.NullFloat64 `db:"estimated_gain"`
	EstimatedGainCostsInc sql.NullFloat64 `db:"estimated_gain_costs_inc"`
	EstimatedValue        sql.NullFloat64 `db:"estimated_value"`
	AnnualBalance         float64         `db:"annual_balance"`
	MonthBalance          float64         `db:"month_balance"`
}

type PortfoliosExt []PortfolioExt
