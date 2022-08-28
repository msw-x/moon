package str

import (
	"strings"
	"unicode"
)

func Trim(s, prefix, suffix string) string {
	return strings.TrimSuffix(strings.TrimPrefix(s, prefix), suffix)
}

func TrimQuotes(s string) string {
	return Trim(s, `"`, `"`)
}

func TrimSquareBrackets(s string) string {
	return Trim(s, "[", "]")
}

func TrimFigureBrackets(s string) string {
	return Trim(s, "{", "}")
}

func TrimBackCarriage(s string) string {
	return strings.TrimSuffix(s, "\r")
}

func TrimBackWhitespaces(s string) string {
	return strings.TrimRightFunc(s, func(r rune) bool {
		return unicode.IsSpace(r)
	})
}

func TrimFrontWhitespaces(s string) string {
	return strings.TrimLeftFunc(s, func(r rune) bool {
		return unicode.IsSpace(r)
	})
}

func TrimWhitespaces(s string) string {
	return TrimFrontWhitespaces(TrimBackWhitespaces(s))
}
