package models

import (
	"database/sql"
)

type PortfolioSummary struct {
	PortfolioId           sql.NullString `db:"portfolio_id"`
	Name                  sql.NullString `db:"name"`
	Currency              sql.NullString `db:"currency"`
	CacheValue            sql.NullString `db:"cache_value"`
	Outgoings             sql.NullString `db:"outgoings"`
	Incomings             sql.NullString `db:"incomings"`
	GainOfSoldShares      sql.NullString `db:"gain_of_sold_shares"`
	Commision             sql.NullString `db:"commision"`
	Tax                   sql.NullString `db:"tax"`
	GainOfOwnedShares     sql.NullString `db:"gain_of_owned_shares"`
	EstimatedGain         sql.NullString `db:"estimated_gain"`
	EstimatedGainCostsInc sql.NullString `db:"estimated_gain_costs_inc"`
}
