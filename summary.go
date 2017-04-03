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
			os.DisplayName(),
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

	portfoliosExt, err := repository.PortfoliosExt()
	LogError(err)

	for _, pe := range portfoliosExt {
		data = append(data, []string{
			pe.PortfolioId.String,
			pe.Name.String,
			pe.CacheValue.String,
			pe.GainOfSoldShares.String,
			pe.Commision.String,
			pe.Tax.String,
			pe.GainOfOwnedShares.String,
			pe.EstimatedGain.String,
			pe.EstimatedGainCostsInc.String,
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
