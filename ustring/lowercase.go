package ustring

import "strings"

func TitleLowerCase(s string) string {
	if len(s) > 0 {
		s = strings.ToLower(string(s[0])) + s[1:]
	}
	return s
}
