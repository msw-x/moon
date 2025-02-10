package uerr

import (
	"errors"
	"fmt"
	"reflect"
)

func Unwind(err error) (s string) {
	for err != nil {
		t := reflect.TypeOf(err)
		n := "?"
		if t != nil {
			n = t.String()
		}
		if s != "" {
			s += "\n"
		}
		s += fmt.Sprintf("[%s] %v", n, err)
		err = errors.Unwrap(err)
	}
	return
}
