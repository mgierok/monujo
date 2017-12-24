package repository

import (
	"github.com/mgierok/monujo/repository/entity"
)

func (r *Repository) PortfoliosExt() (entity.PortfoliosExt, error) {
	portfolios := entity.PortfoliosExt{}
	err := r.db.Select(&portfolios, "SELECT portfolio_id, name, currency, cache_value, gain_of_sold_shares, commision, tax, gain_of_owned_shares, estimated_gain, estimated_gain_costs_inc, estimated_value, annual_balance, month_balance FROM portfolios_ext ORDER BY portfolio_id")
	return portfolios, err
}

func (r *Repository) Portfolios() (entity.Portfolios, error) {
	portfolios := entity.Portfolios{}
	err := r.db.Select(&portfolios, "SELECT portfolio_id, name, currency FROM portfolios ORDER BY portfolio_id ASC")
	return portfolios, err
}
