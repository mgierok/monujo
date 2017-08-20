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
		[]interface{}{"S", "Summary"},
		[]interface{}{"PT", "Put transaction"},
		[]interface{}{"LT", "List transactions"},
		[]interface{}{"PO", "Put operation"},
		[]interface{}{"LO", "List operations"},
		[]interface{}{"U", "Update Quotes"},
		[]interface{}{"Q", "Quit"},
	}

	console.DrawTable([]string{}, data)

	var action string
	fmt.Scanln(&action)
	action = strings.ToUpper(action)
	console.Clear()

	if action == "S" {
		runAction(Summary)
	} else if action == "PT" {
		runAction(PutTransaction)
	} else if action == "LT" {
		runAction(Transactions)
	} else if action == "PO" {
		runAction(PutOperation)
	} else if action == "LO" {
		runAction(Operations)
	} else if action == "U" {
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
