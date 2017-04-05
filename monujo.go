package main

import (
	"fmt"

	"github.com/mgierok/monujo/repository"
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
	repository.SetDb(db)
	defer db.Close()

	if action == 1 {
		Summary()
	} else if action == 2 {
		PutTransaction()
	}
}
