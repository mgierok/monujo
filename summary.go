package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/olekukonko/tablewriter"
)

type OwnedStockSummary struct {
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

func Summary(db *sqlx.DB) {
	ownedStocksSummary := []OwnedStockSummary{}
	err := db.Select(&ownedStocksSummary, "SELECT portfolio_id, portfolio_name, ticker, short_name, shares, last_price, market_value, currency, exchange_rate, last_price_base_currency, market_value_base_currency, average_price, average_price_base_currency, gain, percentage_gain, gain_base_currency, percentage_gain_base_currency  FROM owned_shares_summary")
	LogError(err)

	var data [][]string

	for _, oss := range ownedStocksSummary {
		var stock string
		if oss.ShortName.String == "" {
			stock = strings.Trim(oss.Ticker.String, " ")
		} else {
			stock = oss.ShortName.String
		}

		data = append(data, []string{
			oss.PortfolioName.String,
			stock,
			oss.Shares.String,
			oss.LastPrice.String,
			oss.AveragePrice.String,
			oss.LastPriceBaseCurrency.String,
			oss.AveragePriceBaseCurrency.String,
			oss.Gain.String,
			oss.GainBaseCurrency.String,
			oss.PercentageGain.String,
			oss.PercentageGainBaseCurrency.String,
		})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"Portfolio Name",
		"Stock",
		"Shares",
		"Last Price",
		"Average Price",
		"Last Price BC",
		"Average Price BC",
		"Gain",
		"Gain BC",
		"Gain%",
		"Gain BC%",
	})

	table.AppendBulk(data)
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)
	table.Render()

	data = data[0:0]
	fmt.Println("")
	fmt.Println("")

	portfoliosSummary := []PortfolioSummary{}
	err = db.Select(&portfoliosSummary, "SELECT portfolio_id, name, currency, cache_value, outgoings, incomings, gain_of_sold_shares, commision, tax, gain_of_owned_shares, estimated_gain, estimated_gain_costs_inc FROM portfolios_summary")
	LogError(err)

	for _, ps := range portfoliosSummary {
		data = append(data, []string{
			ps.PortfolioId.String,
			ps.Name.String,
			ps.CacheValue.String,
			ps.GainOfSoldShares.String,
			ps.Commision.String,
			ps.Tax.String,
			ps.GainOfOwnedShares.String,
			ps.EstimatedGain.String,
			ps.EstimatedGainCostsInc.String,
		})
	}

	table = tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"Portfolio Id",
		"Portfolio Name",
		"Cache Value",
		"Gain of Sold Shares",
		"Commision",
		"Tax",
		"Gain Of Ownded Shares",
		"Estimated Gain",
		"Estimated Gain Costs Inc.",
	})

	table.AppendBulk(data)
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)
	table.Render()
}
