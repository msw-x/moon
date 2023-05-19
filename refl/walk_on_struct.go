package refl

import "reflect"

func WalkOnStruct(s any, fn func(reflect.Value, reflect.StructField)) {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return
		}
		v = v.Elem()
	}
	walkOnStruct(v, fn)
}

func walkOnStruct(r reflect.Value, fn func(reflect.Value, reflect.StructField)) {
	if r.Type().Kind() != reflect.Struct {
		panic("structValues: object is not struct")
	}
	for i := 0; i < r.NumField(); i++ {
		f := r.Type().Field(i)
		if f.Anonymous && f.Type.Kind() == reflect.Struct {
			walkOnStruct(r.Field(i), fn)
			continue
		}
		fn(r.Field(i), f)
	}
}
