package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mgierok/monujo/console"
	"github.com/mgierok/monujo/log"
	"github.com/mgierok/monujo/repository"
	"github.com/mgierok/monujo/repository/entity"
)

func PutOperation() {

	var o entity.Operation
	getO(oPortfolioId, &o)
	getO(oDate, &o)
	getO(oOperationType, &o)
	getO(oValue, &o)
	getO(oDescription, &o)
	getO(oCommision, &o)

	summary := [][]interface{}{
		[]interface{}{"Portfolio ID", o.PortfolioId},
		[]interface{}{"Date", o.Date},
		[]interface{}{"Operation type", o.Type},
		[]interface{}{"Value", o.Value},
		[]interface{}{"Description", o.Description},
		[]interface{}{"Commision", o.Commision},
	}

	console.Clear()
	console.DrawTable([]string{}, summary)
	fmt.Println("")

	fmt.Println("Type 'Y' to insert or 'N' to abort")
	if confirm() {
		operationId, err := repository.StoreOperation(o)
		log.PanicIfError(err)

		fmt.Printf("Operation has been recorded with an ID: %d\n", operationId)
	} else {
		fmt.Println("Operation has not been recorded")
	}

}

func getO(f func(*entity.Operation), o *entity.Operation) {
	console.Clear()
	f(o)
}

func oPortfolioId(e *entity.Operation) {
	fmt.Println("Choose portfolio")
	fmt.Println("")

	portfolios, err := repository.Portfolios()
	log.PanicIfError(err)

	header := []string{
		"Portfolio Id",
		"Portfolio Name",
	}

	var dict = make(map[int64]string)
	var data [][]interface{}
	for _, p := range portfolios {
		data = append(data, []interface{}{p.PortfolioId, p.Name})
		dict[p.PortfolioId] = p.Name
	}

	console.DrawTable(header, data)
	fmt.Println("")

	var input string
	fmt.Print("Portfolio ID: ")
	fmt.Scanln(&input)

	p, err := strconv.ParseInt(input, 10, 64)

	if nil != err {
		fmt.Printf("\n%s is not a valid portfolio ID\n\n", input)
		getO(oPortfolioId, e)
		return
	} else {
		_, exists := dict[p]
		if exists {
			e.PortfolioId = p
		} else {
			fmt.Printf("\n%d is not a valid portfolio ID\n\n", p)
			getO(oPortfolioId, e)
			return
		}
	}
}

func oDate(e *entity.Operation) {
	const layout = "2006-01-02"
	now := time.Now()
	var input string

	fmt.Printf("Date (default: %q): ", now.Format(layout))
	fmt.Scanln(&input)
	input = strings.Trim(input, " ")

	if input == "" {
		e.Date = now
	} else {
		t, err := time.Parse(layout, input)
		if err != nil {
			fmt.Println(err)
			fmt.Printf("\n%q is not a valid date format\n\n", input)
			getO(oDate, e)
			return
		} else {
			e.Date = t
		}
	}
}

func oOperationType(e *entity.Operation) {
	fmt.Println("Choose operation type")
	fmt.Println("")

	ots, err := repository.FinancialOperationTypes()
	log.PanicIfError(err)

	header := []string{
		"Operation type",
	}

	var dict = make(map[string]string)
	var data [][]interface{}
	for _, ot := range ots {
		dict[ot.Type] = ot.Type
		data = append(data, []interface{}{ot.Type})
	}

	console.DrawTable(header, data)
	fmt.Println("")

	var ot string
	fmt.Print("Type: ")
	fmt.Scanln(&ot)

	ot = strings.ToLower(ot)

	_, exists := dict[ot]
	if exists {
		e.Type = ot
	} else {
		fmt.Printf("\n%s is not a valid operation type\n\n", ot)
		getO(oOperationType, e)
		return
	}
}

func oValue(e *entity.Operation) {
	fmt.Print("Value: ")
	var input string
	fmt.Scanln(&input)

	v, err := strconv.ParseFloat(input, 64)

	if err != nil {
		fmt.Printf("\n%s is not a valid value\n\n", input)
		getO(oValue, e)
		return
	}

	e.Value = v
}

func oDescription(e *entity.Operation) {
	fmt.Print("Description: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()

	input = strings.TrimSpace(input)

	e.Description = input
}

func oCommision(e *entity.Operation) {
	fmt.Print("Commision (default: 0): ")
	var input string
	fmt.Scanln(&input)

	c, err := strconv.ParseFloat(input, 64)

	if input == "" {
		c = 0
	} else {
		if err != nil {
			fmt.Printf("\n%s is not a valid commision value\n\n", input)
			getO(oCommision, e)
			return
		}
	}

	e.Commision = c
}
