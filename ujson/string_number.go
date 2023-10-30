package ujson

import (
	"fmt"

	"github.com/msw-x/moon/parse"
	"github.com/msw-x/moon/ufmt"
)

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

func (o StringInt64) ValueOr(v int64) int64 {
	x, err := o.Value()
	if err == nil {
		return x
	}
	return v
}

func (o *StringInt64) Set(v int64) {
	*o = StringInt64(fmt.Sprint(v))
}

func (o StringInt64) ValueOrDefault() int64 {
	var v int64
	return o.ValueOr(v)
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

func (o StringFloat64) ValueOr(v float64) float64 {
	x, err := o.Value()
	if err == nil {
		return x
	}
	return v
}

func (o StringFloat64) ValueOrDefault() float64 {
	var v float64
	return o.ValueOr(v)
}

func (o *StringFloat64) Set(v float64) {
	*o = StringFloat64(fmt.Sprint(v))
}

func (o *StringFloat64) SetWithPrecision(v float64, precison int) {
	*o = StringFloat64(ufmt.Float64(v, precison))
}
