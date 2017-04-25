package main

import (
	"fmt"
	"strings"

	"github.com/mgierok/monujo/repository"
	"github.com/mgierok/monujo/console"
)

func main() {
	db := GetDbConnection()
	repository.SetDb(db)
	defer db.Close()

	mainMenu()
}

func mainMenu() {
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

	if action == "1" {
		runAction(Summary)
	} else if action == "2" {
		runAction(PutTransaction)
	} else if action == "3" {
		runAction(Transactions)
	} else if action == "Q" {
		return;
	} else {
		mainMenu()
	}
}

func runAction(f func()) {
	console.Clear()
	f()

	// type something to continue TODO how to detect enter key?
	var input string
	fmt.Scanln(&input)

	mainMenu()
}
