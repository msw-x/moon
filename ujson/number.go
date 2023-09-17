package ujson

import (
	"github.com/msw-x/moon/parse"
	"github.com/msw-x/moon/umath"
)

type Int64 float64

func (o *Int64) UnmarshalJSON(b []byte) error {
	return unmarshalNumber(b, parse.Int64, o)
}

func (o Int64) MarshalJSON() ([]byte, error) {
	return quote(o), nil
}

func (o Int64) Value() int64 {
	return int64(o)
}

func (o Int64) Ptr() *int64 {
	v := o.Value()
	return &v
}

func (o Int64) IsZero() bool {
	return o.Value() == 0
}

func (o Int64) NilIfZero() *int64 {
	if o.IsZero() {
		return nil
	}
	return o.Ptr()
}

func (o *Int64) Set(v int64) {
	*o = Int64(v)
}

type Float64 float64

func (o *Float64) UnmarshalJSON(b []byte) error {
	return unmarshalNumber(b, parse.Float64, o)
}

func (o Float64) MarshalJSON() ([]byte, error) {
	return quote(o), nil
}

func (o Float64) Value() float64 {
	return float64(o)
}

func (o Float64) Ptr() *float64 {
	v := o.Value()
	return &v
}

func (o Float64) IsZero() bool {
	return o.Value() == 0
}

func (o Float64) NilIfZero() *float64 {
	if o.IsZero() {
		return nil
	}
	return o.Ptr()
}

func (o *Float64) Set(v float64) {
	*o = Float64(v)
}

func unmarshalNumber[N Int64 | Float64, T umath.Number](b []byte, parse func(string) (T, error), n *N) error {
	s := unquote(b)
	if s == "" || s == "null" {
		*n = 0
		return nil
	}
	v, err := parse(s)
	*n = N(v)
	return err
}
