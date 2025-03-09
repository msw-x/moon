package db

import "testing"

func TestPk(t *testing.T) {
	test := func(v any, e string) {
		r, err := PkName(v)
		if err == nil {
			if r != e {
				t.Errorf("PkName(%+v) = %s; expected %s", v, r, e)
			}
		} else {
			t.Errorf("PkName(%+v): %v", v, err)
		}
	}
	test(new(struct {
		A int `bun:",pk"`
		B string
		C []float32
	}), "a")
	test(new(struct {
		A int `bun:"name,pk"`
		B string
		C []float32
	}), "name")
	test(new(struct {
		A int    `bun:"id,pk"`
		B string `bun:"name,pk"`
		C []float32
	}), "id,name")
}
