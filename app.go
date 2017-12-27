package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/mgierok/monujo/log"
)

type Screen interface {
	DrawTable(header []string, data [][]interface{})
	Clear()
}

type Input interface {
	String(name string, args ...string) string
	Float(name string, args ...float64) float64
	Date(name string, args ...time.Time) time.Time
}

type app struct {
	config     Config
	screen     Screen
	input      Input
	repository Repository
}

func NewApp(c *Config, r *Repository, s Screen, i Input) (*app, error) {
	a := new(app)
	a.config = *c
	a.screen = s
	a.input = i
	a.repository = *r

	return a, nil
}

func (a *app) Run() {
	a.mainMenu()
}

func (a *app) mainMenu() {
	fmt.Println("Choose action")
	data := [][]interface{}{
		[]interface{}{"S", "Summary"},
		[]interface{}{"PT", "Put transaction"},
		[]interface{}{"LT", "List transactions"},
		[]interface{}{"PO", "Put operation"},
		[]interface{}{"LO", "List operations"},
		[]interface{}{"U", "Update Quotes"},
		[]interface{}{"Q", "Quit"},
	}

	a.screen.DrawTable([]string{}, data)

	var action string
	fmt.Scanln(&action)
	action = strings.ToUpper(action)
	a.screen.Clear()

	if action == "S" {
		a.summary()
	} else if action == "PT" {
		a.putTransaction()
	} else if action == "LT" {
		a.listTransactions()
	} else if action == "PO" {
		a.putOperation()
	} else if action == "LO" {
		a.listOperations()
	} else if action == "U" {
		a.update()
	} else if action == "Q" {
		return
	}

	var input string
	fmt.Scanln(&input)

	a.screen.Clear()
	a.mainMenu()
}

func (a *app) summary() {
	ownedStocks, err := a.repository.OwnedStocks()
	log.PanicIfError(err)

	var data [][]interface{}

	for _, os := range ownedStocks {
		data = append(data, []interface{}{
			os.PortfolioName,
			os.DisplayName(),
			os.Shares,
			os.LastPrice,
			os.AveragePrice,
			os.InvestmentBaseCurrency,
			os.MarketValueBaseCurrency,
			os.GainBaseCurrency,
			os.PercentageGainBaseCurrency,
		})
	}

	header := []string{
		"Portfolio Name",
		"Stock",
		"Shares",
		"Last Price",
		"Average Price",
		"Investment BC",
		"Market Value BC",
		"Gain BC",
		"Gain BC%",
	}

	a.screen.DrawTable(header, data)

	data = data[0:0]
	fmt.Println("")
	fmt.Println("")

	portfoliosExt, err := a.repository.PortfoliosExt()
	log.PanicIfError(err)

	for _, pe := range portfoliosExt {
		data = append(data, []interface{}{
			pe.Name,
			pe.CacheValue,
			pe.GainOfSoldShares,
			pe.GainOfOwnedShares,
			pe.EstimatedGain,
			pe.EstimatedGainCostsInc,
			pe.EstimatedValue,
			pe.AnnualBalance,
			pe.MonthBalance,
		})
	}

	header = []string{
		"Portfolio Name",
		"Cache Value",
		"Gain of Sold Shares",
		"Gain Of Ownded Shares",
		"Estimated Gain",
		"Estimated Gain Costs Inc.",
		"Estimated Value",
		"Annual Balance",
		"Month Balance",
	}

	a.screen.DrawTable(header, data)
}

func (a *app) listTransactions() {
	portfolio := a.portfolio()

	transactions, err := a.repository.PortfolioTransactions(portfolio)
	log.PanicIfError(err)

	var data [][]interface{}

	for _, t := range transactions {
		data = append(data, []interface{}{
			t.TransactionId,
			t.PortfolioId,
			t.Date,
			t.Ticker,
			t.Price,
			t.Currency,
			t.Shares,
			t.Commision,
			t.ExchangeRate,
			t.Tax,
		})
	}

	header := []string{
		"Transaction ID",
		"Portfolio ID",
		"Date",
		"Ticker",
		"Price",
		"Currency",
		"Shares",
		"Commision",
		"Exchange Rate",
		"Tax",
	}

	a.screen.DrawTable(header, data)
	fmt.Println("")

	if !a.yesOrNo("Do you want to delete single transaction?") {
		return
	}

	transaction := a.pickTransaction(transactions)
	err = a.repository.DeleteTransaction(transaction)
	log.PanicIfError(err)
	fmt.Println("Transaction has been removed")
}

func (a *app) listOperations() {
	portfolio := a.portfolio()

	operations, err := a.repository.PortfolioOperations(portfolio)
	log.PanicIfError(err)

	var data [][]interface{}

	for _, o := range operations {
		data = append(data, []interface{}{
			o.OperationId,
			o.PortfolioId,
			o.Date,
			o.Type,
			o.Value,
			o.Description,
			o.Commision,
		})
	}

	header := []string{
		"Operation ID",
		"Portfolio ID",
		"Date",
		"Type",
		"Value",
		"Description",
		"Commision",
	}

	a.screen.DrawTable(header, data)
	fmt.Println("")

	if !a.yesOrNo("Do you want to delete single financial operation?") {
		return
	}

	operation := a.pickOperation(operations)
	err = a.repository.DeleteOperation(operation)
	log.PanicIfError(err)
	fmt.Println("Operation has been removed")
}

func (a *app) update() {
	sources := a.pickSource()
	quotes, err := a.repository.UpdateQuotes(sources)
	log.PanicIfError(err)

	for q := range quotes {
		_, err = a.repository.StoreLatestQuote(q)
		if err == nil {
			fmt.Printf("Ticker: %s Quote: %f\n", q.Ticker, q.Close)
		} else {
			fmt.Printf("Update failed for %s\n", q.Ticker)
		}
	}
}

func (a *app) putOperation() {
	var o Operation
	o.PortfolioId = a.portfolio().PortfolioId
	a.screen.Clear()
	o.Date = a.input.Date("Date", time.Now())
	a.screen.Clear()
	o.Type = a.financialOperationType().Type
	a.screen.Clear()
	o.Value = a.input.Float("Value")
	a.screen.Clear()
	o.Description = a.input.String("Description", "")
	a.screen.Clear()
	o.Commision = a.input.Float("Commision", 0)
	a.screen.Clear()
	o.Tax = a.input.Float("Tax", 0)
	a.screen.Clear()

	summary := [][]interface{}{
		[]interface{}{"Portfolio ID", o.PortfolioId},
		[]interface{}{"Date", o.Date},
		[]interface{}{"Operation type", o.Type},
		[]interface{}{"Value", o.Value},
		[]interface{}{"Description", o.Description},
		[]interface{}{"Commision", o.Commision},
		[]interface{}{"Tax", o.Tax},
	}

	a.screen.Clear()
	a.screen.DrawTable([]string{}, summary)
	fmt.Println("")

	if a.yesOrNo("Do you want to store this operation?") {
		operationId, err := a.repository.StoreOperation(o)
		log.PanicIfError(err)

		fmt.Printf("Operation has been recorded with an ID: %d\n", operationId)
	} else {
		fmt.Println("Operation has not been recorded")
	}

}

func (a *app) putTransaction() {
	var t Transaction
	t.PortfolioId = a.portfolio().PortfolioId
	a.screen.Clear()
	t.Date = a.input.Date("Date", time.Now())
	a.screen.Clear()
	t.Ticker = a.input.String("Ticker")
	a.screen.Clear()
	t.Price = a.input.Float("Price")
	a.screen.Clear()
	t.Currency = a.pickCurrency()
	a.screen.Clear()
	t.Shares = a.shares()
	a.screen.Clear()
	t.Commision = a.input.Float("Commision", 0)
	a.screen.Clear()
	t.ExchangeRate = a.input.Float("Exchange rate", 1)
	a.screen.Clear()
	t.Tax = a.input.Float("Tax", 0)
	a.screen.Clear()

	summary := [][]interface{}{
		[]interface{}{"Portfolio ID", t.PortfolioId},
		[]interface{}{"Date", t.Date},
		[]interface{}{"Ticker", t.Ticker},
		[]interface{}{"Price", t.Price},
		[]interface{}{"Currency", t.Currency},
		[]interface{}{"Shares", t.Shares},
		[]interface{}{"Commision", t.Commision},
		[]interface{}{"Exchange Rate", t.ExchangeRate},
		[]interface{}{"Tax", t.Tax},
	}

	a.screen.Clear()
	a.screen.DrawTable([]string{}, summary)
	fmt.Println("")

	if a.yesOrNo("Do you want to store this transaction?") {
		transactionId, err := a.repository.StoreTransaction(t)
		log.PanicIfError(err)

		fmt.Printf("Transaction has been recorded with an ID: %d\n", transactionId)

		a.securityDetails(t.Ticker)
	} else {
		fmt.Println("Transaction has not been recorded")
	}
}

func (a *app) pickTransaction(transactions Transactions) Transaction {
	var input string
	fmt.Print("Transaction ID: ")
	fmt.Scanln(&input)

	transactionId, err := strconv.ParseInt(input, 10, 64)

	if nil != err {
		fmt.Printf("\n%s is not a valid transaction ID\n\n", input)
		return a.pickTransaction(transactions)
	} else {
		for _, t := range transactions {
			if t.TransactionId == transactionId {
				return t
			}
		}

		fmt.Printf("\n%s is not a valid transaction ID\n\n", input)
		return a.pickTransaction(transactions)
	}
}

func (a *app) pickOperation(operations Operations) Operation {
	var input string
	fmt.Print("Operation ID: ")
	fmt.Scanln(&input)

	operationId, err := strconv.ParseInt(input, 10, 64)

	if nil != err {
		fmt.Printf("\n%s is not a valid operation ID\n\n", input)
		return a.pickOperation(operations)
	} else {
		for _, o := range operations {
			if o.OperationId == operationId {
				return o
			}
		}

		fmt.Printf("\n%s is not a valid operation ID\n\n", input)
		return a.pickOperation(operations)
	}
}

func (a *app) yesOrNo(question string) bool {
	fmt.Println(question)
	fmt.Println("(Y)es or (N)o?")

	var input string
	fmt.Scanln(&input)
	input = strings.ToUpper(input)

	if "Y" == input {
		return true
	} else if "N" == input {
		return false
	}

	return a.yesOrNo(question)
}

func (a *app) portfolio() Portfolio {
	fmt.Println("Choose portfolio")
	fmt.Println("")

	portfolios, err := a.repository.Portfolios()
	log.PanicIfError(err)

	header := []string{
		"Portfolio Id",
		"Portfolio Name",
	}

	var data [][]interface{}
	for _, p := range portfolios {
		data = append(data, []interface{}{p.PortfolioId, p.Name})
	}

	a.screen.DrawTable(header, data)
	fmt.Println("")

	var input string
	fmt.Print("Portfolio ID: ")
	fmt.Scanln(&input)

	portfolioId, err := strconv.ParseInt(input, 10, 64)

	if nil != err {
		fmt.Printf("\n%s is not a valid portfolio ID\n\n", input)
		return a.portfolio()
	} else {
		for _, p := range portfolios {
			if p.PortfolioId == portfolioId {
				return p
			}
		}

		fmt.Printf("\n%s is not a valid portfolio ID\n\n", input)
		return a.portfolio()
	}
}

func (a *app) pickSource() Sources {
	fmt.Println("Choose from which source you want to update quotes")
	fmt.Println("")

	dict := map[string]string{
		"A": "All",
		"Q": "Quit",
	}
	data := [][]interface{}{
		[]interface{}{"A", "All"},
	}
	i := 1

	sources, err := a.repository.Sources()
	log.PanicIfError(err)

	for _, s := range sources {
		dict[strconv.Itoa(i)] = s.Name
		data = append(data, []interface{}{strconv.Itoa(i), s.Name})
		i++
	}
	data = append(data, []interface{}{"Q", "Quit"})

	a.screen.DrawTable([]string{}, data)
	fmt.Println("")

	var input string
	fmt.Scanln(&input)
	a.screen.Clear()

	input = strings.ToUpper(input)

	_, exists := dict[input]
	if exists {
		if input == "A" {
			return sources
		} else if input == "Q" {
			return Sources{}
		} else {
			for _, s := range sources {
				if s.Name == dict[input] {
					return Sources{s}
				}
			}
		}
	}
	return a.pickSource()
}

func (a *app) financialOperationType() FinancialOperationType {
	fmt.Println("Choose operation type")
	fmt.Println("")

	ots, err := a.repository.FinancialOperationTypes()
	log.PanicIfError(err)

	header := []string{
		"Operation type",
	}

	var dict = make(map[string]FinancialOperationType)
	var data [][]interface{}
	for _, ot := range ots {
		dict[ot.Type] = ot
		data = append(data, []interface{}{ot.Type})
	}

	a.screen.DrawTable(header, data)
	fmt.Println("")

	fmt.Print("Type: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	ot := scanner.Text()

	ot = strings.TrimSpace(ot)
	ot = strings.ToLower(ot)

	_, exists := dict[ot]
	if exists {
		return dict[ot]
	} else {
		fmt.Printf("\n%s is not a valid operation type\n\n", ot)
		return a.financialOperationType()
	}
}

func (a *app) securityDetails(ticker string) {
	exists, err := a.repository.SecurityExists(ticker)
	log.PanicIfError(err)
	if exists {
		return
	}

	if !a.yesOrNo(fmt.Sprintf("Would you like to add %s security detials to the database?", strings.TrimSpace(ticker))) {
		return
	}

	s := Security{
		Ticker: ticker,
	}
	s.ShortName = a.input.String("Short name")
	s.FullName = a.input.String("Full name")
	s.Market = a.input.String("Market")
	s.Leverage = a.input.Float("Leverage", 1)
	s.QuotesSource = a.input.String("Quotes source")
	tb := a.input.String("Ticker Bankier", "")
	s.TickerBankier = sql.NullString{String: tb, Valid: true}

	t, err := a.repository.StoreSecurity(s)
	log.PanicIfError(err)

	fmt.Printf("Security details of %s has been stored\n", strings.TrimSpace(t))
}

func (a *app) pickCurrency() string {
	fmt.Println("Choose currency")
	fmt.Println("")

	currencies, err := a.repository.Currencies()
	log.PanicIfError(err)

	header := []string{
		"Currency",
	}

	var dict = make(map[string]string)
	var data [][]interface{}
	for _, c := range currencies {
		dict[c.Symbol] = c.Symbol
		data = append(data, []interface{}{c.Symbol})
	}

	a.screen.DrawTable(header, data)
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
		return a.pickCurrency()
	}
}

func (a *app) shares() float64 {
	s := a.input.Float("Shares")

	isShort := a.isShort()
	if (isShort && s > 0) || (!isShort && s < 0) {
		return 0 - s
	}

	return s
}

func (a *app) isShort() bool {
	var input string
	fmt.Println("(B)UY or (S)ELL?")
	fmt.Scanln(&input)
	input = strings.ToUpper(input)

	if "S" == input {
		return true
	} else if "B" == input {
		return false
	}

	return a.isShort()
}

func (a *app) Dump(dumptype string, file string) {
	if len(file) == 0 {
		fmt.Println("Output file is not set")
		return
	}

	var cmd *exec.Cmd
	if dumptype == "schema" {
		cmd = exec.Command(
			a.config.Sys.Pgdump,
			"--host",
			a.config.Db.Host,
			"--port",
			a.config.Db.Port,
			"--username",
			a.config.Db.User,
			"--no-password",
			"--format",
			"plain",
			"--schema-only",
			"--no-owner",
			"--no-privileges",
			"--no-tablespaces",
			"--no-unlogged-table-data",
			"--file",
			file,
			a.config.Db.Dbname,
		)
	} else if dumptype == "data" {
		cmd = exec.Command(
			a.config.Sys.Pgdump,
			"--host",
			a.config.Db.Host,
			"--port",
			a.config.Db.Port,
			"--username",
			a.config.Db.User,
			"--no-password",
			"--format",
			"plain",
			"--data-only",
			"--inserts",
			"--disable-triggers",
			"--no-owner",
			"--no-privileges",
			"--no-tablespaces",
			"--no-unlogged-table-data",
			"--file",
			file,
			a.config.Db.Dbname,
		)
	} else {
		fmt.Println("Invalid dump type, please specify 'schema' or 'data'")
		return
	}

	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(string(stdout))
}
