package main

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
)

func main() {
	fmt.Println("Choose action")
	table := tablewriter.NewWriter(os.Stdout)
	table.AppendBulk(
		[][]string{
			[]string{"1", "Summary"},
			[]string{"2", "Put transaction"},
		},
	)
	table.Render()

	var action int
	fmt.Scanln(&action)

	db := GetDbConnection()
	defer db.Close()

	if action == 1 {
		Summary(db)
	} else if action == 2 {
		PutTransaction(db)
	}
}
