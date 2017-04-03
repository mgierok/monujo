package entities

import (
	"database/sql"
	"strings"
)

type Stock struct {
	Ticker    sql.NullString `db:"ticker"`
	ShortName sql.NullString `db:"short_name"`
	LastPrice sql.NullString `db:"last_price"`
	Currency  sql.NullString `db:"currency"`
}

type OwnedStock struct {
	Stock
	PortfolioId                sql.NullString `db:"portfolio_id"`
	PortfolioName              sql.NullString `db:"portfolio_name"`
	Shares                     sql.NullString `db:"shares"`
	MarketValue                sql.NullString `db:"market_value"`
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

type OwnedStocks []OwnedStock

func (stock *Stock) GetStockName() string {
	if stock.ShortName.String == "" {
		return strings.Trim(stock.Ticker.String, " ")
	} else {
		return stock.ShortName.String
	}
}
