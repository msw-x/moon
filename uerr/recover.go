package uerr

import (
	"fmt"

	"github.com/msw-x/moon/rt"
)

func Recover(onError func(string)) {
	if r := recover(); r != nil {
		if onError != nil {
			onError(fmt.Sprint(r) + "\n" + rt.StackTrace(1))
		}
	}
}
