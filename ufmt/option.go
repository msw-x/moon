package ufmt

import "fmt"

func Option[T any](ptr *T) string {
	if ptr == nil {
		return "nil"
	}
	return fmt.Sprint(*ptr)
}
