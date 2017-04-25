package main

import (
	"fmt"
	"strings"

	"github.com/mgierok/monujo/repository"
	"github.com/mgierok/monujo/console"
)

func main() {
	fmt.Println("Choose action")
	data := [][]interface{}{
		[]interface{}{"1", "Summary"},
		[]interface{}{"2", "Put transaction"},
		[]interface{}{"3", "List transactions"},
		[]interface{}{"Q", "Quit"},
	}

	console.DrawTable([]string{}, data)

	var action string
	fmt.Scanln(&action)
	action = strings.ToUpper(action)
	console.Clear()

	db := GetDbConnection()
	repository.SetDb(db)
	defer db.Close()

	if action == "1" {
		Summary()
	} else if action == "2" {
		PutTransaction()
	} else if action == "3" {
		Transactions()
	} else if action == "Q" {
		return;
	}
}
