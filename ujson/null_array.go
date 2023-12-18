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
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return
		} else {
			v = v.Elem()
		}
	}
	switch v.Kind() {
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			fv := v.Field(i)
			ft := v.Type().Field(i)
			if ft.IsExported() {
				initNilArray(fv)
			}
		}
	case reflect.Slice:
		if v.IsNil() && v.CanSet() {
			elemType := v.Type().Elem()
			s := reflect.MakeSlice(reflect.SliceOf(elemType), 0, 0)
			v.Set(s)
		}
	}
}
