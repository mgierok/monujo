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

		var stock string
		if shortName.String == "" {
			stock = strings.Trim(ticker.String, " ")
		} else {
			stock = shortName.String
		}

		data = append(data, []string{
			portfolioName.String,
			stock,
			shares.String,
			lastPrice.String,
			averagePrice.String,
			lastPriceBaseCurrency.String,
			averagePriceBaseCurrency.String,
			gain.String,
			gainBaseCurrency.String,
			percentageGain.String,
			percentageGainBaseCurrency.String,
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
			cacheValue.String,
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
