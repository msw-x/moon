package ujson

import (
	"encoding/json"
	"slices"
)

func MarshalLowerCase(v any) (bytes []byte, err error) {
	bytes, err = json.Marshal(v)
	if err == nil {
		ToLowerCase(bytes)
	}
	return
}

func ToLowerCase(j []byte) {
	var keyIndex int
	var waitColon bool
	for i, c := range j {
		if containsAny(c, "\n\n\t ") {
			continue
		}
		if containsAny(c, "{}[],\\-_") {
			keyIndex = 0
			waitColon = false
			continue
		}
		if c == '"' {
			if keyIndex == 0 {
				keyIndex = i + 1
			} else {
				waitColon = true
			}
			continue
		}
		if waitColon && c == ':' {
			c = j[keyIndex]
			if c >= 'A' && c <= 'Z' {
				j[keyIndex] = c + 32
			}
			keyIndex = 0
			waitColon = false
		}
	}
}

func containsAny(v byte, chars string) bool {
	return slices.Contains([]byte(chars), v)
}
