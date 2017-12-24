package console

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func (c *Console) Float(name string, args ...float64) float64 {
	if len(args) == 1 {
		fmt.Printf("%s (default: %v): ", name, args[0])
	} else {
		fmt.Printf("%s: ", name)
	}

	var input string
	fmt.Scanln(&input)

	if input == "" {
		return args[0]
	}

	f, err := strconv.ParseFloat(input, 64)

	if err != nil {
		fmt.Printf("\n%s is not a valid %s\n\n", input, strings.ToLower(name))
		return c.Float(name, args...)
	}

	return f
}

func (c *Console) String(name string, args ...string) string {
	if len(args) == 1 {
		fmt.Printf("%s (default: %v): ", name, args[0])
	} else {
		fmt.Printf("%s: ", name)
	}
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()
	input = strings.TrimSpace(input)

	if input == "" {
		return args[0]
	}

	return input
}

func (c *Console) Date(name string, args ...time.Time) time.Time {
	const layout = "2006-01-02"

	if len(args) == 1 {
		fmt.Printf("%s (default: %q): ", name, args[0].Format(layout))
	} else {
		fmt.Printf("%s: ", name)
	}

	var input string
	fmt.Scanln(&input)
	input = strings.Trim(input, " ")

	if input == "" {
		return args[0]
	} else {
		t, err := time.Parse(layout, input)
		if err != nil {
			fmt.Println(err)
			fmt.Printf("\n%q is not a valid date format\n\n", input)
			return c.Date(name, args...)

		} else {
			return t
		}
	}
}
