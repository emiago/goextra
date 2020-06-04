package mysqlhelp

import (
	"fmt"
	"reflect"
	"strings"
)

/* Query filters, used for building queries*/

type QueryFilter struct {
	column    string
	value     interface{}
	operation string
}

func (qf *QueryFilter) ToSQL() string {
	vstr := ""

	// dlog.WithField("field", prettyFormat(unknown)).Info("TOSQL")

	switch fval := qf.value.(type) {
	case []interface{}:
		if len(fval) == 0 {
			break
		}

		for _, val := range fval {
			switch v := val.(type) {
			case string:
				vstr += fmt.Sprintf(",'%s'", v)
			case int:
				vstr += fmt.Sprintf(",%d", v)
			default:
				vstr = ""
				break
			}
		}
		if len(vstr) > 0 {
			//Skip ','
			vstr = fmt.Sprintf("(%s)", vstr[1:])
		}
	case string:
		vstr = fmt.Sprintf("'%s'", fval)
	case int, int64:
		vstr = fmt.Sprintf("%d", fval)
	case float32, float64:
		vstr = fmt.Sprintf("%f", fval)
	default:
		//Try to convert
		switch reflect.ValueOf(qf.value).Kind() {
		case reflect.Int, reflect.Int64:
			vstr = fmt.Sprintf("%d", fval)
		case reflect.String:
			vstr = fmt.Sprintf("'%s'", fval)
		default:
			vstr = fmt.Sprintf("'%v'", fval)
		}
	}

	return vstr

}

func QueryFilterBuild(filters []*QueryFilter) string {
	q := ""
	and := ""
	for _, v := range filters {
		q += fmt.Sprintf("%s %s %s %s", and, v.column, v.operation, v.ToSQL())
		and = " AND"
	}
	return strings.TrimSpace(q)
}
