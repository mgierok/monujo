package main

import (
	"os"

	"github.com/olekukonko/tablewriter"
)

func DrawTable(header []string, data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)
	table.AppendBulk(data)
	table.SetRowLine(true)
	table.Render()
}
