package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/mgierok/monujo/repository/models"
)

func GetAllPortfoliosSummary(db *sqlx.DB) ([]models.PortfolioSummary, error) {
	portfoliosSummary := []models.PortfolioSummary{}
	err := db.Select(&portfoliosSummary, "SELECT portfolio_id, name, currency, cache_value, outgoings, incomings, gain_of_sold_shares, commision, tax, gain_of_owned_shares, estimated_gain, estimated_gain_costs_inc FROM portfolios_summary")
	return portfoliosSummary, err
}
