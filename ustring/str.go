package ustring

import (
	"strings"
)

func SplitPair(s string, sep string) (string, string) {
	lst := strings.Split(s, sep)
	if len(lst) == 0 {
		return "", ""
	} else if len(lst) == 1 {
		return lst[0], ""
	}
	return lst[0], strings.TrimPrefix(s, lst[0]+sep)
}

func EqualSlices(a, b []string) bool {
	if (a == nil) != (b == nil) {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
