package ujson

import "encoding/json"

func MarshalLowerCase(v any) (bytes []byte, err error) {
	bytes, err = json.Marshal(v)
	if err == nil {
		ToLowerCase(bytes)
	}
	return
}

func ToLowerCase(j []byte) {
	var waitQuote bool
	var waitKey bool
	for i, c := range j {
		if containsAny(c, "\n\n\t ") {
			continue
		}
		if c == '\\' {
			waitQuote = false
			waitKey = false
			continue
		}
		if containsAny(c, "{,") {
			waitQuote = true
			waitKey = true
			continue
		}
		if waitQuote && c == '"' {
			waitQuote = false
			waitKey = true
			continue
		}
		if waitKey {
			if c >= 'A' && c <= 'Z' {
				j[i] = c + 32
			}
		}
		waitQuote = false
		waitKey = false
	}
}

func containsAny(v byte, chars string) bool {
	for _, c := range []byte(chars) {
		if v == c {
			return true
		}
	}
	return false
}
