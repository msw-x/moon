package ufmt

import (
	"fmt"
	"strconv"
	"testing"
)

func TestFloat64(t *testing.T) {
	test := func(v float64, p int, e string) {
		r := Float64(v, p)
		if r != e {
			t.Errorf("Float64(%v, %v) = %s; expected %s", v, p, r, e)
		}
	}
	test2 := func(v float64, p int, e string) {
		test(v, p, e)
		if v == 0 && p == 0 {
			test(-v, p, e)
		} else {
			test(-v, p, "-"+e)
		}
		test(v, 0, strconv.Itoa(int(v)))
		test(v, -p, strconv.Itoa(int(v)))
		test(-v, -p, strconv.Itoa(int(-v)))
	}
	for v := 0.; v != 100; v++ {
		for p := 0; p != 20; p++ {
			test2(v, p, fmt.Sprint(v))
		}
	}
	test2(0.0000009, 2, "0")
	test2(93845.0000009, 2, "93845")
	test2(0.0000009, 6, "0.000001")
	test2(3950.0000009, 6, "3950.000001")
	test2(0.000000924, 6, "0.000001")
	test2(384.000000924, 6, "384.000001")
	test2(0.1, 2, "0.1")
	test2(120, 4, "120")
	test2(120, 500, "120")
	test2(120.1749582, 0, "120")
	test2(120.1749582, 1, "120.2")
	test2(120.1749582, 2, "120.17")
	test2(120.1749582, 3, "120.175")
	test2(120.1749582, 4, "120.175")
	test2(120.1749582, 5, "120.17496")
	test2(120.1749582, 6, "120.174958")
	test2(120.1749582, 7, "120.1749582")
	test2(120.1749582, 8, "120.1749582")
}

func TestDelicateFloat64(t *testing.T) {
	test := func(v float64, p int, e string) {
		r := DelicateFloat64(v, p)
		if r != e {
			t.Errorf("DelicateFloat64(%v, %v) = %s; expected %s", v, p, r, e)
		}
	}
	test2 := func(v float64, p int, e string) {
		test(v, p, e)
		test(-v, p, "-"+e)
	}
	test2(120., 1, "120")
	test2(1., 1, "1")
	test2(2., 1, "2")
	test2(7.38900004, 1, "7.4")
	test2(8., 1, "8")
	test2(10., 1, "10")
	test2(10.1, 1, "10.1")
	test2(1.3, 1, "1.3")
	test2(1.356, 1, "1.4")
	test2(9.35637, 1, "9.4")
	test2(0.7, 1, "0.7")
	test2(0.000999191, 1, "0.001")
	test2(2.000999191, 1, "2")
	test2(2.000999191, 4, "2.001")
	test2(38449.000999191, 1, "38449")
	test2(38449.000999191, 3, "38449.001")
	test2(38449.000999191, 5, "38449.001")
	test2(2984858.439000027, 1, "2984858.4")
	test2(4.00000000004, 1, "4")
	test2(0.0000000000091, 1, "0.000000000009")
	test2(0.0000000000091, 4, "0.0000000000091")
	test2(4858639895.34, 1, "4858639895.3")
	test2(4858639895.34, 2, "4858639895.34")
	test2(4858639895.34, 6, "4858639895.34")
}
