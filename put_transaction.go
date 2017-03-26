package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"github.com/olekukonko/tablewriter"
)

func PutTransaction() {
	db := GetDb()
	defer db.Close()

	portfolioId := choosePortfolio(db)
	date := provideDateOfTransaction()
	ticker := provideTicker()
	price := providePrice()
	typeOfTransaction := chooseTypeOfTransaction(db)
	currency := chooseCurrency(db)
	numberOfShares := provideNumberOfShares()
	commision := provideCommision()
	exchangeRate := provideExchangeRate()
	tax := provideTax()

	table := tablewriter.NewWriter(os.Stdout)
	table.AppendBulk(
		[][]string{
			[]string{"Portfolio ID", strconv.FormatInt(portfolioId, 10)},
			[]string{"Date", date},
			[]string{"Ticker", ticker},
			[]string{"Price", strconv.FormatFloat(price, 'f', -1, 64)},
			[]string{"Type", typeOfTransaction},
			[]string{"Currency", currency},
			[]string{"Shares", strconv.FormatFloat(numberOfShares, 'f', -1, 64)},
			[]string{"Commision", strconv.FormatFloat(commision, 'f', -1, 64)},
			[]string{"Exchange Rate", strconv.FormatFloat(exchangeRate, 'f', -1, 64)},
			[]string{"Tax", strconv.FormatFloat(tax, 'f', -1, 64)},
		},
	)
	table.Render()

	var transactionId int
	err := db.QueryRow(`INSERT INTO transactions(portfolio_id, date, ticker, price, type, currency, shares, commision, exchange_rate, tax)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING transaction_id`, portfolioId, date, ticker, price, typeOfTransaction, currency, numberOfShares, commision, exchangeRate, tax).Scan(&transactionId)
	LogError(err)

	fmt.Printf("Transaction has been recorded with an ID: %d", transactionId)
}

func choosePortfolio(db *sql.DB) int64 {
	fmt.Println("Choose portfolio")
	fmt.Println("")

	rows, err := db.Query("SELECT portfolio_id, name FROM portfolios ORDER BY portfolio_id ASC")
	LogError(err)

	var portfolios = make(map[int64]string)

	for rows.Next() {
		var portfolioId sql.NullInt64
		var portfolioName sql.NullString
		err = rows.Scan(
			&portfolioId,
			&portfolioName,
		)
		LogError(err)

		portfolios[portfolioId.Int64] = portfolioName.String
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"Portfolio Id",
		"Portfolio Name",
	})

	for portfolioId, portfolioName := range portfolios {
		table.Append([]string{strconv.FormatInt(portfolioId, 10), portfolioName})
	}

	table.SetAutoMergeCells(true)
	table.SetRowLine(true)
	table.Render()
	fmt.Println("")

	var portfolioId int64
	fmt.Print("Portfolio ID: ")
	fmt.Scanln(&portfolioId)

	_, exists := portfolios[portfolioId]
	if exists {
		return portfolioId
	} else {
		fmt.Printf("\n%d is not a valid portfolio ID\n\n", portfolioId)
		return choosePortfolio(db)
	}
}

func chooseTypeOfTransaction(db *sql.DB) string {
	fmt.Println("Choose type of transaction")
	fmt.Println("")

	rows, err := db.Query("SELECT transaction_operation_type FROM transaction_operation_types ORDER BY transaction_operation_type ASC")
	LogError(err)

	var typeOfTransactions = make(map[string]string)

	for rows.Next() {
		var typeOfTransaction string

		err = rows.Scan(
			&typeOfTransaction,
		)
		LogError(err)

		typeOfTransactions[typeOfTransaction] = typeOfTransaction
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"Transaction Type",
	})

	for _, typeOfTransaction := range typeOfTransactions {
		table.Append([]string{typeOfTransaction})
	}

	table.SetAutoMergeCells(true)
	table.SetRowLine(true)
	table.Render()
	fmt.Println("")

	var typeOfTransaction string
	fmt.Print("Transaction type: ")
	fmt.Scanln(&typeOfTransaction)

	_, exists := typeOfTransactions[typeOfTransaction]
	if exists {
		return typeOfTransaction
	} else {
		fmt.Printf("\n%s is not a valid transaction type\n\n", typeOfTransaction)
		return chooseTypeOfTransaction(db)
	}
}

func provideDateOfTransaction() string {
	const layout = "2006-01-02"
	var now = time.Now().Format(layout)
	var date string

	fmt.Printf("Date (default: %q): ", now)
	fmt.Scanln(&date)

	if date == "" {
		date = now
	} else {
		_, err := time.Parse(layout, date)
		if err != nil {
			fmt.Println(err)
			fmt.Printf("\n%q is not a valid date format\n\n", date)
			date = provideDateOfTransaction()
		}
	}

	return date
}

func provideTicker() string {
	fmt.Print("Ticker: ")
	var ticker string
	fmt.Scanln(&ticker)

	if ticker == "" {
		return provideTicker()
	}

	return strings.ToUpper(ticker)
}

func providePrice() float64 {
	fmt.Print("Price: ")
	var price string
	fmt.Scanln(&price)

	floatPrice, err := strconv.ParseFloat(price, 64)

	if err != nil {
		fmt.Printf("\n%s is not a valid price value\n\n", price)
		return providePrice()
	}

	return floatPrice
}

func provideNumberOfShares() float64 {
	fmt.Print("Number of shares: ")
	var shares string
	fmt.Scanln(&shares)

	floatShares, err := strconv.ParseFloat(shares, 64)

	if err != nil {
		fmt.Printf("\n%s is not a valid share number value\n\n", shares)
		return provideNumberOfShares()
	}

	return floatShares
}

func provideExchangeRate() float64 {
	fmt.Print("Exchange rate (default: 1):")
	var exchangeRate string
	fmt.Scanln(&exchangeRate)

	floatExchangeRate, err := strconv.ParseFloat(exchangeRate, 64)

	if exchangeRate == "" {
		floatExchangeRate = 1
	} else {
		if err != nil {
			fmt.Printf("\n%s is not a valid exchange rate value\n\n", exchangeRate)
			floatExchangeRate = provideExchangeRate()
		}
	}

	return floatExchangeRate
}

func provideCommision() float64 {
	fmt.Print("Commision (default: 0): ")
	var commision string
	fmt.Scanln(&commision)

	floatCommision, err := strconv.ParseFloat(commision, 64)

	if commision == "" {
		floatCommision = 0
	} else {
		if err != nil {
			fmt.Printf("\n%s is not a valid commision value\n\n", commision)
			floatCommision = provideCommision()
		}
	}

	return floatCommision
}

func provideTax() float64 {
	fmt.Print("Tax (default: 0): ")
	var tax string
	fmt.Scanln(&tax)

	floatTax, err := strconv.ParseFloat(tax, 64)

	if tax == "" {
		floatTax = 0
	} else {
		if err != nil {
			fmt.Printf("\n%s is not a valid tax value\n\n", tax)
			floatTax = provideTax()
		}
	}

	return floatTax
}

func chooseCurrency(db *sql.DB) string {
	fmt.Println("Choose currency")
	fmt.Println("")

	rows, err := db.Query("SELECT currency FROM currencies ORDER BY currency ASC")
	LogError(err)

	var currencies = make(map[string]string)

	for rows.Next() {
		var currency string

		err = rows.Scan(
			&currency,
		)
		LogError(err)

		currencies[currency] = currency
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"Currency",
	})

	for _, currency := range currencies {
		table.Append([]string{currency})
	}

	table.SetAutoMergeCells(true)
	table.SetRowLine(true)
	table.Render()
	fmt.Println("")

	var currency string
	fmt.Print("Currency: ")
	fmt.Scanln(&currency)

	currency = strings.ToUpper(currency)

	_, exists := currencies[currency]
	if exists {
		return currency
	} else {
		fmt.Printf("\n%s is not a valid currency\n\n", currency)
		return chooseCurrency(db)
	}
}
