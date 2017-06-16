package entity

import "fmt"

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
