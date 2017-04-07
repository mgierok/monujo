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
	get(portfolioId, &t)
	get(date, &t)
	get(ticker, &t)
	get(price, &t)
	get(transactionOperationType, &t)
	get(currency, &t)
	get(shares, &t)
	get(commision, &t)
	get(exchangeRate, &t)
	get(tax, &t)

	summary := [][]interface{}{
		[]interface{}{"Portfolio ID", t.PortfolioId},
		[]interface{}{"Date", t.Date},
		[]interface{}{"Ticker", t.Ticker},
		[]interface{}{"Price", t.Price},
		[]interface{}{"Type", t.TransactionOperationType},
		[]interface{}{"Currency", t.Currency},
		[]interface{}{"Shares", t.Shares},
		[]interface{}{"Commision", t.Commision},
		[]interface{}{"Exchange Rate", t.ExchangeRate},
		[]interface{}{"Tax", t.Tax},
	}

	Clear()
	DrawTable([]string{}, summary)

	transactionId, err := repository.StoreTransaction(t)
	LogError(err)

	fmt.Printf("Transaction has been recorded with an ID: %d\n", transactionId)
}

func get(f func(*entities.Transaction), t *entities.Transaction) {
	Clear()
	f(t)
}

func portfolioId(e *entities.Transaction) {
	fmt.Println("Choose portfolio")
	fmt.Println("")

	portfolios, err := repository.Portfolios()
	LogError(err)

	header := []string{
		"Portfolio Id",
		"Portfolio Name",
	}

	var dict = make(map[int64]string)
	var data [][]interface{}
	for _, p := range portfolios {
		data = append(data, []interface{}{p.PortfolioId, p.Name})
		dict[p.PortfolioId] = p.Name
	}

	DrawTable(header, data)
	fmt.Println("")

	var input string
	fmt.Print("Portfolio ID: ")
	fmt.Scanln(&input)

	p, err := strconv.ParseInt(input, 10, 64)

	if nil != err {
		fmt.Printf("\n%sd is not a valid portfolio ID\n\n", input)
		get(portfolioId, e)
		return
	} else {
		_, exists := dict[p]
		if exists {
			e.PortfolioId = p
		} else {
			fmt.Printf("\n%d is not a valid portfolio ID\n\n", p)
			get(portfolioId, e)
			return
		}
	}
}

func date(e *entities.Transaction) {
	const layout = "2006-01-02"
	var now = time.Now().Format(layout)
	var d string

	fmt.Printf("Date (default: %q): ", now)
	fmt.Scanln(&d)

	d = strings.Trim(d, " ")
	if d == "" {
		d = now
	} else {
		_, err := time.Parse(layout, d)
		if err != nil {
			fmt.Println(err)
			fmt.Printf("\n%q is not a valid date format\n\n", d)
			get(date, e)
			return
		}
	}

	e.Date = d
}

func ticker(e *entities.Transaction) {
	fmt.Print("Ticker: ")
	var t string
	fmt.Scanln(&t)

	t = strings.Trim(t, " ")
	if t == "" {
		get(ticker, e)
		return
	}

	e.Ticker = strings.ToUpper(t)
}

func price(e *entities.Transaction) {
	fmt.Print("Price: ")
	var input string
	fmt.Scanln(&input)

	p, err := strconv.ParseFloat(input, 64)

	if err != nil {
		fmt.Printf("\n%s is not a valid price value\n\n", input)
		get(price, e)
		return
	}

	e.Price = p
}

func transactionOperationType(e *entities.Transaction) {
	fmt.Println("Choose type of transaction")
	fmt.Println("")

	operationTypes, err := repository.TransactionalOperationTypes()
	LogError(err)

	header := []string{
		"Transaction Type",
	}

	var dict = make(map[string]string)
	var data [][]interface{}
	for _, ot := range operationTypes {
		dict[ot.Type] = ot.Type
		data = append(data, []interface{}{ot.Type})
	}

	DrawTable(header, data)
	fmt.Println("")

	var ot string
	fmt.Print("Transaction type: ")
	fmt.Scanln(&ot)

	_, exists := dict[ot]
	if exists {
		e.TransactionOperationType = ot
	} else {
		fmt.Printf("\n%s is not a valid transaction type\n\n", ot)
		get(transactionOperationType, e)
		return
	}
}

func currency(e *entities.Transaction) {
	fmt.Println("Choose currency")
	fmt.Println("")

	currencies, err := repository.Currencies()
	LogError(err)

	header := []string{
		"Currency",
	}

	var dict = make(map[string]string)
	var data [][]interface{}
	for _, c := range currencies {
		dict[c.Symbol] = c.Symbol
		data = append(data, []interface{}{c.Symbol})
	}

	DrawTable(header, data)
	fmt.Println("")

	var c string
	fmt.Print("Currency: ")
	fmt.Scanln(&c)

	c = strings.ToUpper(c)

	_, exists := dict[c]
	if exists {
		e.Currency = c
	} else {
		fmt.Printf("\n%s is not a valid currency\n\n", c)
		get(currency, e)
		return
	}
}

func shares(e *entities.Transaction) {
	fmt.Print("Number of shares: ")
	var input string
	fmt.Scanln(&input)

	s, err := strconv.ParseFloat(input, 64)

	if err != nil {
		fmt.Printf("\n%s is not a valid share number value\n\n", input)
		get(shares, e)
		return
	}

	e.Shares = s
}

func exchangeRate(e *entities.Transaction) {
	fmt.Print("Exchange rate (default: 1):")
	var input string
	fmt.Scanln(&input)

	er, err := strconv.ParseFloat(input, 64)

	if input == "" {
		er = 1
	} else {
		if err != nil {
			fmt.Printf("\n%s is not a valid exchange rate value\n\n", input)
			get(exchangeRate, e)
			return
		}
	}

	e.ExchangeRate = er
}

func commision(e *entities.Transaction) {
	fmt.Print("Commision (default: 0): ")
	var input string
	fmt.Scanln(&input)

	c, err := strconv.ParseFloat(input, 64)

	if input == "" {
		c = 0
	} else {
		if err != nil {
			fmt.Printf("\n%s is not a valid commision value\n\n", input)
			get(commision, e)
			return
		}
	}

	e.Commision = c
}

func tax(e *entities.Transaction) {
	fmt.Print("Tax (default: 0): ")
	var input string
	fmt.Scanln(&input)

	t, err := strconv.ParseFloat(input, 64)

	if input == "" {
		t = 0
	} else {
		if err != nil {
			fmt.Printf("\n%s is not a valid tax value\n\n", input)
			get(tax, e)
			return
		}
	}

	e.Tax = t
}
