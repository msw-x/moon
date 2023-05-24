package uhttp

import "fmt"

func IsEmpty(s string) bool {
	return s == "" || s == "0"
}

func Marshal(v any, omitempty bool) (s string, omit bool) {
	s = fmt.Sprint(v)
	omit = omitempty && IsEmpty(s)
	return
}
