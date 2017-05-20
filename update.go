package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/mgierok/monujo/log"
	"github.com/mgierok/monujo/repository"
	"github.com/mgierok/monujo/repository/entity"
)

func Update() {
	ownedStocks, err := repository.OwnedStocks()
	log.PanicIfError(err)

	var importMap = make(map[string][]string)
	tickers := ownedStocks.DistinctTickers()
	securities, err := repository.Securities(tickers)
	log.PanicIfError(err)

	for _, t := range tickers {
		for _, s := range securities {
			if s.Ticker == t {
				importMap[s.QuotesSource] = append(importMap[s.QuotesSource], t)
			}
		}
	}

	for source, tickers := range importMap {
		var quotes entity.Quotes
		if source == "stooq" {
			quotes = stooq(tickers)
		}

		for _, q := range quotes {
			_, err = repository.StoreLatestQuote(q)
			if err == nil {
				fmt.Printf("Ticker: %s Quote: %f\n", q.Ticker, q.Close)
			} else {
				fmt.Printf("Update failed for %s\n", q.Ticker)
			}
		}
	}
}

func stooq(tickers []string) entity.Quotes {
	var quotes entity.Quotes

	const layout = "20060102"
	now := time.Now()
	var client http.Client

	for _, t := range tickers {
		resp, err := client.Get(
			fmt.Sprintf(
				"https://stooq.pl/q/d/l/?s=%s&d1=%s&d2=%s&i=d",
				strings.Trim(strings.ToLower(t), " "),
				now.AddDate(0, 0, -7).Format(layout),
				now.Format(layout),
			),
		)
		if err != nil {
			fmt.Printf("Update failed for %s\n", t)
		} else {
			body, _ := ioutil.ReadAll(resp.Body)
			csvBody := string(body)

			r := csv.NewReader(strings.NewReader(csvBody))

			records, err := r.ReadAll()
			if err != nil {
				fmt.Println(err)
			}

			last := len(records) - 1
			quote := entity.Quote{
				Ticker:  t,
				Volume:  0,
				OpenInt: 0,
			}
			quote.Date, _ = time.Parse("2006-01-02", records[last][0])
			quote.Open, _ = strconv.ParseFloat(records[last][1], 64)
			quote.High, _ = strconv.ParseFloat(records[last][2], 64)
			quote.Low, _ = strconv.ParseFloat(records[last][3], 64)
			quote.Close, _ = strconv.ParseFloat(records[last][4], 64)

			quotes = append(quotes, quote)
		}
	}

	return quotes
}
