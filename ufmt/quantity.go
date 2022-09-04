package ufmt

import "github.com/msw-x/moon/umath"

func ByteSize[V umath.AnyInt](v V) string {
	return Int(v, IntCtx{
		Base: 1024,
		Name: "B",
	})
}

func ByteSizeExact[V umath.AnyInt](v V) string {
	return Int(v, IntCtx{
		Base:     1024,
		Name:     "B",
		MaxLevel: 0,
	})
}

func ByteSizeDense[V umath.AnyInt](v V) string {
	return Int(v, IntCtx{
		Base:  1024,
		Name:  "B",
		Dense: true,
	})
}

func ByteSpeed[V umath.AnyInt](v V) string {
	return Int(v, IntCtx{
		Base: 1024,
		Name: "B/s",
	})
}

func ByteSpeedDense[V umath.AnyInt](v V) string {
	return Int(v, IntCtx{
		Base:  1024,
		Name:  "B/s",
		Dense: true,
	})
}

func BitSize[V umath.AnyInt](v V) string {
	return Int(v, IntCtx{
		Base: 1024,
		Name: "b",
	})
}

func BitSizeExact[V umath.AnyInt](v V) string {
	return Int(v, IntCtx{
		Base:     1024,
		Name:     "b",
		MaxLevel: 0,
	})
}

func BitSizeDense[V umath.AnyInt](v V) string {
	return Int(v, IntCtx{
		Base:  1024,
		Name:  "b",
		Dense: true,
	})
}

func BitSpeed[V umath.AnyInt](v V) string {
	return Int(v, IntCtx{
		Base: 1024,
		Name: "b/s",
	})
}

func BitSpeedDense[V umath.AnyInt](v V) string {
	return Int(v, IntCtx{
		Base:  1024,
		Name:  "b/s",
		Dense: true,
	})
}
