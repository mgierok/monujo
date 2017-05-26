package entity

import (
	"fmt"
	"strings"

	"database/sql"
)

type Stock struct {
	Ticker    string          `db:"ticker"`
	ShortName sql.NullString  `db:"short_name"`
	LastPrice sql.NullFloat64 `db:"last_price"`
	Currency  string          `db:"currency"`
}

type OwnedStock struct {
	Stock
	PortfolioId                int64           `db:"portfolio_id"`
	PortfolioName              string          `db:"portfolio_name"`
	Shares                     float64         `db:"shares"`
	ExchangeRate               sql.NullFloat64 `db:"exchange_rate"`
	LastPriceBaseCurrency      sql.NullFloat64 `db:"last_price_base_currency"`
	AveragePrice               float64         `db:"average_price"`
	AveragePriceBaseCurrency   float64         `db:"average_price_base_currency"`
	Gain                       sql.NullFloat64 `db:"gain"`
	PercentageGain             sql.NullFloat64 `db:"percentage_gain"`
	GainBaseCurrency           sql.NullFloat64 `db:"gain_base_currency"`
	PercentageGainBaseCurrency sql.NullFloat64 `db:"percentage_gain_base_currency"`
}

type OwnedStocks []OwnedStock

func (stock *Stock) DisplayName() string {
	if stock.ShortName.Valid {
		return fmt.Sprintf("%s (%s)", stock.ShortName.String, strings.Trim(stock.Ticker, " "))
	} else {
		return strings.Trim(stock.Ticker, " ")
	}
}
