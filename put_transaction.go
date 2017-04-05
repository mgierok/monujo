package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/mgierok/monujo/repository"
	"github.com/mgierok/monujo/repository/entities"
)

func PutTransaction() {

	var t entities.Transaction
	t.PortfolioId = choosePortfolio()
	t.Date = provideDateOfTransaction()
	t.Ticker = provideTicker()
	t.Price = providePrice()
	t.TransactionOperationType = chooseTypeOfTransaction()
	t.Currency = chooseCurrency()
	t.Shares = provideNumberOfShares()
	t.Commision = provideCommision()
	t.ExchangeRate = provideExchangeRate()
	t.Tax = provideTax()

	summary := [][]string{
		[]string{"Portfolio ID", strconv.FormatInt(t.PortfolioId, 10)},
		[]string{"Date", t.Date},
		[]string{"Ticker", t.Ticker},
		[]string{"Price", strconv.FormatFloat(t.Price, 'f', -1, 64)},
		[]string{"Type", t.TransactionOperationType},
		[]string{"Currency", t.Currency},
		[]string{"Shares", strconv.FormatFloat(t.Shares, 'f', -1, 64)},
		[]string{"Commision", strconv.FormatFloat(t.Commision, 'f', -1, 64)},
		[]string{"Exchange Rate", strconv.FormatFloat(t.ExchangeRate, 'f', -1, 64)},
		[]string{"Tax", strconv.FormatFloat(t.Tax, 'f', -1, 64)},
	}
	DrawTable([]string{}, summary)

	transactionId, err := repository.StoreTransaction(t)
	LogError(err)

	fmt.Printf("Transaction has been recorded with an ID: %d\n", transactionId)
}

func choosePortfolio() int64 {
	fmt.Println("Choose portfolio")
	fmt.Println("")

	portfolios, err := repository.Portfolios()
	LogError(err)

	header := []string{
		"Portfolio Id",
		"Portfolio Name",
	}

	var portfoliosDict = make(map[int64]string)
	var data [][]string
	for _, p := range portfolios {
		data = append(data, []string{p.PortfolioId.String, p.Name.String})
		portfolioId, _ := strconv.ParseInt(p.PortfolioId.String, 10, 64)
		portfoliosDict[portfolioId] = p.Name.String
	}

	DrawTable(header, data)
	fmt.Println("")

	var portfolioIdString string
	fmt.Print("Portfolio ID: ")
	fmt.Scanln(&portfolioIdString)

	portfolioId, err := strconv.ParseInt(portfolioIdString, 10, 64)

	if nil != err {
		fmt.Printf("\n%sd is not a valid portfolio ID\n\n", portfolioIdString)
		return choosePortfolio()
	}

	_, exists := portfoliosDict[portfolioId]
	if exists {
		return portfolioId
	} else {
		fmt.Printf("\n%d is not a valid portfolio ID\n\n", portfolioId)
		return choosePortfolio()
	}
}

func chooseTypeOfTransaction() string {
	fmt.Println("Choose type of transaction")
	fmt.Println("")

	operationTypes, err := repository.TransactionalOperationTypes()
	LogError(err)

	header := []string{
		"Transaction Type",
	}

	var operationTypesDict = make(map[string]string)
	var data [][]string
	for _, ot := range operationTypes {
		operationTypesDict[ot.Type.String] = ot.Type.String
		data = append(data, []string{ot.Type.String})
	}

	DrawTable(header, data)
	fmt.Println("")

	var operationType string
	fmt.Print("Transaction type: ")
	fmt.Scanln(&operationType)

	_, exists := operationTypesDict[operationType]
	if exists {
		return operationType
	} else {
		fmt.Printf("\n%s is not a valid transaction type\n\n", operationType)
		return chooseTypeOfTransaction()
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

func chooseCurrency() string {
	fmt.Println("Choose currency")
	fmt.Println("")

	currencies, err := repository.Currencies()
	LogError(err)

	header := []string{
		"Currency",
	}

	var currenciesDict = make(map[string]string)
	var data [][]string
	for _, c := range currencies {
		currenciesDict[c.Symbol] = c.Symbol
		data = append(data, []string{c.Symbol})
	}

	DrawTable(header, data)
	fmt.Println("")

	var currency string
	fmt.Print("Currency: ")
	fmt.Scanln(&currency)

	currency = strings.ToUpper(currency)

	_, exists := currenciesDict[currency]
	if exists {
		return currency
	} else {
		fmt.Printf("\n%s is not a valid currency\n\n", currency)
		return chooseCurrency()
	}
}
