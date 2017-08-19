package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/mgierok/monujo/config"
	"github.com/mgierok/monujo/console"
	"github.com/mgierok/monujo/db"
)

func main() {
	var env string
	var dump string
	var file string
	flag.StringVar(&env, "env", "", "force environment")
	flag.StringVar(&dump, "dump", "", "dump 'data' or 'schema'")
	flag.StringVar(&file, "file", "", "where to store the dump")
	flag.Parse()

	config.MustInitialize(env)
	db.MustInitialize()

	defer db.Connection().Close()

	if len(dump) > 0 {
		db.Dump(dump, file)
	} else {
		mainMenu()
	}
}

func mainMenu() {
	fmt.Println("Choose action")
	data := [][]interface{}{
		[]interface{}{"1", "Summary"},
		[]interface{}{"2", "Put transaction"},
		[]interface{}{"3", "List transactions"},
		[]interface{}{"4", "Update Quotes"},
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
	} else if action == "4" {
		runAction(Update)
	} else if action == "Q" {
		return
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
