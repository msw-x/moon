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
	type X struct {
		Idx int `bun:",pk"`
	}
	test(new(struct {
		A int `bun:",pk"`
		B string
		C []float32
	}), "a")
	test([]struct {
		A int `bun:",pk"`
		B string
		C []float32
	}{}, "a")
	test(new(struct {
		A int
		B string
		C []float32
		X
	}), "idx")
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
	test(new(struct {
		A int    `bun:"id,pk"`
		B string `bun:"name,pk"`
		C []float32
		X
	}), "id,name,idx")
	test([]struct {
		A int    `bun:"id,pk"`
		B string `bun:"name,pk"`
		C []float32
		X
	}{}, "id,name,idx")

	testErr := func(v any, e string) {
		e = "pk name fail: " + e
		_, r := PkName(v)
		if r == nil {
			t.Errorf("PkName(%+v): without an error", v)
		} else {
			if r.Error() != e {
				t.Errorf("PkName(%+v) = %s; expected %s", v, r, e)
			}
		}
	}
	testErr(5, "model is not pointer")
	testErr(struct {
		A int `bun:"id,pk"`
		B string
	}{}, "model is not pointer")
	testErr(new(struct {
		A int
		B string
	}), "it not found")
}
