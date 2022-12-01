package ufmt

import "golang.org/x/exp/constraints"

func ByteSize[V constraints.Integer](v V) string {
	return Int(v, IntCtx{
		Base: 1024,
		Name: "B",
	})
}

func ByteSizeExact[V constraints.Integer](v V) string {
	return Int(v, IntCtx{
		Base:     1024,
		Name:     "B",
		MaxLevel: 0,
	})
}

func ByteSizeDense[V constraints.Integer](v V) string {
	return Int(v, IntCtx{
		Base:  1024,
		Name:  "B",
		Dense: true,
	})
}

func ByteSpeed[V constraints.Integer](v V) string {
	return Int(v, IntCtx{
		Base: 1024,
		Name: "B/s",
	})
}

func ByteSpeedDense[V constraints.Integer](v V) string {
	return Int(v, IntCtx{
		Base:  1024,
		Name:  "B/s",
		Dense: true,
	})
}

func BitSize[V constraints.Integer](v V) string {
	return Int(v, IntCtx{
		Base: 1024,
		Name: "b",
	})
}

func BitSizeExact[V constraints.Integer](v V) string {
	return Int(v, IntCtx{
		Base:     1024,
		Name:     "b",
		MaxLevel: 0,
	})
}

func BitSizeDense[V constraints.Integer](v V) string {
	return Int(v, IntCtx{
		Base:  1024,
		Name:  "b",
		Dense: true,
	})
}

func BitSpeed[V constraints.Integer](v V) string {
	return Int(v, IntCtx{
		Base: 1024,
		Name: "b/s",
	})
}

func BitSpeedDense[V constraints.Integer](v V) string {
	return Int(v, IntCtx{
		Base:  1024,
		Name:  "b/s",
		Dense: true,
	})
}
