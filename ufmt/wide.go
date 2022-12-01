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
