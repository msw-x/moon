package ufmt

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/msw-x/moon/ustring"
	"golang.org/x/exp/constraints"
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

func (o *IntCtx) init() {
	if o.Base == nil {
		o.base = 1000
	} else {
		o.base = o.Base.(int)
	}
	if o.Precision == nil {
		o.precision = 1
	} else {
		o.precision = o.Precision.(int)
	}
	if o.MinLevel == nil {
		o.minLevel = 0
	} else {
		o.minLevel = o.MinLevel.(int)
	}
	if o.MaxLevel == nil {
		o.maxLevel = len(IntPrefix)
	} else {
		o.maxLevel = o.MaxLevel.(int)
	}
}

func Int[V constraints.Integer](v V, ctx IntCtx) string {
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
		_, fractionalPart := ustring.SplitPair(fmt.Sprintf("%."+strconv.Itoa(ctx.precision)+"f", f), ".")
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
