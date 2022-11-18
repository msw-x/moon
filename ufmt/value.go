package ufmt

import (
	"fmt"
	"reflect"
)

func Value(v any) string {
	return reflectValue(reflect.ValueOf(v))
}

func reflectValue(v reflect.Value) string {
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return "nil"
		} else {
			v = v.Elem()
		}
	}
	var l []string
	var f string
	put := func(s string) {
		l = append(l, s)
	}
	switch v.Kind() {
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			fv := v.Field(i)
			ft := v.Type().Field(i)
			if ft.IsExported() {
				put(fmt.Sprintf("%s:%s", ft.Name, reflectValue(fv)))
			}
		}
		f = "{%s}"
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			put(reflectValue(v.Index(i)))
		}
		f = "[%s]"
	case reflect.Map:
		for _, e := range v.MapKeys() {
			put(fmt.Sprintf("%s:%s", reflectValue(e), reflectValue(v.MapIndex(e))))
		}
		f = "[%s]"
	default:
		return fmt.Sprint(v.Interface())
	}
	return fmt.Sprintf(f, JoinSlice(l))
}
