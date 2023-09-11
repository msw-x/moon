package refl

import "reflect"

func FindField(s any, name string) (val reflect.Value, fld reflect.StructField, ok bool) {
	FindOnStruct(s, func(v reflect.Value, f reflect.StructField) bool {
		ok = f.Name == name
		if ok {
			val = v
			fld = f
		}
		return ok
	})
	return
}
