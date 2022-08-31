package ufmt

import "github.com/msw-x/moon/umath"

func ByteSize[V umath.AnyInt](v V) string {
	return Int(v, IntCtx{
		Name: "B",
	})
}

func ByteSizeDense[V umath.AnyInt](v V) string {
	return Int(v, IntCtx{
		Dense: true,
		Name:  "B",
	})
}

func ByteSpeed[V umath.AnyInt](v V) string {
	return Int(v, IntCtx{
		Name: "B/s",
	})
}

func ByteSpeedDense[V umath.AnyInt](v V) string {
	return Int(v, IntCtx{
		Dense: true,
		Name:  "B/s",
	})
}

func BitSize[V umath.AnyInt](v V) string {
	return Int(v, IntCtx{
		Name: "b",
	})
}

func BitSizeDense[V umath.AnyInt](v V) string {
	return Int(v, IntCtx{
		Dense: true,
		Name:  "b",
	})
}

func BitSpeed[V umath.AnyInt](v V) string {
	return Int(v, IntCtx{
		Name: "b/s",
	})
}

func BitSpeedDense[V umath.AnyInt](v V) string {
	return Int(v, IntCtx{
		Dense: true,
		Name:  "b/s",
	})
}
