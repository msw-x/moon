package refl

import "reflect"

func Apply(src any, dst any) {
	WalkOnStruct(src, func(v reflect.Value, f reflect.StructField) {
		if f.Type.Kind() == reflect.Pointer {
			if !v.IsNil() {
				if d, _, found := FindField(dst, f.Name); found {
					if d.CanSet() {
						d.Set(v.Elem())
					}
				}
			}
		}
	})
}
