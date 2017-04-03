package repository

import (
	"github.com/mgierok/monujo/repository/entities"
)

func PortfoliosExt() (entities.PortfoliosExt, error) {
	portfolios := entities.PortfoliosExt{}
	err := Db().Select(&portfolios, "SELECT portfolio_id, name, currency, cache_value, outgoings, incomings, gain_of_sold_shares, commision, tax, gain_of_owned_shares, estimated_gain, estimated_gain_costs_inc FROM portfolios_ext")
	return portfolios, err
}
