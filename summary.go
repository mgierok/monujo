package main

import (
	"fmt"

	"github.com/mgierok/monujo/repository"
	"github.com/mgierok/monujo/console"
	"github.com/mgierok/monujo/log"
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

	console.DrawTable(header, data)
}
