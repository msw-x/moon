package ufmt

import "testing"

func TestInt(t *testing.T) {
	test := func(v int, e string) {
		p := "IntCtx{Precision: 0, Dense: true}"
		r := Int(v, IntCtx{Precision: 0, Dense: true})
		if r != e {
			t.Errorf("Float64(%v, %v) = %s; expected %s", v, p, r, e)
		}
	}
	test(0, "0")
	test(2, "2")
	test(17, "17")
	test(96, "96")
	test(99, "99")
	test(100, "100")
	test(707, "707")
	test(999, "999")
	test(1000, "1K")
	test(1001, "1K")
	test(1200, "1K")
	test(1900, "1K")
	test(1999, "1K")
	test(2000, "2K")
}
