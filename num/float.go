package num

import "github.com/msw-x/moon/ufmt"

type Float64d1 float64

type Float64d2 float64

type Float64d3 float64

type Float64d4 float64

type Float64d5 float64

type Float64d6 float64

func (o *Float64d1) Set(v float64) {
	*o = Float64d1(v)
}

func (o *Float64d2) Set(v float64) {
	*o = Float64d2(v)
}

func (o *Float64d3) Set(v float64) {
	*o = Float64d3(v)
}

func (o *Float64d4) Set(v float64) {
	*o = Float64d4(v)
}

func (o *Float64d5) Set(v float64) {
	*o = Float64d5(v)
}

func (o *Float64d6) Set(v float64) {
	*o = Float64d6(v)
}

func (o Float64d1) Value() float64 {
	return float64(o)
}

func (o Float64d2) Value() float64 {
	return float64(o)
}

func (o Float64d3) Value() float64 {
	return float64(o)
}

func (o Float64d4) Value() float64 {
	return float64(o)
}

func (o Float64d5) Value() float64 {
	return float64(o)
}

func (o Float64d6) Value() float64 {
	return float64(o)
}

func (o Float64d1) String() string {
	return ufmt.DelicateFloat64(o.Value(), 1)
}

func (o Float64d2) String() string {
	return ufmt.DelicateFloat64(o.Value(), 2)
}

func (o Float64d3) String() string {
	return ufmt.DelicateFloat64(o.Value(), 3)
}

func (o Float64d4) String() string {
	return ufmt.DelicateFloat64(o.Value(), 4)
}

func (o Float64d5) String() string {
	return ufmt.DelicateFloat64(o.Value(), 5)
}

func (o Float64d6) String() string {
	return ufmt.DelicateFloat64(o.Value(), 6)
}

func (o Float64d1) MarshalJSON() ([]byte, error) {
	return []byte(o.String()), nil
}

func (o Float64d2) MarshalJSON() ([]byte, error) {
	return []byte(o.String()), nil
}

func (o Float64d3) MarshalJSON() ([]byte, error) {
	return []byte(o.String()), nil
}

func (o Float64d4) MarshalJSON() ([]byte, error) {
	return []byte(o.String()), nil
}

func (o Float64d5) MarshalJSON() ([]byte, error) {
	return []byte(o.String()), nil
}

func (o Float64d6) MarshalJSON() ([]byte, error) {
	return []byte(o.String()), nil
}
