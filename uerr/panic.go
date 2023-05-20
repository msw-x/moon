package uerr

import (
	"fmt"

	"github.com/msw-x/moon/ufmt"
)

func Panic(v ...any) {
	panic(ufmt.Join(v...))
}

func Panicf(format string, v ...any) {
	panic(fmt.Errorf(format, v...))
}
