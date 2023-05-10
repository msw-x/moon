package db

import "github.com/msw-x/moon/diff"

func Diff[T any](a, b T) []string {
	return diff.Struct(a, b, "bun")
}
