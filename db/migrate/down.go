package migrate

import (
	"fmt"
	"slices"
	"strings"
)

func GenerateDown(upSql string) (string, error) {
	return GenerateDownWithContext(upSql, NewContext())
}

func GenerateDownLine(up string) (string, error) {
	return GenerateDownLineWithContext(up, nil)
}

func GenerateDownWithContext(upSql string, c *Context) (down string, err error) {
	var l []string
	for _, s := range strings.Split(upSql, "\n") {
		s, err = GenerateDownLineWithContext(s, c)
		if err != nil {
			break
		}
		if s != "" {
			l = append(l, s)
		}
	}
	if err == nil {
		slices.Reverse(l)
		down = strings.Join(l, "\n")
	}
	return
}

func GenerateDownLineWithContext(up string, c *Context) (down string, err error) {
	g := NewDownGenerator(c)
	l := strings.Fields(up)
	for _, f := range l {
		if g.Process(f) {
			down = g.Command
			break
		}
	}
	err = g.Error
	if err == nil && c != nil {
		err = c.Process(l)
	}
	return
}

type DownGenerator struct {
	Process func(string) bool
	Command string
	Error   error

	c *Context
	l []string
}

func NewDownGenerator(c *Context) *DownGenerator {
	o := new(DownGenerator)
	o.Process = o.init
	o.c = c
	return o
}

func (o *DownGenerator) insert(token string, f func(string) bool) {
	if o.Error != nil {
		return
	}
	if f == nil {
		if !(strings.HasPrefix(token, "(") && strings.HasSuffix(token, ")")) {
			token = trimNameTail(token)
		}
	}
	if token != "" {
		o.l = append(o.l, token)
	}
	if f == nil {
		o.Command = strings.Join(o.l, " ") + ";"
	}
	o.Process = f
}

func (o *DownGenerator) replace(token string, f func(string) bool) {
	o.l[len(o.l)-1] = token
	o.insert("", f)
}

func (o *DownGenerator) last() string {
	return o.l[len(o.l)-1]
}

func (o *DownGenerator) lastIs(s string) bool {
	return len(o.l) > 0 && o.last() == s
}

func (o *DownGenerator) init(token string) (r bool) {
	v := strings.ToUpper(token)
	switch v {
	case "CREATE":
		o.insert("DROP", o.create)
	case "ALTER":
		o.insert(v, o.alter)
	default:
		r = true
	}
	return
}

func (o *DownGenerator) create(token string) (r bool) {
	v := strings.ToUpper(token)
	switch v {
	case "UNIQUE", "OR", "REPLACE":
	case "SCHEMA", "SEQUENCE", "TYPE", "PROCEDURE", "FUNCTION", "TABLE":
		o.insert(v, nil)
		o.insert("IF EXISTS", o.name)
	case "INDEX":
		o.insert(v, nil)
		o.insert("IF EXISTS", o.nameOfIndex)
	default:
		r = true
	}
	return
}

func (o *DownGenerator) alter(token string) (r bool) {
	v := strings.ToUpper(token)
	switch v {
	case "ADD":
		o.insert("DROP", o.add)
	case "DROP":
		o.insert(v, o.alter)
	case "COLUMN":
		if o.lastIs("DROP") {
			o.replace("ADD", o.name)
		} else {
			o.insert(v, o.alter)
		}
	case "CONSTRAINT":
		if o.lastIs("DROP") {
			o.replace("ADD", o.constraint)
		}
	default:
		if o.lastIs("COLUMN") {
			o.insert(token, o.column)
		} else {
			o.insert(token, o.alter)
		}
	}
	return
}

func (o *DownGenerator) add(token string) bool {
	if token == "PRIMARY" {
		o.insert("CONSTRAINT", o.primary)
		return false
	}
	o.insert("COLUMN", nil)
	return o.name(token)
}

func (o *DownGenerator) primary(token string) bool {
	o.insert("", o.key)
	return false
}

func (o *DownGenerator) key(token string) bool {
	if strings.HasSuffix(token, ",") {
		o.insert("", o.key)
		return false
	}
	table := o.l[2]
	l := strings.Split(table, ".")
	tableLocal := l[len(l)-1]
	o.insert(tableLocal+"_pkey", nil)
	return true
}

func (o *DownGenerator) column(token string) (r bool) {
	v := strings.ToUpper(token)
	switch v {
	case "TYPE":
		o.insert(v, o.column)
	default:
		if o.lastIs("TYPE") {
			columnType, _, err := o.c.schema.ColumnType(o.l[2], o.l[5])
			if err == nil {
				o.insert(columnType, nil)
			} else {
				o.Error = err
			}
			r = true
		}
	}
	return
}

func (o *DownGenerator) constraint(token string) (r bool) {
	table := o.l[2]
	l := strings.Split(table, ".")
	tableLocal := l[len(l)-1]
	token = trimNameTail(token)
	token = strings.TrimPrefix(token, tableLocal+"_")
	column := strings.TrimSuffix(token, "_key")
	constraint, err := o.c.schema.RemoveColumnKeyConstraint(table, column)
	if err == nil {
		o.insert(constraint, nil)
	} else {
		o.Error = err
	}
	o.insert(fmt.Sprintf("(%s)", column), nil)
	return true
}

func (o *DownGenerator) name(token string) (r bool) {
	v := strings.ToUpper(token)
	switch v {
	case "IF", "NOT", "EXISTS":
	default:
		lastIsAdd := o.lastIs("ADD")
		o.insert(token, nil)
		r = true
		if lastIsAdd {
			columnType, columnConstraints, err := o.c.schema.ColumnType(o.l[2], o.l[4])
			if err == nil {
				o.insert(columnType, nil)
				if columnConstraints != "" {
					o.insert(columnConstraints, nil)
				}
			} else {
				o.Error = err
			}
		}
	}
	return
}

func (o *DownGenerator) nameOfIndex(token string) (r bool) {
	v := strings.ToUpper(token)
	switch v {
	case "CONCURRENTLY":
	default:
		o.insert(token, nil)
		r = true
	}
	return
}

func trimNameTail(s string) string {
	i := strings.IndexAny(s, "{(;")
	if i != -1 {
		s = s[:i]
	}
	return s
}
