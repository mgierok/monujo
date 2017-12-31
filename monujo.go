package monujo

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

type monujo struct {
	config     Config
	screen     Screen
	input      Input
	repository Repository
}

func New(c Config, r *Repository, s Screen, i Input) (*monujo, error) {
	m := new(monujo)
	m.config = c
	m.screen = s
	m.input = i
	m.repository = *r

	return m, nil
}

func (m *monujo) Run() {
	m.mainMenu()
}

func (m *monujo) mainMenu() {
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

	m.screen.DrawTable([]string{}, data)

	var action string
	fmt.Scanln(&action)
	action = strings.ToUpper(action)
	m.screen.Clear()

	if action == "S" {
		m.summary()
	} else if action == "PT" {
		m.putTransaction()
	} else if action == "LT" {
		m.listTransactions()
	} else if action == "PO" {
		m.putOperation()
	} else if action == "LO" {
		m.listOperations()
	} else if action == "U" {
		m.update()
	} else if action == "Q" {
		return
	}

	var input string
	fmt.Scanln(&input)

	m.screen.Clear()
	m.mainMenu()
}

func (m *monujo) summary() {
	ownedStocks, err := m.repository.OwnedStocks()
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

	m.screen.DrawTable(header, data)

	data = data[0:0]
	fmt.Println("")
	fmt.Println("")

	portfoliosExt, err := m.repository.PortfoliosExt()
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

	m.screen.DrawTable(header, data)
}

func (m *monujo) listTransactions() {
	portfolio := m.portfolio()

	transactions, err := m.repository.PortfolioTransactions(portfolio)
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

	m.screen.DrawTable(header, data)
	fmt.Println("")

	if !m.yesOrNo("Do you want to delete single transaction?") {
		return
	}

	transaction := m.pickTransaction(transactions)
	err = m.repository.DeleteTransaction(transaction)
	log.PanicIfError(err)
	fmt.Println("Transaction has been removed")
}

func (m *monujo) listOperations() {
	portfolio := m.portfolio()

	operations, err := m.repository.PortfolioOperations(portfolio)
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

	m.screen.DrawTable(header, data)
	fmt.Println("")

	if !m.yesOrNo("Do you want to delete single financial operation?") {
		return
	}

	operation := m.pickOperation(operations)
	err = m.repository.DeleteOperation(operation)
	log.PanicIfError(err)
	fmt.Println("Operation has been removed")
}

func (m *monujo) update() {
	sources := m.pickSource()
	quotes, err := m.repository.UpdateQuotes(sources)
	log.PanicIfError(err)

	for q := range quotes {
		_, err = m.repository.StoreLatestQuote(q)
		if err == nil {
			fmt.Printf("Ticker: %s Quote: %f\n", q.Ticker, q.Close)
		} else {
			fmt.Printf("Update failed for %s\n", q.Ticker)
		}
	}
}

func (m *monujo) putOperation() {
	var o Operation
	o.PortfolioId = m.portfolio().PortfolioId
	m.screen.Clear()
	o.Date = m.input.Date("Date", time.Now())
	m.screen.Clear()
	o.Type = m.financialOperationType().Type
	m.screen.Clear()
	o.Value = m.input.Float("Value")
	m.screen.Clear()
	o.Description = m.input.String("Description", "")
	m.screen.Clear()
	o.Commision = m.input.Float("Commision", 0)
	m.screen.Clear()
	o.Tax = m.input.Float("Tax", 0)
	m.screen.Clear()

	summary := [][]interface{}{
		[]interface{}{"Portfolio ID", o.PortfolioId},
		[]interface{}{"Date", o.Date},
		[]interface{}{"Operation type", o.Type},
		[]interface{}{"Value", o.Value},
		[]interface{}{"Description", o.Description},
		[]interface{}{"Commision", o.Commision},
		[]interface{}{"Tax", o.Tax},
	}

	m.screen.Clear()
	m.screen.DrawTable([]string{}, summary)
	fmt.Println("")

	if m.yesOrNo("Do you want to store this operation?") {
		operationId, err := m.repository.StoreOperation(o)
		log.PanicIfError(err)

		fmt.Printf("Operation has been recorded with an ID: %d\n", operationId)
	} else {
		fmt.Println("Operation has not been recorded")
	}

}

func (m *monujo) putTransaction() {
	var t Transaction
	t.PortfolioId = m.portfolio().PortfolioId
	m.screen.Clear()
	t.Date = m.input.Date("Date", time.Now())
	m.screen.Clear()
	t.Ticker = m.input.String("Ticker")
	m.screen.Clear()
	t.Price = m.input.Float("Price")
	m.screen.Clear()
	t.Currency = m.pickCurrency()
	m.screen.Clear()
	t.Shares = m.shares()
	m.screen.Clear()
	t.Commision = m.input.Float("Commision", 0)
	m.screen.Clear()
	t.ExchangeRate = m.input.Float("Exchange rate", 1)
	m.screen.Clear()
	t.Tax = m.input.Float("Tax", 0)
	m.screen.Clear()

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

	m.screen.Clear()
	m.screen.DrawTable([]string{}, summary)
	fmt.Println("")

	if m.yesOrNo("Do you want to store this transaction?") {
		transactionId, err := m.repository.StoreTransaction(t)
		log.PanicIfError(err)

		fmt.Printf("Transaction has been recorded with an ID: %d\n", transactionId)

		m.securityDetails(t.Ticker)
	} else {
		fmt.Println("Transaction has not been recorded")
	}
}

func (m *monujo) pickTransaction(transactions Transactions) Transaction {
	var input string
	fmt.Print("Transaction ID: ")
	fmt.Scanln(&input)

	transactionId, err := strconv.ParseInt(input, 10, 64)

	if nil != err {
		fmt.Printf("\n%s is not a valid transaction ID\n\n", input)
		return m.pickTransaction(transactions)
	} else {
		for _, t := range transactions {
			if t.TransactionId == transactionId {
				return t
			}
		}

		fmt.Printf("\n%s is not a valid transaction ID\n\n", input)
		return m.pickTransaction(transactions)
	}
}

func (m *monujo) pickOperation(operations Operations) Operation {
	var input string
	fmt.Print("Operation ID: ")
	fmt.Scanln(&input)

	operationId, err := strconv.ParseInt(input, 10, 64)

	if nil != err {
		fmt.Printf("\n%s is not a valid operation ID\n\n", input)
		return m.pickOperation(operations)
	} else {
		for _, o := range operations {
			if o.OperationId == operationId {
				return o
			}
		}

		fmt.Printf("\n%s is not a valid operation ID\n\n", input)
		return m.pickOperation(operations)
	}
}

func (m *monujo) yesOrNo(question string) bool {
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

	return m.yesOrNo(question)
}

func (m *monujo) portfolio() Portfolio {
	fmt.Println("Choose portfolio")
	fmt.Println("")

	portfolios, err := m.repository.Portfolios()
	log.PanicIfError(err)

	header := []string{
		"Portfolio Id",
		"Portfolio Name",
	}

	var data [][]interface{}
	for _, p := range portfolios {
		data = append(data, []interface{}{p.PortfolioId, p.Name})
	}

	m.screen.DrawTable(header, data)
	fmt.Println("")

	var input string
	fmt.Print("Portfolio ID: ")
	fmt.Scanln(&input)

	portfolioId, err := strconv.ParseInt(input, 10, 64)

	if nil != err {
		fmt.Printf("\n%s is not a valid portfolio ID\n\n", input)
		return m.portfolio()
	} else {
		for _, p := range portfolios {
			if p.PortfolioId == portfolioId {
				return p
			}
		}

		fmt.Printf("\n%s is not a valid portfolio ID\n\n", input)
		return m.portfolio()
	}
}

func (m *monujo) pickSource() Sources {
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

	sources, err := m.repository.Sources()
	log.PanicIfError(err)

	for _, s := range sources {
		dict[strconv.Itoa(i)] = s.Name
		data = append(data, []interface{}{strconv.Itoa(i), s.Name})
		i++
	}
	data = append(data, []interface{}{"Q", "Quit"})

	m.screen.DrawTable([]string{}, data)
	fmt.Println("")

	var input string
	fmt.Scanln(&input)
	m.screen.Clear()

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
	return m.pickSource()
}

func (m *monujo) financialOperationType() FinancialOperationType {
	fmt.Println("Choose operation type")
	fmt.Println("")

	ots, err := m.repository.FinancialOperationTypes()
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

	m.screen.DrawTable(header, data)
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
		return m.financialOperationType()
	}
}

func (m *monujo) securityDetails(ticker string) {
	exists, err := m.repository.SecurityExists(ticker)
	log.PanicIfError(err)
	if exists {
		return
	}

	if !m.yesOrNo(fmt.Sprintf("Would you like to add %s security detials to the database?", strings.TrimSpace(ticker))) {
		return
	}

	s := Security{
		Ticker: ticker,
	}
	s.ShortName = m.input.String("Short name")
	s.FullName = m.input.String("Full name")
	s.Market = m.input.String("Market")
	s.Leverage = m.input.Float("Leverage", 1)
	s.QuotesSource = m.input.String("Quotes source")
	tb := m.input.String("Ticker Bankier", "")
	s.TickerBankier = sql.NullString{String: tb, Valid: true}

	t, err := m.repository.StoreSecurity(s)
	log.PanicIfError(err)

	fmt.Printf("Security details of %s has been stored\n", strings.TrimSpace(t))
}

func (m *monujo) pickCurrency() string {
	fmt.Println("Choose currency")
	fmt.Println("")

	currencies, err := m.repository.Currencies()
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

	m.screen.DrawTable(header, data)
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
		return m.pickCurrency()
	}
}

func (m *monujo) shares() float64 {
	s := m.input.Float("Shares")

	isShort := m.isShort()
	if (isShort && s > 0) || (!isShort && s < 0) {
		return 0 - s
	}

	return s
}

func (m *monujo) isShort() bool {
	var input string
	fmt.Println("(B)UY or (S)ELL?")
	fmt.Scanln(&input)
	input = strings.ToUpper(input)

	if "S" == input {
		return true
	} else if "B" == input {
		return false
	}

	return m.isShort()
}

func (m *monujo) Dump(dumptype string, file string) {
	dbConf := m.config.Db()
	sysConf := m.config.Sys()
	if len(file) == 0 {
		fmt.Println("Output file is not set")
		return
	}

	var cmd *exec.Cmd
	if dumptype == "schema" {
		cmd = exec.Command(
			sysConf.Pgdump,
			"--host",
			dbConf.Host,
			"--port",
			dbConf.Port,
			"--username",
			dbConf.User,
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
			dbConf.Dbname,
		)
	} else if dumptype == "data" {
		cmd = exec.Command(
			sysConf.Pgdump,
			"--host",
			dbConf.Host,
			"--port",
			dbConf.Port,
			"--username",
			dbConf.User,
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
			dbConf.Dbname,
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
