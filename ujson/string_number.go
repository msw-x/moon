package ujson

import "github.com/msw-x/moon/parse"

type StringNumber string

func (o StringNumber) Empty() bool {
	return o == ""
}

func (o StringNumber) Exists() bool {
	return !o.Empty()
}

func (o StringNumber) Int64() (int64, error) {
	return parse.Int64(string(o))
}

func (o StringNumber) Float64() (float64, error) {
	return parse.Float64(string(o))
}
