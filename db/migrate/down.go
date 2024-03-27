package migrate

import (
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

func (o *DownGenerator) add(token string, f func(string) bool) {
	if o.Error != nil {
		return
	}
	if f == nil {
		token = trimNameTail(token)
	}
	o.l = append(o.l, token)
	if f == nil {
		o.Command = strings.Join(o.l, " ") + ";"
	}
	o.Process = f
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
		o.add("DROP", o.create)
	case "ALTER":
		o.add(v, o.alter)
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
		o.add(v, nil)
		o.add("IF EXISTS", o.name)
	case "INDEX":
		o.add(v, nil)
		o.add("IF EXISTS", o.nameOfIndex)
	default:
		r = true
	}
	return
}

func (o *DownGenerator) alter(token string) (r bool) {
	v := strings.ToUpper(token)
	switch v {
	case "ADD":
		o.add("DROP", nil)
		o.add("COLUMN", o.name)
	case "COLUMN":
		o.add(v, o.alter)
	default:
		if o.lastIs("COLUMN") {
			o.add(token, o.column)
		} else {
			o.add(token, o.alter)
		}
	}
	return
}

func (o *DownGenerator) column(token string) (r bool) {
	v := strings.ToUpper(token)
	switch v {
	case "TYPE":
		o.add(v, o.column)
	default:
		if o.lastIs("TYPE") {
			columnType, err := o.c.schema.ColumnType(o.l[2], o.l[5])
			if err == nil {
				o.add(columnType, nil)
			} else {
				o.Error = err
			}
			r = true
		}
	}
	return
}

func (o *DownGenerator) name(token string) (r bool) {
	v := strings.ToUpper(token)
	switch v {
	case "IF", "NOT", "EXISTS":
	default:
		o.add(token, nil)
		r = true
	}
	return
}

func (o *DownGenerator) nameOfIndex(token string) (r bool) {
	v := strings.ToUpper(token)
	switch v {
	case "CONCURRENTLY":
	default:
		o.add(token, nil)
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
