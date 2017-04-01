package main

import (
	"fmt"
)

func main() {
	fmt.Println("Choose action")
	data := [][]string{
		[]string{"1", "Summary"},
		[]string{"2", "Put transaction"},
	}

	DrawTable([]string{}, data)

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
