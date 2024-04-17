package ufmt

import (
	"testing"
)

func TestWideFloat64(t *testing.T) {
	test := func(v float64, p int, e string) {
		r := WideFloat(v, p)
		if r != e {
			t.Errorf("WideFloat(%v, %v) = %s; expected %s", v, p, r, e)
		}
	}
	test(0.0009, 2, "0")
	test(93845.0009, 2, "93 845")
	test(1243.009, 2, "1 243.01")
	test(3950.325, 2, "3 950.32")
	test(31.325, 2, "31.32")
}
