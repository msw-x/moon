package ufmt

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/msw-x/moon/umath"
)

func Float64(v float64, precision int) (s string) {
	if precision > 0 {
		f := "%." + strconv.Itoa(precision) + "f"
		s = fmt.Sprintf(f, v)
		s = strings.TrimRight(s, "0")
		s = strings.TrimRight(s, ".")
	} else {
		s = strconv.Itoa(int(v))
	}
	return
}

func DelicateFloat64(v float64, precision int) (s string) {
	if v == 0 {
		return "0"
	}
	if math.Abs(v) < 1 {
		precision += -umath.Order(v) - 1
	}
	return Float64(v, precision)
}
