package uerr

import (
	"fmt"

	"github.com/msw-x/moon/rt"
	"github.com/msw-x/moon/ufmt"
)

func Strict(err error, v ...any) {
	if err != nil {
		if len(v) > 0 {
			Panicf("%s: %v", ufmt.Join(v...), err)
		} else {
			Panic(err)
		}
	}
}

func Strictf(err error, format string, v ...any) {
	Strict(err, fmt.Sprintf(format, v...))
}

func StrictFn(err error, fn any) {
	Strict(err, rt.FuncName(fn))
}
