package migrate

import (
	"strings"
)

type Context struct {
	schema      Schema
	activeTable string
}

func NewContext() *Context {
	o := new(Context)
	return o
}

func (o *Context) Process(l []string) error {
	n := len(l)
	if n == 0 {
		return nil
	}
	if strings.ContainsAny(l[n-1], "};") {
		o.activeTable = ""
	}
	l[n-1] = trimNameTail(l[n-1])
	if l[n-1] == "" {
		n--
	}
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
							return o.schema.AddColumn(tableName, l[4], l[5])
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
			o.schema.AddColumn(o.activeTable, l[0], l[1])
		}
	}
	return nil
}
