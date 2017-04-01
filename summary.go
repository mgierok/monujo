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

type OwnedSharesSummary struct {
	portfolioId                sql.NullString
	portfolioName              sql.NullString
	ticker                     sql.NullString
	shortName                  sql.NullString
	shares                     sql.NullString
	lastPrice                  sql.NullString
	marketValue                sql.NullString
	currency                   sql.NullString
	exchangeRate               sql.NullString
	lastPriceBaseCurrency      sql.NullString
	marketValueBaseCurrency    sql.NullString
	averagePrice               sql.NullString
	averagePriceBaseCurrency   sql.NullString
	gain                       sql.NullString
	percentageGain             sql.NullString
	gainBaseCurrency           sql.NullString
	percentageGainBaseCurrency sql.NullString
}

type PortfolioSummary struct {
	portfolioId           sql.NullString
	name                  sql.NullString
	currency              sql.NullString
	cacheValue            sql.NullString
	outgoings             sql.NullString
	incomings             sql.NullString
	gainOfSoldShares      sql.NullString
	commision             sql.NullString
	tax                   sql.NullString
	gainOfOwnedShares     sql.NullString
	estimatedGain         sql.NullString
	estimatedGainCostsInc sql.NullString
}

func Summary(db *sqlx.DB) {
	rows, err := db.Query("SELECT portfolio_id, portfolio_name, ticker, short_name, shares, last_price, market_value, currency, exchange_rate, last_price_base_currency, market_value_base_currency, average_price, average_price_base_currency, gain, percentage_gain, gain_base_currency, percentage_gain_base_currency  FROM owned_shares_summary")
	LogError(err)

	var data [][]string

	for rows.Next() {
		var oss OwnedSharesSummary

		err = rows.Scan(
			&oss.portfolioId,
			&oss.portfolioName,
			&oss.ticker,
			&oss.shortName,
			&oss.shares,
			&oss.lastPrice,
			&oss.marketValue,
			&oss.currency,
			&oss.exchangeRate,
			&oss.lastPriceBaseCurrency,
			&oss.marketValueBaseCurrency,
			&oss.averagePrice,
			&oss.averagePriceBaseCurrency,
			&oss.gain,
			&oss.percentageGain,
			&oss.gainBaseCurrency,
			&oss.percentageGainBaseCurrency,
		)
		LogError(err)

		var stock string
		if oss.shortName.String == "" {
			stock = strings.Trim(oss.ticker.String, " ")
		} else {
			stock = oss.shortName.String
		}

		data = append(data, []string{
			oss.portfolioName.String,
			stock,
			oss.shares.String,
			oss.lastPrice.String,
			oss.averagePrice.String,
			oss.lastPriceBaseCurrency.String,
			oss.averagePriceBaseCurrency.String,
			oss.gain.String,
			oss.gainBaseCurrency.String,
			oss.percentageGain.String,
			oss.percentageGainBaseCurrency.String,
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

	rows, err = db.Query("SELECT portfolio_id, name, currency, cache_value, outgoings, incomings, gain_of_sold_shares, commision, tax, gain_of_owned_shares, estimated_gain, estimated_gain_costs_inc FROM portfolios_summary")
	LogError(err)

	for rows.Next() {
		var ps PortfolioSummary

		err = rows.Scan(
			&ps.portfolioId,
			&ps.name,
			&ps.currency,
			&ps.cacheValue,
			&ps.outgoings,
			&ps.incomings,
			&ps.gainOfSoldShares,
			&ps.commision,
			&ps.tax,
			&ps.gainOfOwnedShares,
			&ps.estimatedGain,
			&ps.estimatedGainCostsInc,
		)
		LogError(err)

		data = append(data, []string{
			ps.portfolioId.String,
			ps.name.String,
			ps.cacheValue.String,
			ps.gainOfSoldShares.String,
			ps.commision.String,
			ps.tax.String,
			ps.gainOfOwnedShares.String,
			ps.estimatedGain.String,
			ps.estimatedGainCostsInc.String,
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
