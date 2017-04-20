package repository

import (
	"github.com/mgierok/monujo/repository/entities"
)

func OwnedStocks() (entities.OwnedStocks, error) {
	stocks := entities.OwnedStocks{}
	err := Db().Select(&stocks, "SELECT portfolio_id, portfolio_name, ticker, short_name, shares, last_price, currency, exchange_rate, last_price_base_currency, average_price, average_price_base_currency, gain, percentage_gain, gain_base_currency, percentage_gain_base_currency  FROM owned_stocks")
	return stocks, err
}
