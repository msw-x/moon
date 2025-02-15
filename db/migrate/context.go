package migrate

import (
	"strings"

	"github.com/msw-x/moon/ufmt"
)

type Context struct {
	schema      Schema
	activeTable string
}

func NewContext() *Context {
	o := new(Context)
	return o
}

func (o *Context) String() string {
	return o.schema.String()
}

func (o *Context) Pretty() string {
	return o.schema.Pretty()
}

func (o *Context) Process(l []string) error {
	n := len(l)
	if n == 0 {
		return nil
	}
	last := l[n-1]
	if strings.ContainsAny(last, "};") {
		o.activeTable = ""
	}
	if last == "(" || last == "," || last == ";" {
		n--
		l = l[:n]
	}
	l[n-1] = strings.TrimSuffix(l[n-1], ",")
	l[n-1] = strings.TrimSuffix(l[n-1], ";")
	if o.activeTable == "" {
		if n > 2 {
			if l[1] == "TABLE" {
				if l[0] == "CREATE" {
					tableName := l[n-1]
					if tableName == "" && n > 3 {
						tableName = l[n-2]
					}
					o.activeTable = tableName
					return o.schema.AddTable(tableName)
				}
				if n > 5 {
					if l[0] == "ALTER" {
						tableName := l[2]
						if l[3] == "ADD" {
							if l[4] == "PRIMARY" {
								///
								return nil
							} else {
								columnName := l[4]
								columnType := l[5]
								constraints := ""
								if n > 6 {
									constraints = ufmt.JoinSlice(l[6:])
								}
								return o.addColumn(tableName, columnName, columnType, constraints)
							}
						}
						if l[3] == "ALTER" && l[4] == "COLUMN" && l[6] == "TYPE" {
							return o.schema.AlterColumnType(tableName, l[5], l[7])
						}
					}
				}
			}
		}
	} else {
		if n > 1 {
			columnName := l[0]
			columnType := l[1]
			constraints := ""
			if n > 2 {
				constraints = ufmt.JoinSlice(l[2:])
			}
			return o.addColumn(o.activeTable, columnName, columnType, constraints)
		}
	}
	return nil
}

func (o *Context) addColumn(tableName string, columnName string, columnType string, columnConstraints string) error {
	return o.schema.AddColumn(tableName, columnName, columnType, columnConstraints)
}
