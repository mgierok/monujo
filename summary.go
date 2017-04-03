package main

import (
	"fmt"

	_ "github.com/lib/pq"
	"github.com/mgierok/monujo/repository"
)

func Summary() {
	ownedStocks, err := repository.OwnedStocks()
	LogError(err)

	var data [][]string

	for _, os := range ownedStocks {

		data = append(data, []string{
			os.PortfolioName.String,
			os.GetStockName(),
			os.Shares.String,
			os.LastPrice.String,
			os.AveragePrice.String,
			os.LastPriceBaseCurrency.String,
			os.AveragePriceBaseCurrency.String,
			os.Gain.String,
			os.GainBaseCurrency.String,
			os.PercentageGain.String,
			os.PercentageGainBaseCurrency.String,
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

	DrawTable(header, data)

	data = data[0:0]
	fmt.Println("")
	fmt.Println("")

	portfoliosSummary, err := repository.GetAllPortfoliosSummary()
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
	}

	DrawTable(header, data)
}
