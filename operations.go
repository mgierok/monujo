package main

import (
	"fmt"
	"strconv"

	"github.com/mgierok/monujo/console"
	"github.com/mgierok/monujo/log"
	"github.com/mgierok/monujo/repository"
	"github.com/mgierok/monujo/repository/entity"
)

func Operations() {
	portfolio := portfolio()

	operations, err := repository.PortfolioOperations(portfolio)
	log.PanicIfError(err)

	var data [][]interface{}

	for _, o := range operations {
		data = append(data, []interface{}{
			o.OperationId,
			o.PortfolioId,
			o.Date,
			o.Type,
			o.Value,
			o.Description,
			o.Commision,
		})
	}

	header := []string{
		"Operation ID",
		"Portfolio ID",
		"Date",
		"Type",
		"Value",
		"Description",
		"Commision",
	}

	console.DrawTable(header, data)
	fmt.Println("")

	if !careToDelete() {
		return
	}

	operation := operation(operations)
	err = repository.DeleteOperation(operation)
	log.PanicIfError(err)
	fmt.Println("Operation has been removed")
}

func operation(operations entity.Operations) entity.Operation {
	var input string
	fmt.Print("Operation ID: ")
	fmt.Scanln(&input)

	operationId, err := strconv.ParseInt(input, 10, 64)

	if nil != err {
		fmt.Printf("\n%s is not a valid operation ID\n\n", input)
		return operation(operations)
	} else {
		for _, o := range operations {
			if o.OperationId == operationId {
				return o
			}
		}

		fmt.Printf("\n%s is not a valid operation ID\n\n", input)
		return operation(operations)
	}
}
