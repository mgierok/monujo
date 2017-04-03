package entities

import (
	"database/sql"
	"strings"
)

type OwnedStock struct {
	PortfolioId                sql.NullString `db:"portfolio_id"`
	PortfolioName              sql.NullString `db:"portfolio_name"`
	Ticker                     sql.NullString `db:"ticker"`
	ShortName                  sql.NullString `db:"short_name"`
	Shares                     sql.NullString `db:"shares"`
	LastPrice                  sql.NullString `db:"last_price"`
	MarketValue                sql.NullString `db:"market_value"`
	Currency                   sql.NullString `db:"currency"`
	ExchangeRate               sql.NullString `db:"exchange_rate"`
	LastPriceBaseCurrency      sql.NullString `db:"last_price_base_currency"`
	MarketValueBaseCurrency    sql.NullString `db:"market_value_base_currency"`
	AveragePrice               sql.NullString `db:"average_price"`
	AveragePriceBaseCurrency   sql.NullString `db:"average_price_base_currency"`
	Gain                       sql.NullString `db:"gain"`
	PercentageGain             sql.NullString `db:"percentage_gain"`
	GainBaseCurrency           sql.NullString `db:"gain_base_currency"`
	PercentageGainBaseCurrency sql.NullString `db:"percentage_gain_base_currency"`
}

func (stock OwnedStock) GetStockName() string {
	if stock.ShortName.String == "" {
		return strings.Trim(stock.Ticker.String, " ")
	} else {
		return stock.ShortName.String
	}
}
