package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/mgierok/monujo/console"
	"github.com/mgierok/monujo/log"
	"github.com/mgierok/monujo/repository"
	"github.com/mgierok/monujo/repository/entity"
)

var availableSources = map[string]func([]string) entity.Quotes{
	"stooq":    stooq,
	"google":   google,
	"ingturbo": ingturbo,
}

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

	for source, f := range sources() {
		if len(importMap[source]) > 0 {
			var quotes entity.Quotes
			if source == source {
				quotes = f(importMap[source])
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
}

func sources() map[string]func([]string) entity.Quotes {
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
	for s, _ := range availableSources {
		dict[strconv.Itoa(i)] = s
		data = append(data, []interface{}{strconv.Itoa(i), s})
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
			return availableSources
		} else if input == "Q" {
			return map[string]func([]string) entity.Quotes{}
		} else {
			return map[string]func([]string) entity.Quotes{
				dict[input]: availableSources[dict[input]],
			}
		}
	} else {
		return sources()
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

func ingturbo(tickers []string) entity.Quotes {
	type response struct {
		BidQuotes [][]float64 `json:"BidQuotes"`
	}
	var quotes entity.Quotes

	var client http.Client

	for _, t := range tickers {
		subt := strings.Trim(strings.ToLower(t), " ")
		resp, err := client.Get(
			fmt.Sprintf(
				"https://www.ingturbo.pl/services/product/PLINGNV%s/chart?period=intraday",
				subt[len(t)-5:],
			),
		)
		if err != nil {
			fmt.Printf("Update failed for %s\n", t)
		} else {
			body, _ := ioutil.ReadAll(resp.Body)
			var r response
			_ = json.Unmarshal(body, &r)
			v := r.BidQuotes[len(r.BidQuotes)-1][1]
			quote := entity.Quote{
				Ticker:  t,
				Date:    time.Now(),
				Open:    v,
				High:    v,
				Low:     v,
				Close:   v,
				Volume:  0,
				OpenInt: 0,
			}

			quotes = append(quotes, quote)
		}
	}

	return quotes
}

func google(tickers []string) entity.Quotes {
	type gQuote struct {
		Ticker   string `json:"t"`
		Exchange string `json:"e"`
		Quote    string `json:"l_fix"`
		Date     string `json:"lt_dts"`
	}
	var quotes entity.Quotes

	securities, err := repository.Securities(tickers)
	log.PanicIfError(err)

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
		_ = json.Unmarshal(body, &gQuotes)

		for _, gQuote := range gQuotes {
			v, _ := strconv.ParseFloat(gQuote.Quote, 64)
			quote := entity.Quote{
				Ticker:  gmap[gQuote.Exchange+":"+gQuote.Ticker],
				Open:    v,
				High:    v,
				Low:     v,
				Close:   v,
				Volume:  0,
				OpenInt: 0,
			}
			quote.Date, _ = time.Parse("2006-01-02T15:04:05Z", gQuote.Date)

			quotes = append(quotes, quote)
		}
	}

	return quotes
}
