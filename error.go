package moon

import (
	"fmt"
	"moon/ufmt"
)

func Panic(v ...any) {
	panic(ufmt.Join(v...))
}

func Panicf(format string, v ...any) {
	panic(fmt.Errorf(format, v...))
}

func Check(err error, v ...any) {
	if err != nil {
		if len(v) > 0 {
			Panicf("%s: %v", ufmt.Join(v...), err)
		} else {
			Panic(err)
		}
	}
}

func Checkf(err error, format string, v ...any) {
	Check(err, fmt.Sprintf(format, v...))
}

func Recover(onError func(string)) {
	if r := recover(); r != nil {
		onError(fmt.Sprint(r))
	}
}
