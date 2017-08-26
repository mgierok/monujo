package entity

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/json-iterator/go"
)

type Source struct {
	Name string
}

type Sources []Source

func (s Source) Update(securities Securities, quotes chan Quote, wg *sync.WaitGroup) {
	defer wg.Done()
	if s.Name == "stooq" {
		stooq(securities, quotes)
	} else if s.Name == "ingturbo" {
		ingturbo(securities, quotes)
	} else if s.Name == "google" {
		google(securities, quotes)
	}
}

func stooq(securities Securities, quotes chan Quote) {
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
	// close(quotes)
}

func ingturbo(securities Securities, quotes chan Quote) {
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
	// close(quotes)
}

func google(securities Securities, quotes chan Quote) {
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
	// close(quotes)
}
