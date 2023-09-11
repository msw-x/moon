package refl

import "reflect"

func FindOnStruct(s any, fn func(reflect.Value, reflect.StructField) bool) bool {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return false
		}
		v = v.Elem()
	}
	return findOnStruct(v, fn)
}

func findOnStruct(r reflect.Value, fn func(reflect.Value, reflect.StructField) bool) (found bool) {
	if r.Type().Kind() != reflect.Struct {
		panic("findOnStruct: object is not struct")
	}
	for i := 0; i < r.NumField(); i++ {
		if found {
			break
		}
		f := r.Type().Field(i)
		if f.Anonymous && f.Type.Kind() == reflect.Struct {
			found = findOnStruct(r.Field(i), fn)
			continue
		}
		found = fn(r.Field(i), f)
	}
	return
}
