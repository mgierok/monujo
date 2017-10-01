package action

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/mgierok/monujo/console"
	"github.com/mgierok/monujo/log"
	"github.com/mgierok/monujo/repository"
	"github.com/mgierok/monujo/repository/entity"
)

func Summary() {
	ownedStocks, err := repository.OwnedStocks()
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

	console.DrawTable(header, data)

	data = data[0:0]
	fmt.Println("")
	fmt.Println("")

	portfoliosExt, err := repository.PortfoliosExt()
	log.PanicIfError(err)

	for _, pe := range portfoliosExt {
		data = append(data, []interface{}{
			pe.PortfolioId,
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
		"Portfolio Id",
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

	console.DrawTable(header, data)
}

func ListTransactions() {
	portfolio := portfolio()

	transactions, err := repository.PortfolioTransactions(portfolio)
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

	console.DrawTable(header, data)
	fmt.Println("")

	if !YesOrNo("Do you want to delete single transaction?") {
		return
	}

	transaction := pickTransaction(transactions)
	err = repository.DeleteTransaction(transaction)
	log.PanicIfError(err)
	fmt.Println("Transaction has been removed")
}

func ListOperations() {
	portfolio := portfolio()

	operations, err := repository.PortfolioOperations(portfolio)
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

	console.DrawTable(header, data)
	fmt.Println("")

	if !YesOrNo("Do you want to delete single financial operation?") {
		return
	}

	operation := pickOperation(operations)
	err = repository.DeleteOperation(operation)
	log.PanicIfError(err)
	fmt.Println("Operation has been removed")
}

func Update() {
	ownedStocks, err := repository.OwnedStocks()
	log.PanicIfError(err)
	currencies, err := repository.Currencies()
	log.PanicIfError(err)

	var importMap = make(map[string]entity.Securities)
	tickers := ownedStocks.DistinctTickers()
	tickers = append(tickers, currencies.CurrencyPairs("PLN")...)

	securities, err := repository.Securities(tickers)
	log.PanicIfError(err)

	for _, t := range tickers {
		for _, s := range securities {
			if s.Ticker == t {
				importMap[s.QuotesSource] = append(importMap[s.QuotesSource], s)
			}
		}
	}

	var wg sync.WaitGroup
	quotes := make(chan entity.Quote)
	for _, source := range pickSource() {
		securities := importMap[source.Name]
		if len(securities) > 0 {
			wg.Add(1)
			go source.Update(securities, quotes, &wg)
		}
	}

	go func() {
		wg.Wait()
		close(quotes)
	}()

	for q := range quotes {
		_, err = repository.StoreLatestQuote(q)
		if err == nil {
			fmt.Printf("Ticker: %s Quote: %f\n", q.Ticker, q.Close)
		} else {
			fmt.Printf("Update failed for %s\n", q.Ticker)
		}
	}
}

func PutOperation() {
	var o entity.Operation
	o.PortfolioId = portfolio().PortfolioId
	console.Clear()
	o.Date = console.InputDate("Date", time.Now())
	console.Clear()
	o.Type = financialOperationType().Type
	console.Clear()
	o.Value = console.InputFloat("Value")
	console.Clear()
	o.Description = console.InputString("Description", "")
	console.Clear()
	o.Commision = console.InputFloat("Commision", 0)
	console.Clear()
	o.Tax = console.InputFloat("Tax", 0)
	console.Clear()

	summary := [][]interface{}{
		[]interface{}{"Portfolio ID", o.PortfolioId},
		[]interface{}{"Date", o.Date},
		[]interface{}{"Operation type", o.Type},
		[]interface{}{"Value", o.Value},
		[]interface{}{"Description", o.Description},
		[]interface{}{"Commision", o.Commision},
		[]interface{}{"Tax", o.Tax},
	}

	console.Clear()
	console.DrawTable([]string{}, summary)
	fmt.Println("")

	if YesOrNo("Do you want to store this operation?") {
		operationId, err := repository.StoreOperation(o)
		log.PanicIfError(err)

		fmt.Printf("Operation has been recorded with an ID: %d\n", operationId)
	} else {
		fmt.Println("Operation has not been recorded")
	}

}

func PutTransaction() {
	var t entity.Transaction
	t.PortfolioId = portfolio().PortfolioId
	console.Clear()
	t.Date = console.InputDate("Date", time.Now())
	console.Clear()
	t.Ticker = console.InputString("Ticker")
	console.Clear()
	t.Price = console.InputFloat("Price")
	console.Clear()
	t.Currency = pickCurrency()
	console.Clear()
	t.Shares = shares()
	console.Clear()
	t.Commision = console.InputFloat("Commision", 0)
	console.Clear()
	t.ExchangeRate = console.InputFloat("Exchange rate", 1)
	console.Clear()
	t.Tax = console.InputFloat("Tax", 0)
	console.Clear()

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

	console.Clear()
	console.DrawTable([]string{}, summary)
	fmt.Println("")

	if YesOrNo("Do you want to store this transaction?") {
		transactionId, err := repository.StoreTransaction(t)
		log.PanicIfError(err)

		fmt.Printf("Transaction has been recorded with an ID: %d\n", transactionId)

		securityDetails(t.Ticker)
	} else {
		fmt.Println("Transaction has not been recorded")
	}
}

func pickTransaction(transactions entity.Transactions) entity.Transaction {
	var input string
	fmt.Print("Transaction ID: ")
	fmt.Scanln(&input)

	transactionId, err := strconv.ParseInt(input, 10, 64)

	if nil != err {
		fmt.Printf("\n%s is not a valid transaction ID\n\n", input)
		return pickTransaction(transactions)
	} else {
		for _, t := range transactions {
			if t.TransactionId == transactionId {
				return t
			}
		}

		fmt.Printf("\n%s is not a valid transaction ID\n\n", input)
		return pickTransaction(transactions)
	}
}

func pickOperation(operations entity.Operations) entity.Operation {
	var input string
	fmt.Print("Operation ID: ")
	fmt.Scanln(&input)

	operationId, err := strconv.ParseInt(input, 10, 64)

	if nil != err {
		fmt.Printf("\n%s is not a valid operation ID\n\n", input)
		return pickOperation(operations)
	} else {
		for _, o := range operations {
			if o.OperationId == operationId {
				return o
			}
		}

		fmt.Printf("\n%s is not a valid operation ID\n\n", input)
		return pickOperation(operations)
	}
}

func YesOrNo(question string) bool {
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

	return YesOrNo(question)
}

func portfolio() entity.Portfolio {
	fmt.Println("Choose portfolio")
	fmt.Println("")

	portfolios, err := repository.Portfolios()
	log.PanicIfError(err)

	header := []string{
		"Portfolio Id",
		"Portfolio Name",
	}

	var data [][]interface{}
	for _, p := range portfolios {
		data = append(data, []interface{}{p.PortfolioId, p.Name})
	}

	console.DrawTable(header, data)
	fmt.Println("")

	var input string
	fmt.Print("Portfolio ID: ")
	fmt.Scanln(&input)

	portfolioId, err := strconv.ParseInt(input, 10, 64)

	if nil != err {
		fmt.Printf("\n%s is not a valid portfolio ID\n\n", input)
		return portfolio()
	} else {
		for _, p := range portfolios {
			if p.PortfolioId == portfolioId {
				return p
			}
		}

		fmt.Printf("\n%s is not a valid portfolio ID\n\n", input)
		return portfolio()
	}
}

func pickSource() entity.Sources {
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
	for _, s := range repository.Sources() {
		dict[strconv.Itoa(i)] = s.Name
		data = append(data, []interface{}{strconv.Itoa(i), s.Name})
		i++
	}
	data = append(data, []interface{}{"Q", "Quit"})

	console.DrawTable([]string{}, data)
	fmt.Println("")

	var input string
	fmt.Scanln(&input)
	console.Clear()

	input = strings.ToUpper(input)

	_, exists := dict[input]
	if exists {
		if input == "A" {
			return repository.Sources()
		} else if input == "Q" {
			return entity.Sources{}
		} else {
			for _, s := range repository.Sources() {
				if s.Name == dict[input] {
					return entity.Sources{s}
				}
			}
		}
	}
	return pickSource()
}

func financialOperationType() entity.FinancialOperationType {
	fmt.Println("Choose operation type")
	fmt.Println("")

	ots, err := repository.FinancialOperationTypes()
	log.PanicIfError(err)

	header := []string{
		"Operation type",
	}

	var dict = make(map[string]entity.FinancialOperationType)
	var data [][]interface{}
	for _, ot := range ots {
		dict[ot.Type] = ot
		data = append(data, []interface{}{ot.Type})
	}

	console.DrawTable(header, data)
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
		return financialOperationType()
	}
}

func securityDetails(ticker string) {
	exists, err := repository.SecurityExists(ticker)
	log.PanicIfError(err)
	if exists {
		return
	}

	if !YesOrNo(fmt.Sprintf("Would you like to add %s security detials to the database?", strings.TrimSpace(ticker))) {
		return
	}

	s := entity.Security{
		Ticker: ticker,
	}
	s.ShortName = console.InputString("Short name")
	s.FullName = console.InputString("Full name")
	s.Market = console.InputString("Market")
	s.Leverage = console.InputFloat("Leverage", 1)
	s.QuotesSource = console.InputString("Quotes source")

	t, err := repository.StoreSecurity(s)
	log.PanicIfError(err)

	fmt.Printf("Security details of %s has been stored\n", strings.TrimSpace(t))
}

func pickCurrency() string {
	fmt.Println("Choose currency")
	fmt.Println("")

	currencies, err := repository.Currencies()
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

	console.DrawTable(header, data)
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
		return pickCurrency()
	}
}

func shares() float64 {
	s := console.InputFloat("Shares")

	isShort := isShort()
	if (isShort && s > 0) || (!isShort && s < 0) {
		return 0 - s
	}

	return s
}

func isShort() bool {
	var input string
	fmt.Println("(B)UY or (S)ELL?")
	fmt.Scanln(&input)
	input = strings.ToUpper(input)

	if "S" == input {
		return true
	} else if "B" == input {
		return false
	}

	return isShort()
}
