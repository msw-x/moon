package ujson

import (
	"fmt"
	"strings"
)

func quote(v any) []byte {
	return []byte(fmt.Sprintf(`"%v"`, v))
}

func unquote(b []byte) string {
	return strings.Trim(string(b), `"`)
}
