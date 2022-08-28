package ufmt

import (
	"fmt"
	"strings"
)

func JoinWith(splitter string, v ...any) string {
	s := make([]string, len(v))
	for n, a := range v {
		s[n] = fmt.Sprint(a)
	}
	return strings.Join(s, splitter)
}

func Join(v ...any) string {
	return JoinWith(" ", v...)
}
