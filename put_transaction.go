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
	t.PortfolioId = portfolioId()
	t.Date = date()
	t.Ticker = ticker()
	t.Price = price()
	t.TransactionOperationType = transactionOperationType()
	t.Currency = currency()
	t.Shares = shares()
	t.Commision = commision()
	t.ExchangeRate = exchangeRate()
	t.Tax = tax()

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

func portfolioId() int64 {
	fmt.Println("Choose portfolio")
	fmt.Println("")

	portfolios, err := repository.Portfolios()
	LogError(err)

	header := []string{
		"Portfolio Id",
		"Portfolio Name",
	}

	var dict = make(map[int64]string)
	var data [][]string
	for _, p := range portfolios {
		data = append(data, []string{p.PortfolioId.String, p.Name.String})
		portfolioId, _ := strconv.ParseInt(p.PortfolioId.String, 10, 64)
		dict[portfolioId] = p.Name.String
	}

	DrawTable(header, data)
	fmt.Println("")

	var input string
	fmt.Print("Portfolio ID: ")
	fmt.Scanln(&input)

	p, err := strconv.ParseInt(input, 10, 64)

	if nil != err {
		fmt.Printf("\n%sd is not a valid portfolio ID\n\n", input)
		return portfolioId()
	}

	_, exists := dict[p]
	if exists {
		return p
	} else {
		fmt.Printf("\n%d is not a valid portfolio ID\n\n", p)
		return portfolioId()
	}
}

func date() string {
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
			d = date()
		}
	}

	return d
}

func ticker() string {
	fmt.Print("Ticker: ")
	var t string
	fmt.Scanln(&t)

	t = strings.Trim(t, " ")
	if t == "" {
		return ticker()
	}

	return strings.ToUpper(t)
}

func price() float64 {
	fmt.Print("Price: ")
	var input string
	fmt.Scanln(&input)

	p, err := strconv.ParseFloat(input, 64)

	if err != nil {
		fmt.Printf("\n%s is not a valid price value\n\n", input)
		return price()
	}

	return p
}

func transactionOperationType() string {
	fmt.Println("Choose type of transaction")
	fmt.Println("")

	operationTypes, err := repository.TransactionalOperationTypes()
	LogError(err)

	header := []string{
		"Transaction Type",
	}

	var dict = make(map[string]string)
	var data [][]string
	for _, ot := range operationTypes {
		dict[ot.Type] = ot.Type
		data = append(data, []string{ot.Type})
	}

	DrawTable(header, data)
	fmt.Println("")

	var ot string
	fmt.Print("Transaction type: ")
	fmt.Scanln(&ot)

	_, exists := dict[ot]
	if exists {
		return ot
	} else {
		fmt.Printf("\n%s is not a valid transaction type\n\n", ot)
		return transactionOperationType()
	}
}

func currency() string {
	fmt.Println("Choose currency")
	fmt.Println("")

	currencies, err := repository.Currencies()
	LogError(err)

	header := []string{
		"Currency",
	}

	var dict = make(map[string]string)
	var data [][]string
	for _, c := range currencies {
		dict[c.Symbol] = c.Symbol
		data = append(data, []string{c.Symbol})
	}

	DrawTable(header, data)
	fmt.Println("")

	var c string
	fmt.Print("Currency: ")
	fmt.Scanln(&c)

	c = strings.ToUpper(c)

	_, exists := dict[c]
	if exists {
		return c
	} else {
		fmt.Printf("\n%s is not a valid currency\n\n", c)
		return currency()
	}
}

func shares() float64 {
	fmt.Print("Number of shares: ")
	var input string
	fmt.Scanln(&input)

	s, err := strconv.ParseFloat(input, 64)

	if err != nil {
		fmt.Printf("\n%s is not a valid share number value\n\n", input)
		return shares()
	}

	return s
}

func exchangeRate() float64 {
	fmt.Print("Exchange rate (default: 1):")
	var input string
	fmt.Scanln(&input)

	er, err := strconv.ParseFloat(input, 64)

	if input == "" {
		er = 1
	} else {
		if err != nil {
			fmt.Printf("\n%s is not a valid exchange rate value\n\n", input)
			er = exchangeRate()
		}
	}

	return er
}

func commision() float64 {
	fmt.Print("Commision (default: 0): ")
	var input string
	fmt.Scanln(&input)

	c, err := strconv.ParseFloat(input, 64)

	if input == "" {
		c = 0
	} else {
		if err != nil {
			fmt.Printf("\n%s is not a valid commision value\n\n", input)
			c = commision()
		}
	}

	return c
}

func tax() float64 {
	fmt.Print("Tax (default: 0): ")
	var input string
	fmt.Scanln(&input)

	t, err := strconv.ParseFloat(input, 64)

	if input == "" {
		t = 0
	} else {
		if err != nil {
			fmt.Printf("\n%s is not a valid tax value\n\n", input)
			t = tax()
		}
	}

	return t
}
