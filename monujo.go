package main

import (
	"fmt"

	"github.com/mgierok/monujo/repository"
	"github.com/mgierok/monujo/console"
)

func main() {
	fmt.Println("Choose action")
	data := [][]interface{}{
		[]interface{}{"1", "Summary"},
		[]interface{}{"2", "Put transaction"},
	}

	console.DrawTable([]string{}, data)

	var action int
	fmt.Scanln(&action)
	console.Clear()

	db := GetDbConnection()
	repository.SetDb(db)
	defer db.Close()

	if action == 1 {
		Summary()
	} else if action == 2 {
		PutTransaction()
	}
}
