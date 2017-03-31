package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/lib/pq"
	"github.com/olekukonko/tablewriter"
)

func Summary() {
	db := GetDb()
	defer db.Close()

	rows, err := db.Query("SELECT portfolio_id, portfolio_name, ticker, short_name, shares, last_price, market_value, currency, exchange_rate, last_price_base_currency, market_value_base_currency, average_price, average_price_base_currency, gain, percentage_gain, gain_base_currency, percentage_gain_base_currency  FROM owned_shares_summary")
	LogError(err)

	var data [][]string

	for rows.Next() {
		var portfolioId sql.NullString
		var portfolioName sql.NullString
		var ticker sql.NullString
		var shortName sql.NullString
		var shares sql.NullString
		var lastPrice sql.NullString
		var marketValue sql.NullString
		var currency sql.NullString
		var exchangeRate sql.NullString
		var lastPriceBaseCurrency sql.NullString
		var marketValueBaseCurrency sql.NullString
		var averagePrice sql.NullString
		var averagePriceBaseCurrency sql.NullString
		var gain sql.NullString
		var percentageGain sql.NullString
		var gainBaseCurrency sql.NullString
		var percentageGainBaseCurrency sql.NullString

		err = rows.Scan(
			&portfolioId,
			&portfolioName,
			&ticker,
			&shortName,
			&shares,
			&lastPrice,
			&marketValue,
			&currency,
			&exchangeRate,
			&lastPriceBaseCurrency,
			&marketValueBaseCurrency,
			&averagePrice,
			&averagePriceBaseCurrency,
			&gain,
			&percentageGain,
			&gainBaseCurrency,
			&percentageGainBaseCurrency,
		)
		LogError(err)

		data = append(data, []string{
			portfolioId.String,
			portfolioName.String,
			strings.Trim(ticker.String, " "),
			shortName.String,
			shares.String,
			lastPrice.String,
			marketValue.String,
			currency.String,
			exchangeRate.String,
			lastPriceBaseCurrency.String,
			marketValueBaseCurrency.String,
			averagePrice.String,
			averagePriceBaseCurrency.String,
			gain.String,
			percentageGain.String,
			gainBaseCurrency.String,
			percentageGainBaseCurrency.String,
		})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"Portfolio Id",
		"Portfolio Name",
		"Ticker",
		"Short Name",
		"Shares",
		"Last Price",
		"Market Value",
		"Currency",
		"Ex Rate",
		"Last Price BC",
		"MV BC",
		"Average Price",
		"Average Price BC",
		"Gain",
		"Gain%",
		"Gain BC",
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
		var portfolioId sql.NullString
		var name sql.NullString
		var currency sql.NullString
		var cacheValue sql.NullString
		var outgoings sql.NullString
		var incomings sql.NullString
		var gainOfSoldShares sql.NullString
		var commision sql.NullString
		var tax sql.NullString
		var gainOfOwnedShares sql.NullString
		var estimatedGain sql.NullString
		var estimatedGainCostsInc sql.NullString

		err = rows.Scan(
			&portfolioId,
			&name,
			&currency,
			&cacheValue,
			&outgoings,
			&incomings,
			&gainOfSoldShares,
			&commision,
			&tax,
			&gainOfOwnedShares,
			&estimatedGain,
			&estimatedGainCostsInc,
		)
		LogError(err)

		data = append(data, []string{
			portfolioId.String,
			name.String,
			currency.String,
			cacheValue.String,
			outgoings.String,
			incomings.String,
			gainOfSoldShares.String,
			commision.String,
			tax.String,
			gainOfOwnedShares.String,
			estimatedGain.String,
			estimatedGainCostsInc.String,
		})
	}

	table = tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"Portfolio Id",
		"Portfolio Name",
		"Currency",
		"Cache Value",
		"Outgoings",
		"Incomings",
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
