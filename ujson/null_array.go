package ujson

import (
	"encoding/json"
	"reflect"
)

func NullArray(bytes []byte) []byte {
	if string(bytes) == "null" {
		bytes = []byte("[]")
	}
	return bytes
}

func InitNilSlice(v any) {
	initNilSlice(reflect.ValueOf(v))
}

func InitNilSliceFlat(i any) {
	v := reflect.ValueOf(i)
	if v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return
		}
		v = v.Elem()
	}
	if v.Kind() == reflect.Slice {
		caseSlice(v)
	}
}

func initNilSlice(v reflect.Value) {
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface:
		if !v.IsNil() {
			initNilSlice(v.Elem())
		}
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			fv := v.Field(i)
			ft := v.Type().Field(i)
			if ft.IsExported() {
				initNilSlice(fv)
			}
		}
	case reflect.Slice:
		caseSlice(v)
		fallthrough
	case reflect.Array:
		for i := 0; i < v.Len(); i++ {
			initNilSlice(v.Index(i))
		}
	case reflect.Map:
		for _, e := range v.MapKeys() {
			initNilSlice(v.MapIndex(e))
		}
	}
}

func isRaw(v reflect.Value) bool {
	var r json.RawMessage
	return v.Type() == reflect.TypeOf(r)
}

func caseSlice(v reflect.Value) {
	if v.IsNil() {
		if v.CanSet() {
			if !isRaw(v) {
				elemType := v.Type().Elem()
				s := reflect.MakeSlice(reflect.SliceOf(elemType), 0, 0)
				v.Set(s)
			}
		}
	}
}
