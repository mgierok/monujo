package action

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

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
			os.LastPriceBaseCurrency,
			os.AveragePriceBaseCurrency,
			os.Gain,
			os.GainBaseCurrency,
			os.PercentageGain,
			os.PercentageGainBaseCurrency,
		})
	}

	header := []string{
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
			pe.Commision,
			pe.Tax,
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
		"Commision",
		"Tax",
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
