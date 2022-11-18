package moon

import (
	"fmt"

	"github.com/msw-x/moon/rt"
	"github.com/msw-x/moon/ufmt"
)

func Panic(v ...any) {
	panic(ufmt.Join(v...))
}

func Panicf(format string, v ...any) {
	panic(fmt.Errorf(format, v...))
}

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

func Recover(onError func(string)) {
	if r := recover(); r != nil {
		onError(fmt.Sprint(r))
	}
}

func StrictFn(err error, fn any) {
	Strict(err, rt.FuncName(fn))
}
