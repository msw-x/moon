package ujson

import "github.com/msw-x/moon/parse"

type StringNumber string

func (o StringNumber) Empty() bool {
	return o == "" || o == "null"
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

type StringInt64 string

func (o StringInt64) Empty() bool {
	return o == "" || o == "null"
}

func (o StringInt64) Exists() bool {
	return !o.Empty()
}

func (o StringInt64) Value() (int64, error) {
	return parse.Int64(string(o))
}

type StringFloat64 string

func (o StringFloat64) Empty() bool {
	return o == "" || o == "null"
}

func (o StringFloat64) Exists() bool {
	return !o.Empty()
}

func (o StringFloat64) Value() (float64, error) {
	return parse.Float64(string(o))
}
