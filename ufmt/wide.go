package ufmt

import (
	"strconv"

	"golang.org/x/exp/constraints"
)

func WideInt[V constraints.Integer](v V) string {
	s := strconv.FormatInt(int64(v), 10)
	parts := []string{}
	for len(s) > 3 {
		parts = append(parts, s[len(s)-3:])
		s = s[:len(s)-3]
	}
	if len(s) > 0 {
		parts = append(parts, s)
	}
	partsCount := len(parts)
	var ret string
	for i := partsCount - 1; i >= 0; i-- {
		if len(ret) != 0 {
			ret += " "
		}
		ret += parts[i]
	}
	return ret
}

func WideFloat(v float64, precision int) string {
	i := int64(v)
	d := v - float64(i)
	r := ""
	if d > 0 {
		r = Float64(d, precision)
		if len(r) > 0 {
			r = r[1:]
		}
	}
	return WideInt(i) + r
}
