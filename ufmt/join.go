package ufmt

import (
	"fmt"
	"strings"
)

func Join(v ...any) string {
	return JoinWith(" ", v...)
}

func JoinWith(splitter string, v ...any) string {
	return JoinSliceWith(splitter, v[:])
}

func JoinSlice[T any](v []T) string {
	return JoinSliceWith(" ", v)
}

func JoinSliceWith[T any](splitter string, v []T) string {
	s := make([]string, len(v))
	for n, a := range v {
		s[n] = fmt.Sprint(a)
	}
	return strings.Join(s, splitter)
}
