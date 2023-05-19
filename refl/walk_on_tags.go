package refl

import (
	"reflect"
	"strings"
)

func WalkOnTags(s any, tagName string, fn func(v reflect.Value, name string, flags []string)) {
	WalkOnStruct(s, func(v reflect.Value, f reflect.StructField) {
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
		fn(v, name, flags)
	})
}

func WalkOnTagsAny(s any, tagName string, fn func(v any, name string, flags []string)) {
	WalkOnTags(s, tagName, func(v reflect.Value, name string, flags []string) {
		if v.Kind() == reflect.Ptr {
			if v.IsNil() {
				return
			}
			v = v.Elem()
		}
		fn(v.Interface(), name, flags)
	})
}
