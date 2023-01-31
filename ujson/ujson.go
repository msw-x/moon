package ujson

import (
	"encoding/json"
	"strings"
)

func MarshalToLowerCamelCase(v any) (dst []byte, err error) {
	dst, err = json.Marshal(v)
	if err == nil {
		return ToLowerCamelCase(dst)
	}
	return
}

func ToLowerCamelCase(src []byte) (dst []byte, err error) {
	var o any
	err = json.Unmarshal(bytes, &o)
	if err == nil {
		o = toLowerCamelCaseMap(o)
	}
	return json.Marshal(o)
}

func toLowerCamelCaseNames(o any) {
	if f, ok := o.(map[string]any); ok {
		for k, v := range f {
			delete(f, k)
			name := toLowerCamelCase(k)
			f[name] = v
			adaptFieldNames(v)
		}
	}

	if f, ok := o.([]any); ok {
		for _, v := range f {
			toLowerCamelCaseNames(v)
		}
	}
}

func toLowerCamelCase(name string) string {
	if name == strings.ToUpper(name) {
		return strings.ToLower(name)
	}
	return strings.ToLower(name[:1]) + name[1:]
}
