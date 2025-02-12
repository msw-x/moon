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

func JoinSliceWith[T any](splitter string, l []T) string {
	return JoinSliceFuncWith(splitter, l, func(v T) string {
		return fmt.Sprint(v)
	})
}

func JoinSliceFuncWith[T any](splitter string, l []T, f func(T) string) string {
	s := make([]string, len(l))
	for n, v := range l {
		s[n] = f(v)
	}
	return strings.Join(s, splitter)
}

func NotableJoin(v ...any) string {
	return NotableJoinWith(" ", v...)
}

func NotableJoinWith(splitter string, v ...any) string {
	return NotableJoinSliceWith(splitter, v[:])
}

func NotableJoinSlice[T any](v []T) string {
	return NotableJoinSliceWith(" ", v)
}

func NotableJoinSliceWith[T any](splitter string, v []T) string {
	s := []string{}
	for _, a := range v {
		i := fmt.Sprint(a)
		if i != "" {
			s = append(s, i)
		}
	}
	return strings.Join(s, splitter)
}
