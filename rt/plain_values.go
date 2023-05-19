package rt

import (
	"reflect"
	"strings"
)

// For each plain values of struct
func PlainValues(s any, tagName string, fn func(v any, name string, flags []string)) {
	plainValues(reflect.ValueOf(s), tagName, fn)
}

func plainValues(r reflect.Value, tagName string, fn func(v any, name string, flags []string)) {
	if r.Type().Kind() != reflect.Struct {
		panic("structValues: object is not struct")
	}
	for i := 0; i < r.NumField(); i++ {
		f := r.Type().Field(i)
		if f.Anonymous && f.Type.Kind() == reflect.Struct {
			plainValues(r.Field(i), tagName, fn)
			continue
		}
		v := r.Field(i)
		if f.Type.Kind() == reflect.Ptr {
			if v.IsNil() {
				continue
			}
			v = v.Elem()
		}
		tag := f.Tag.Get(tagName)
		flags := strings.Split(tag, ",")
		var name string
		if len(flags) > 0 {
			name = flags[0]
			flags = flags[1:]
		}
		if name == "" {
			name = f.Name
		}
		fn(v.Interface(), name, flags)
	}
}
