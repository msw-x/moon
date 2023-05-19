package refl

import "reflect"

func WalkByPlainStruct(s any, tagName string, fn func(reflect.StructField)) {
	walkByPlainStruct(reflect.ValueOf(s), tagName, fn)
}

func walkByPlainStruct(r reflect.Value, tagName string, fn func(reflect.StructField)) {
	if r.Type().Kind() != reflect.Struct {
		panic("structValues: object is not struct")
	}
	for i := 0; i < r.NumField(); i++ {
		f := r.Type().Field(i)
		if f.Anonymous && f.Type.Kind() == reflect.Struct {
			walkByPlainStruct(r.Field(i), tagName, fn)
			continue
		}
		fn(r.Field(i))
	}
}