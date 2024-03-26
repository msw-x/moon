package migrate

import (
	"slices"
	"strings"
)

func GenerateDown(upSql string) string {
	var l []string
	for _, s := range strings.Split(upSql, "\n") {
		s = GenerateDownLine(s)
		if s != "" {
			l = append(l, s)
		}
	}
	slices.Reverse(l)
	return strings.Join(l, "\n")
}

func GenerateDownLine(up string) (down string) {
	g := NewDownGenerator()
	for _, f := range strings.Fields(up) {
		if g.Process(f) {
			down = g.Command
			break
		}
	}
	return
}

type DownGenerator struct {
	Process func(string) bool
	Command string

	l []string
}

func NewDownGenerator() *DownGenerator {
	o := new(DownGenerator)
	o.Process = o.init
	return o
}

func (o *DownGenerator) add(token string, f func(string) bool) {
	if f == nil {
		token = trimNameTail(token)
	}
	o.l = append(o.l, token)
	if f == nil {
		o.Command = strings.Join(o.l, " ") + ";"
	}
	o.Process = f
}

func (o *DownGenerator) lastIs(s string) bool {
	return len(o.l) > 0 && o.l[len(o.l)-1] == s
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
		o.add("DROP", o.name)
	case "COLUMN":
		o.add(v, o.alter)
	default:
		if o.lastIs("COLUMN") {
			o.add(token, nil)
			r = true
		} else {
			o.add(token, o.alter)
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
