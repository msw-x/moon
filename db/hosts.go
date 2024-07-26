package db

import "strings"

type Hosts struct {
	l []string
	i int
}

func NewHosts(s string) *Hosts {
	o := new(Hosts)
	o.l = strings.Split(s, ",")
	for i, h := range o.l {
		o.l[i] = host(h)
	}
	return o
}

func (o Hosts) Count() int {
	return len(o.l)
}

func (o Hosts) String() string {
	return strings.Join(o.l, " ")
}

func (o Hosts) IsMulti() bool {
	return o.Count() > 1
}

func (o *Hosts) Next() (v string) {
	v = o.l[o.i]
	o.i++
	if o.i == o.Count() {
		o.i = 0
	}
	return
}

func (o *Hosts) Loop(f func(string) bool) {
	for i := 0; i != o.Count(); i++ {
		if f(o.Next()) {
			break
		}
	}
}
