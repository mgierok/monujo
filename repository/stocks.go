package repository

import (
	"github.com/mgierok/monujo/db"
	"github.com/mgierok/monujo/repository/entity"
)

func OwnedStocks() (entity.OwnedStocks, error) {
	stocks := entity.OwnedStocks{}
	err := db.Connection().Select(&stocks,
		`SELECT
			portfolio_id,
			portfolio_name,
			ticker,
			short_name,
			shares,
			last_price,
			currency,
			exchange_rate,
			last_price_base_currency,
			average_price,
			average_price_base_currency,
			average_price_base_currency,
			investment_base_currency,
			market_value_base_currency,
			gain,
			percentage_gain,
			gain_base_currency,
			percentage_gain_base_currency
			FROM owned_stocks
			`)
	return stocks, err
}
