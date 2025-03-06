package ustring

import (
	"strings"
	"unicode"
)

// https://pkg.go.dev/golang.org/x/text/cases

func TitleLowerCase(s string) string {
	if len(s) > 0 {
		s = strings.ToLower(string(s[0])) + s[1:]
	}
	return s
}

func IsUpper(s string) bool {
	for _, r := range s {
		if !unicode.IsUpper(r) && unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func IsLower(s string) bool {
	for _, r := range s {
		if !unicode.IsLower(r) && unicode.IsLetter(r) {
			return false
		}
	}
	return true
}
