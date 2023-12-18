package ujson

import (
	"reflect"
)

func NullArray(bytes []byte) []byte {
	if string(bytes) == "null" {
		bytes = []byte("[]")
	}
	return bytes
}

func InitNilArray(v any) {
	initNilArray(reflect.ValueOf(v))
}

func initNilArray(v reflect.Value) {
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface:
		if !v.IsNil() {
			initNilArray(v.Elem())
		}
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			fv := v.Field(i)
			ft := v.Type().Field(i)
			if ft.IsExported() {
				initNilArray(fv)
			}
		}
	case reflect.Slice:
		if v.IsNil() {
			if v.CanSet() {
				elemType := v.Type().Elem()
				s := reflect.MakeSlice(reflect.SliceOf(elemType), 0, 0)
				v.Set(s)
			}
		}
		fallthrough
	case reflect.Array:
		for i := 0; i < v.Len(); i++ {
			initNilArray(v.Index(i))
		}
	case reflect.Map:
		for _, e := range v.MapKeys() {
			initNilArray(v.MapIndex(e))
		}
	}
}
