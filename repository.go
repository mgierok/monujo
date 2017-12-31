package monujo

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	jsoniter "github.com/json-iterator/go"
)

type Repository struct {
	db     *sqlx.DB
	config Config
}

func NewRepository(db *sqlx.DB, c *Config) (*Repository, error) {
	return &Repository{
		db:     db,
		config: *c,
	}, nil
}

type Currency struct {
	Symbol string `db:"currency"`
}

type Currencies []Currency

func (currencies *Currencies) CurrencyPairs(base string) []string {
	pairs := make([]string, 0, len(*currencies))
	for _, c := range *currencies {
		if c.Symbol == base {
			continue
		}
		pairs = append(pairs, fmt.Sprintf("%-12s", c.Symbol+base))
	}

	return pairs
}

func (r *Repository) Currencies() (Currencies, error) {
	currencies := Currencies{}
	err := r.db.Select(&currencies, "SELECT currency FROM currencies ORDER BY currency ASC")
	return currencies, err
}

type Operation struct {
	OperationId int64     `db:"operation_id"`
	PortfolioId int64     `db:"portfolio_id"`
	Date        time.Time `db:"date"`
	Type        string    `db:"type"`
	Value       float64   `db:"value"`
	Description string    `db:"description"`
	Commision   float64   `db:"commision"`
	Tax         float64   `db:"tax"`
}

type Operations []Operation

func (r *Repository) StoreOperation(operation Operation) (int64, error) {
	stmt, err := r.db.PrepareNamed(`
		INSERT INTO operations (portfolio_id, date, type, value, description, commision, tax)
		VALUES (:portfolio_id, :date, :type, :value, :description, :commision, :tax)
		RETURNING operation_id
	`)

	var operationId int64
	if nil == err {
		err = stmt.Get(&operationId, operation)
	}

	return operationId, err
}

func (r *Repository) PortfolioOperations(portfolio Portfolio) (Operations, error) {
	operations := Operations{}
	err := r.db.Select(&operations,
		`SELECT
		operation_id,
		portfolio_id,
		date,
		type,
		value,
		COALESCE(description, '') AS description,
		commision,
		tax
	FROM operations
	WHERE portfolio_id = $1
	ORDER BY
		date ASC,
		operation_id ASC
	`,
		portfolio.PortfolioId)
	return operations, err
}

func (r *Repository) DeleteOperation(operation Operation) error {
	_, err := r.db.Exec("DELETE FROM operations  WHERE portfolio_id = $1 AND operation_id = $2", operation.PortfolioId, operation.OperationId)
	return err
}

type OperationType struct {
	Type string `db:"type"`
}

type OperationTypes []OperationType

type FinancialOperationType struct {
	OperationType
}

type FinancialOperationTypes []FinancialOperationType

func (r *Repository) FinancialOperationTypes() (FinancialOperationTypes, error) {
	types := FinancialOperationTypes{}
	err := r.db.Select(&types, "SELECT type FROM financial_operation_types ORDER BY type ASC")
	return types, err
}

type Portfolio struct {
	PortfolioId int64  `db:"portfolio_id"`
	Name        string `db:"name"`
	Currency    string `db:"currency"`
}

type Portfolios []Portfolio

type PortfolioExt struct {
	Portfolio
	CacheValue            sql.NullFloat64 `db:"cache_value"`
	GainOfSoldShares      sql.NullFloat64 `db:"gain_of_sold_shares"`
	Commision             sql.NullFloat64 `db:"commision"`
	Tax                   sql.NullFloat64 `db:"tax"`
	GainOfOwnedShares     sql.NullFloat64 `db:"gain_of_owned_shares"`
	EstimatedGain         sql.NullFloat64 `db:"estimated_gain"`
	EstimatedGainCostsInc sql.NullFloat64 `db:"estimated_gain_costs_inc"`
	EstimatedValue        sql.NullFloat64 `db:"estimated_value"`
	AnnualBalance         float64         `db:"annual_balance"`
	MonthBalance          float64         `db:"month_balance"`
}

type PortfoliosExt []PortfolioExt

func (r *Repository) PortfoliosExt() (PortfoliosExt, error) {
	portfolios := PortfoliosExt{}
	err := r.db.Select(&portfolios, "SELECT portfolio_id, name, currency, cache_value, gain_of_sold_shares, commision, tax, gain_of_owned_shares, estimated_gain, estimated_gain_costs_inc, estimated_value, annual_balance, month_balance FROM portfolios_ext ORDER BY portfolio_id")
	return portfolios, err
}

func (r *Repository) Portfolios() (Portfolios, error) {
	portfolios := Portfolios{}
	err := r.db.Select(&portfolios, "SELECT portfolio_id, name, currency FROM portfolios ORDER BY portfolio_id ASC")
	return portfolios, err
}

type Quote struct {
	Ticker  string    `db:"ticker"`
	Date    time.Time `db:"date"`
	Open    float64   `db:"open"`
	High    float64   `db:"high"`
	Low     float64   `db:"low"`
	Close   float64   `db:"close"`
	Volume  float64   `db:"volume"`
	OpenInt float64   `db:"openint"`
}

type Quotes []Quote

func (r *Repository) StoreLatestQuote(quote Quote) (string, error) {
	stmt, err := r.db.PrepareNamed(`
		INSERT INTO latest_quotes (ticker, date, open, high, low, close, volume, openint)
		VALUES (:ticker, :date, :open, :high, :low, :close, :volume, :openint)
		RETURNING ticker
	`)

	var t string
	if nil == err {
		err = stmt.Get(&t, quote)
	}

	return t, err
}

type Security struct {
	Ticker        string         `db:"ticker"`
	ShortName     string         `db:"short_name"`
	FullName      string         `db:"full_name"`
	Market        string         `db:"market"`
	Leverage      float64        `db:"leverage"`
	QuotesSource  string         `db:"quotes_source"`
	TickerBankier sql.NullString `db:"ticker_bankier"`
}

type Securities []Security

func (r *Repository) SecurityExists(ticker string) (bool, error) {
	var exists bool
	err := r.db.Get(
		&exists,
		`SELECT
			COUNT(1)
		FROM securities
		WHERE ticker = $1
		`,
		ticker,
	)
	return exists, err
}

func (r *Repository) Securities(tickers []string) (Securities, error) {
	s := Securities{}
	var query string
	var err error

	if len(tickers) > 0 {
		var args []interface{}
		query, args, err = sqlx.In(
			`SELECT
				ticker,
				short_name,
				full_name,
				market,
				leverage,
				quotes_source,
				ticker_bankier
			FROM securities
			WHERE ticker IN (?)
			`,
			tickers,
		)
		query = r.db.Rebind(query)
		err = r.db.Select(&s, query, args...)
	} else {
		query =
			`SELECT
				ticker,
				short_name,
				full_name,
				market,
				leverage,
				quotes_source,
				ticker_bankier
			FROM securities`
		err = r.db.Select(&s, query)
	}

	return s, err
}

func (r *Repository) StoreSecurity(s Security) (string, error) {
	stmt, err := r.db.PrepareNamed(`
		INSERT INTO securities (ticker, short_name, full_name, market, leverage, quotes_source, ticker_bankier)
		VALUES (:ticker, :short_name, :full_name, :market, :leverage, :quotes_source, :ticker_bankier)
		RETURNING ticker
	`)

	var t string
	if nil == err {
		err = stmt.Get(&t, s)
	}

	return t, err
}

type Source struct {
	Name string `db:"name"`
}

type Stooq Source
type Ingturbo Source
type Google Source
type Alphavantage Source
type Bankier Source

type Sources []Source

func (s Source) Update(securities Securities, quotes chan Quote, wg *sync.WaitGroup, config AppConf) {
	defer wg.Done()
	if s.Name == "stooq" {
		ss := Stooq(s)
		ss.stooq(securities, quotes)
	} else if s.Name == "ingturbo" {
		ss := Ingturbo(s)
		ss.ingturbo(securities, quotes)
	} else if s.Name == "google" {
		ss := Stooq(s)
		ss.stooq(securities, quotes)
	} else if s.Name == "alphavantage" {
		ss := Alphavantage(s)
		ss.alphavantage(securities, quotes, config.Alphavantagekey)
	} else if s.Name == "bankier" {
		ss := Bankier(s)
		ss.update(securities, quotes)
	}
}

func (s Stooq) stooq(securities Securities, quotes chan Quote) {
	const layout = "20060102"
	now := time.Now()
	var client http.Client
	for _, s := range securities {
		resp, err := client.Get(
			fmt.Sprintf(
				"https://stooq.pl/q/d/l/?s=%s&d1=%s&d2=%s&i=d",
				strings.Trim(strings.ToLower(s.Ticker), " "),
				now.AddDate(0, 0, -7).Format(layout),
				now.Format(layout),
			),
		)
		if err != nil {
			fmt.Printf("Update failed for %s\n", s.Ticker)
		} else {
			body, _ := ioutil.ReadAll(resp.Body)
			csvBody := string(body)

			r := csv.NewReader(strings.NewReader(csvBody))

			records, err := r.ReadAll()
			if err != nil {
				fmt.Printf("Ticker: %s Error: %s\n", s.Ticker, err)
			} else if len(records[0]) == 1 {
				fmt.Printf("Ticker: %s Error: %s\n", s.Ticker, records[0][0])
			} else {
				last := len(records) - 1
				quote := Quote{
					Ticker:  s.Ticker,
					Volume:  0,
					OpenInt: 0,
				}
				quote.Date, _ = time.Parse("2006-01-02", records[last][0])
				quote.Open, _ = strconv.ParseFloat(records[last][1], 64)
				quote.High, _ = strconv.ParseFloat(records[last][2], 64)
				quote.Low, _ = strconv.ParseFloat(records[last][3], 64)
				quote.Close, _ = strconv.ParseFloat(records[last][4], 64)

				quotes <- quote
			}
		}
	}
}

func (s Ingturbo) ingturbo(securities Securities, quotes chan Quote) {
	type response struct {
		BidQuotes [][]float64 `json:"BidQuotes"`
	}

	var client http.Client

	for _, s := range securities {
		subt := strings.Trim(strings.ToLower(s.Ticker), " ")
		resp, err := client.Get(
			fmt.Sprintf(
				"https://www.ingturbo.pl/services/product/PLINGNV%s/chart?period=intraday",
				subt[len(s.Ticker)-5:],
			),
		)
		if err != nil {
			fmt.Printf("Update failed for %s\n", s.Ticker)
		} else {
			body, _ := ioutil.ReadAll(resp.Body)
			var r response
			_ = jsoniter.Unmarshal(body, &r)
			if len(r.BidQuotes) == 0 {
				fmt.Printf("Update failed for %s\n", s.Ticker)
			} else {
				v := r.BidQuotes[len(r.BidQuotes)-1][1]
				quote := Quote{
					Ticker:  s.Ticker,
					Date:    time.Now(),
					Open:    v,
					High:    v,
					Low:     v,
					Close:   v,
					Volume:  0,
					OpenInt: 0,
				}

				quotes <- quote
			}
		}
	}
}

func (s Alphavantage) alphavantage(securities Securities, quotes chan Quote, key string) {
	var client http.Client
	for _, s := range securities {
		resp, err := client.Get(
			fmt.Sprintf(
				"https://www.alphavantage.co/query?function=TIME_SERIES_DAILY&apikey=%s&datatype=csv&symbol=%s",
				key,
				strings.TrimSuffix(strings.TrimSpace(s.Ticker), ".US"),
			),
		)
		if err != nil {
			fmt.Printf("Update failed for %s\n", s.Ticker)
		} else {
			body, _ := ioutil.ReadAll(resp.Body)
			csvBody := string(body)

			r := csv.NewReader(strings.NewReader(csvBody))

			records, err := r.ReadAll()
			if err != nil {
				fmt.Printf("Ticker: %s Error: %s\n", s.Ticker, err)
			} else if len(records[0]) == 1 {
				fmt.Printf("Ticker: %s Error: %s\n", s.Ticker, records[0][0])
			} else {
				quote := Quote{
					Ticker:  s.Ticker,
					Volume:  0,
					OpenInt: 0,
				}
				quote.Date, _ = time.Parse("2006-01-02", records[1][0])
				quote.Open, _ = strconv.ParseFloat(records[1][1], 64)
				quote.High, _ = strconv.ParseFloat(records[1][2], 64)
				quote.Low, _ = strconv.ParseFloat(records[1][3], 64)
				quote.Close, _ = strconv.ParseFloat(records[1][4], 64)

				quotes <- quote
			}
		}
	}
}

func (s Google) google(securities Securities, quotes chan Quote) {
	type gQuote struct {
		Ticker   string `json:"t"`
		Exchange string `json:"e"`
		Quote    string `json:"l_fix"`
		QuoteC   string `json:"l"`
		Date     string `json:"lt_dts"`
	}

	var gtickers []string
	var gmap = make(map[string]string)
	for _, s := range securities {
		gticker := s.Market + ":" + strings.TrimSuffix(strings.TrimSpace(s.Ticker), ".US")
		gtickers = append(gtickers, gticker)
		gmap[gticker] = s.Ticker
	}

	var client http.Client
	resp, err := client.Get(
		fmt.Sprintf(
			"https://finance.google.com/finance/info?client=ig&q=%s",
			strings.Join(gtickers, ","),
		),
	)
	if err != nil {
		fmt.Println("Update from Google failed")
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		body = body[4:] // remove comment sign at the beginning of response

		var gQuotes []gQuote
		_ = jsoniter.Unmarshal(body, &gQuotes)

		for _, gQuote := range gQuotes {
			v, _ := strconv.ParseFloat(gQuote.Quote, 64)
			if v == 0 {
				v, _ = strconv.ParseFloat(gQuote.QuoteC, 64)
			}
			quote := Quote{
				Ticker:  gmap[gQuote.Exchange+":"+gQuote.Ticker],
				Open:    v,
				High:    v,
				Low:     v,
				Close:   v,
				Volume:  0,
				OpenInt: 0,
			}
			quote.Date, _ = time.Parse("2006-01-02T15:04:05Z", gQuote.Date)

			quotes <- quote
		}
	}
}

func (s Bankier) update(securities Securities, quotes chan Quote) {
	type bQuote struct {
		Open   float64
		High   float64
		Low    float64
		Close  float64
		Volume float64
		Date   time.Time
	}
	var bQuotes = make(map[string]bQuote)
	var client http.Client
	var toFloat = func(s string) float64 {
		s = strings.Replace(s, "&nbsp;", "", -1)
		s = strings.Replace(s, ",", ".", -1)
		v, _ := strconv.ParseFloat(s, 64)
		return v
	}

	regex, _ := regexp.Compile(`(?sU)<td class="colWalor textNowrap">.+<a title=".+" href=".+">(.+)</a>.+<td class="colKurs change.+">(.+)</td>.+<td class="colObrot">(.+)</td>.+<td class="colOtwarcie">(.+)</td>.+<td class="calMaxi">(.+)</td>.+<td class="calMini">(.+)</td>.+<td class="colAktualizacja">(.+)</td>`)
	urls := [3]string{
		"https://www.bankier.pl/gielda/notowania/akcje",
		"https://www.bankier.pl/gielda/notowania/new-connect",
		"https://www.bankier.pl/gielda/notowania/futures",
	}
	for _, url := range urls {
		resp, err := client.Get(url)
		if err != nil {
			fmt.Println(err)
			fmt.Printf("Unable to read %s\n", url)
		} else {
			body, _ := ioutil.ReadAll(resp.Body)
			matches := regex.FindAllStringSubmatch(string(body), -1)

			for _, row := range matches {
				close := toFloat(row[2])
				volume := toFloat(row[3])
				open := toFloat(row[4])
				high := toFloat(row[5])
				low := toFloat(row[6])
				date, _ := time.Parse("2006.01.02 15:04", time.Now().Format("2006")+"."+row[7])

				bQuotes[strings.ToUpper(row[1])] = bQuote{
					Open:   open,
					High:   high,
					Low:    low,
					Close:  close,
					Volume: volume,
					Date:   date,
				}
			}
		}
	}

	regex, _ = regexp.Compile(`(?sU)<td class="colTicker"><a href="/fundusze/notowania/(.+)">.+<td class="colKurs">(.+)</td>.+<td class="colAktualizacja textNowrap">(.+)</td>`)
	url := "https://www.bankier.pl/fundusze/notowania/wszystkie"
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println(err)
		fmt.Printf("Unable to read %s\n", url)
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		matches := regex.FindAllStringSubmatch(string(body), -1)

		for _, row := range matches {
			v := toFloat(row[2])
			date, _ := time.Parse("2006-01-02", row[3])

			bQuotes[strings.ToUpper(row[1])] = bQuote{
				Open:   v,
				High:   v,
				Low:    v,
				Close:  v,
				Volume: 0,
				Date:   date,
			}
		}
	}

	for _, s := range securities {
		q, ok := bQuotes[strings.Trim(s.TickerBankier.String, " ")]
		if ok {
			quote := Quote{
				Ticker:  s.Ticker,
				Open:    q.Open,
				High:    q.High,
				Low:     q.Low,
				Close:   q.Close,
				Volume:  q.Volume,
				OpenInt: 0,
				Date:    q.Date,
			}

			quotes <- quote
		} else {
			fmt.Printf("Update failed for %s\n", s.Ticker)
		}
	}
}

func (r *Repository) Sources() (Sources, error) {
	s := Sources{}
	err := r.db.Select(&s,
		`SELECT
			DISTINCT quotes_source AS name
			FROM securities
			ORDER BY quotes_source
			`)
	return s, err
}

func (r *Repository) UpdateQuotes(sources Sources) (chan Quote, error) {
	quotes := make(chan Quote)

	ownedStocks, err := r.OwnedStocks()
	if err != nil {
		return quotes, err
	}
	currencies, err := r.Currencies()
	if err != nil {
		return quotes, err
	}

	var importMap = make(map[string]Securities)
	tickers := ownedStocks.DistinctTickers()
	tickers = append(tickers, currencies.CurrencyPairs("PLN")...)

	securities, err := r.Securities(tickers)
	if err != nil {
		return quotes, err
	}

	for _, t := range tickers {
		for _, s := range securities {
			if s.Ticker == t {
				importMap[s.QuotesSource] = append(importMap[s.QuotesSource], s)
			}
		}
	}

	var wg sync.WaitGroup
	for _, source := range sources {
		securities := importMap[source.Name]
		if len(securities) > 0 {
			wg.Add(1)
			go source.Update(securities, quotes, &wg, r.config.App)
		}
	}

	go func() {
		wg.Wait()
		close(quotes)
	}()

	return quotes, nil
}

type Stock struct {
	Ticker    string          `db:"ticker"`
	ShortName sql.NullString  `db:"short_name"`
	LastPrice sql.NullFloat64 `db:"last_price"`
	Currency  string          `db:"currency"`
}

type OwnedStock struct {
	Stock
	PortfolioId                int64           `db:"portfolio_id"`
	PortfolioName              string          `db:"portfolio_name"`
	Shares                     float64         `db:"shares"`
	ExchangeRate               sql.NullFloat64 `db:"exchange_rate"`
	LastPriceBaseCurrency      sql.NullFloat64 `db:"last_price_base_currency"`
	AveragePrice               float64         `db:"average_price"`
	AveragePriceBaseCurrency   float64         `db:"average_price_base_currency"`
	InvestmentBaseCurrency     sql.NullFloat64 `db:"investment_base_currency"`
	MarketValueBaseCurrency    sql.NullFloat64 `db:"market_value_base_currency"`
	Gain                       sql.NullFloat64 `db:"gain"`
	PercentageGain             sql.NullFloat64 `db:"percentage_gain"`
	GainBaseCurrency           sql.NullFloat64 `db:"gain_base_currency"`
	PercentageGainBaseCurrency sql.NullFloat64 `db:"percentage_gain_base_currency"`
}

type OwnedStocks []OwnedStock

func (stock *Stock) DisplayName() string {
	if stock.ShortName.Valid {
		return fmt.Sprintf("%s (%s)", stock.ShortName.String, strings.Trim(stock.Ticker, " "))
	} else {
		return strings.Trim(stock.Ticker, " ")
	}
}

func (stocks *OwnedStocks) DistinctTickers() []string {
	t := make(map[string]struct{})
	for _, stock := range *stocks {
		t[stock.Ticker] = struct{}{}
	}

	keys := make([]string, 0, len(t))
	for k := range t {
		keys = append(keys, k)
	}

	return keys
}

func (r *Repository) OwnedStocks() (OwnedStocks, error) {
	stocks := OwnedStocks{}
	err := r.db.Select(&stocks,
		`SELECT
			portfolio_id,
			portfolio_name,
			ticker,
			short_name,
			shares,
			last_price,
			currency,
			exchange_rate,
			last_price_base_currency,
			average_price,
			average_price_base_currency,
			average_price_base_currency,
			investment_base_currency,
			market_value_base_currency,
			gain,
			percentage_gain,
			gain_base_currency,
			percentage_gain_base_currency
			FROM owned_stocks
			ORDER BY portfolio_id
			`)
	return stocks, err
}

type Transaction struct {
	TransactionId int64     `db:"transaction_id"`
	PortfolioId   int64     `db:"portfolio_id"`
	Date          time.Time `db:"date"`
	Ticker        string    `db:"ticker"`
	Price         float64   `db:"price"`
	Currency      string    `db:"currency"`
	Shares        float64   `db:"shares"`
	Commision     float64   `db:"commision"`
	ExchangeRate  float64   `db:"exchange_rate"`
	Tax           float64   `db:"tax"`
}

type Transactions []Transaction

func (r *Repository) StoreTransaction(transaction Transaction) (int64, error) {
	stmt, err := r.db.PrepareNamed(`
		INSERT INTO transactions (portfolio_id, date, ticker, price, currency, shares, commision, exchange_rate, tax)
		VALUES (:portfolio_id, :date, :ticker, :price, :currency, :shares, :commision, :exchange_rate, :tax)
		RETURNING transaction_id
	`)

	var transactionId int64
	if nil == err {
		err = stmt.Get(&transactionId, transaction)
	}

	return transactionId, err
}

func (r *Repository) PortfolioTransactions(portfolio Portfolio) (Transactions, error) {
	transactions := Transactions{}
	err := r.db.Select(&transactions,
		`SELECT
		transaction_id,
		portfolio_id,
		date,
		ticker,
		price,
		currency,
		shares,
		commision,
		exchange_rate,
		tax
	FROM transactions
	WHERE portfolio_id = $1
	ORDER BY
		date ASC,
		transaction_id ASC
	`,
		portfolio.PortfolioId)
	return transactions, err
}

func (r *Repository) DeleteTransaction(transaction Transaction) error {
	_, err := r.db.Exec("DELETE FROM transactions WHERE portfolio_id = $1 AND transaction_id = $2", transaction.PortfolioId, transaction.TransactionId)
	return err
}
