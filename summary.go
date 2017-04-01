package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/mgierok/monujo/repository"
	"github.com/olekukonko/tablewriter"
)

func Summary(db *sqlx.DB) {
	ownedStocksSummary, err := repository.GetAllOwnedStocksSummary(db)
	LogError(err)

	var data [][]string

	for _, oss := range ownedStocksSummary {
		var stock string
		if oss.ShortName.String == "" {
			stock = strings.Trim(oss.Ticker.String, " ")
		} else {
			stock = oss.ShortName.String
		}

		data = append(data, []string{
			oss.PortfolioName.String,
			stock,
			oss.Shares.String,
			oss.LastPrice.String,
			oss.AveragePrice.String,
			oss.LastPriceBaseCurrency.String,
			oss.AveragePriceBaseCurrency.String,
			oss.Gain.String,
			oss.GainBaseCurrency.String,
			oss.PercentageGain.String,
			oss.PercentageGainBaseCurrency.String,
		})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
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
	})

	table.AppendBulk(data)
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)
	table.Render()

	data = data[0:0]
	fmt.Println("")
	fmt.Println("")

	portfoliosSummary, err := repository.GetAllPortfoliosSummary(db)
	LogError(err)

	for _, ps := range portfoliosSummary {
		data = append(data, []string{
			ps.PortfolioId.String,
			ps.Name.String,
			ps.CacheValue.String,
			ps.GainOfSoldShares.String,
			ps.Commision.String,
			ps.Tax.String,
			ps.GainOfOwnedShares.String,
			ps.EstimatedGain.String,
			ps.EstimatedGainCostsInc.String,
		})
	}

	table = tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"Portfolio Id",
		"Portfolio Name",
		"Cache Value",
		"Gain of Sold Shares",
		"Commision",
		"Tax",
		"Gain Of Ownded Shares",
		"Estimated Gain",
		"Estimated Gain Costs Inc.",
	})

	table.AppendBulk(data)
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)
	table.Render()
}
