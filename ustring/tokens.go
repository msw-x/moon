package ustring

import "strings"

func TransformTokens(s string, sep string, fn func(string) string) string {
	lines := strings.Split(s, sep)
	for n, line := range lines {
		lines[n] = fn(line)
	}
	s = strings.Join(lines, sep)
	return s
}

func TransformLines(s string, fn func(string) string) string {
	return TransformTokens(s, "\n", fn)
}
