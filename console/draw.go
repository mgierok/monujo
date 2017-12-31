package console

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/olekukonko/tablewriter"
)

func (c *Console) PrintTable(header []string, data [][]interface{}) {
	newData := make([][]string, len(data))
	for i, r := range data {
		newData[i] = make([]string, len(data[i]))
		for j, e := range r {
			newData[i][j] = c.toString(e)
		}
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)
	table.AppendBulk(newData)
	table.SetRowLine(true)
	table.Render()
}

func (c *Console) PrintText(format string, a ...interface{}) {
	fmt.Printf(format, a...)
}

func (c *Console) NewLine(n ...int) {
	var lines int
	if len(n) == 0 {
		lines = 1
	} else {
		lines = n[0]
	}

	for i := 0; i < lines; i++ {
		c.PrintText("\n")
	}
}

func (c *Console) Clear() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func (c *Console) toString(v interface{}) string {
	s := ""
	switch v.(type) {
	case sql.NullFloat64:
		if vv := v.(sql.NullFloat64); vv.Valid {
			s = strconv.FormatFloat(vv.Float64, 'f', -1, 64)
		}
	case float64:
		s = strconv.FormatFloat(v.(float64), 'f', -1, 64)
	case sql.NullInt64:
		if vv := v.(sql.NullInt64); vv.Valid {
			s = strconv.FormatInt(vv.Int64, 10)
		}
	case int64:
		s = strconv.FormatInt(v.(int64), 10)
	case time.Time:
		s = v.(time.Time).Format("2006-01-02")
	case string:
		s = v.(string)
	}

	return s
}
