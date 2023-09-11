package refl

import "reflect"

func WalkOnStruct(s any, fn func(reflect.Value, reflect.StructField)) {
	FindOnStruct(s, func(v reflect.Value, f reflect.StructField) bool {
		fn(v, f)
		return false
	})
}
