package entity

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/json-iterator/go"
	"github.com/mgierok/monujo/config"
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
		stooq(securities, quotes)
	} else if s.Name == "alphavantage" {
		alphavantage(securities, quotes)
	} else if s.Name == "bankier" {
		bankier(securities, quotes)
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
}

func alphavantage(securities Securities, quotes chan Quote) {
	var client http.Client
	for _, s := range securities {
		resp, err := client.Get(
			fmt.Sprintf(
				"https://www.alphavantage.co/query?function=TIME_SERIES_DAILY&apikey=%s&datatype=csv&symbol=%s",
				config.App().Alphavantagekey,
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
}

func bankier(securities Securities, quotes chan Quote) {
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
	urls := [2]string{
		"https://www.bankier.pl/gielda/notowania/akcje",
		"https://www.bankier.pl/gielda/notowania/new-connect",
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
