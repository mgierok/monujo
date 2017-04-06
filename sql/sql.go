package sql

import (
	"database/sql"
	"reflect"
	"strconv"
)

func Sprint(v interface{}) string {
	s := ""
	t := reflect.TypeOf(v).String()
	if "sql.NullFloat64" == t {
		f := v.(sql.NullFloat64)
		if f.Valid {
			s = strconv.FormatFloat(f.Float64, 'f', -1, 64)
		}
	} else if "float64" == t {
		f := v.(float64)
		s = strconv.FormatFloat(f, 'f', -1, 64)
	} else if "sql.NullInt64" == t {
		i := v.(sql.NullInt64)
		if i.Valid {
			s = strconv.FormatInt(i.Int64, 10)
		}
	} else if "int64" == t {
		i := v.(int64)
		s = strconv.FormatInt(i, 10)
	}

	return s
}
