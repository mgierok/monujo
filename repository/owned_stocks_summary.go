package repository

import (
	"github.com/mgierok/monujo/repository/entities"
)

func GetAllOwnedStocksSummary() ([]entities.OwnedStockSummary, error) {
	ownedStocksSummary := []entities.OwnedStockSummary{}
	err := Db().Select(&ownedStocksSummary, "SELECT portfolio_id, portfolio_name, ticker, short_name, shares, last_price, market_value, currency, exchange_rate, last_price_base_currency, market_value_base_currency, average_price, average_price_base_currency, gain, percentage_gain, gain_base_currency, percentage_gain_base_currency  FROM owned_shares_summary")
	return ownedStocksSummary, err
}
