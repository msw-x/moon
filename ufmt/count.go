package ufmt

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/msw-x/moon/str"
	"github.com/msw-x/moon/umath"
)

const IntPrefix = "KMGTPE"

type IntCtx struct {
	Base          any
	Precision     any
	FixFractional bool
	MinLevel      any
	MaxLevel      any
	Dense         bool
	Name          string

	base      int
	precision int
	minLevel  int
	maxLevel  int
}

func (this *IntCtx) init() {
	if this.Base == nil {
		this.base = 1000
	}
	this.base = this.Base.(int)
	if this.Precision == nil {
		this.precision = 1
	}
	this.precision = this.Precision.(int)
	if this.MinLevel == nil {
		this.minLevel = 0
	}
	this.minLevel = this.MinLevel.(int)
	if this.MaxLevel == nil {
		this.maxLevel = len(IntPrefix)
	}
	this.maxLevel = this.MaxLevel.(int)
}

func Int[V umath.AnyInt](v V, ctx IntCtx) string {
	ctx.init()
	n := v
	base := V(ctx.base)
	level := 0
	for level < ctx.minLevel {
		n /= base
		level++
	}
	for n >= base && level < ctx.maxLevel {
		n /= base
		level++
	}
	s := ""
	if ctx.Dense {
		s = fmt.Sprintf("%d", n)
	} else {
		s = WideInt(n)
	}
	if ctx.precision > 0 {
		div := V(math.Pow(float64(base), float64(level)))
		f := float64(v) / float64(div)
		_, fractionalPart := str.SplitPair(fmt.Sprintf("%."+strconv.Itoa(ctx.precision)+"f", f), ".")
		if !ctx.FixFractional {
			for strings.HasSuffix(fractionalPart, "0") {
				fractionalPart = strings.TrimSuffix(fractionalPart, "0")
			}
		}
		if len(fractionalPart) > 0 {
			s = fmt.Sprintf("%s.%s", s, fractionalPart)
		}
	}
	prefix := ""
	if level > 0 {
		prefix = fmt.Sprintf("%c", IntPrefix[level-1])
	}
	space := ""
	if !ctx.Dense {
		space = " "
	}
	return fmt.Sprintf("%s%s%s%s", s, space, prefix, ctx.Name)
}

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
