package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/chrisftw/ezconf"
	_ "github.com/lib/pq"
	"github.com/olekukonko/tablewriter"
)

func main() {
	dbinfo := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s sslmode=disable",
		ezconf.Get("db", "host"),
		ezconf.Get("db", "user"),
		ezconf.Get("db", "password"),
		ezconf.Get("db", "dbname"),
	)
	db, err := sql.Open("postgres", dbinfo)

	if err != nil {
		log.Fatal("Error: The data source arguments are not valid")
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Error: Could not establish a connection with the database")
	}

	rows, err := db.Query("SELECT portfolio_id, portfolio_name, ticker, short_name, shares, last_price, market_value, currency, exchange_rate, market_value_base_currency, average_price, gain, percentage_gain, gain_base_currency, percentage_gain_base_currency  FROM owned_shares_summary")
	logError(err)

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
		var marketValueBaseCurrency sql.NullString
		var averagePrice sql.NullString
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
			&marketValueBaseCurrency,
			&averagePrice,
			&gain,
			&percentageGain,
			&gainBaseCurrency,
			&percentageGainBaseCurrency,
		)
		logError(err)

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
			marketValueBaseCurrency.String,
			averagePrice.String,
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
		"MV BC",
		"Average Price",
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

	rows, err = db.Query("SELECT portfolio_id, name, currency, cache_value, outgoings, incomings, gain_of_sold_shares, commision, tax, gain_of_owned_shares FROM portfolios_summary")
	logError(err)

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
		)
		logError(err)

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
	})

	table.AppendBulk(data)
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)
	table.Render()
}

func logError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
