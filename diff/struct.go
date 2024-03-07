package diff

import (
	"reflect"
	"strings"
	"time"
)

func Struct[T any](old, new T, tag string) []string {
	refOld := reflect.ValueOf(old)
	refNew := reflect.ValueOf(new)
	if refOld.Type().Kind() != reflect.Struct {
		panic("diff.Struct: object is not struct")
	}
	return impl(refOld, refNew, tag)
}

func impl(refOld, refNew reflect.Value, tag string) []string {
	var names []string
	for i := 0; i < refOld.NumField(); i++ {
		fieldOld := refOld.Field(i)
		fieldNew := refNew.Field(i)
		f := refOld.Type().Field(i)
		if !f.IsExported() {
			continue
		}
		if f.Anonymous && f.Type.Kind() == reflect.Struct {
			names = append(names, impl(fieldOld, fieldNew, tag)...)
			continue
		}
		name := getTag(refOld.Type().Field(i), tag)
		if name == "" {
			name = f.Name
		}
		if fieldOld.Kind() == reflect.Ptr {
			if fieldOld.IsNil() == fieldNew.IsNil() {
				if fieldOld.IsNil() {
					continue
				} else {
					fieldOld = fieldOld.Elem()
					fieldNew = fieldNew.Elem()
				}
			} else {
				names = append(names, name)
				continue
			}
		}
		iold := fieldOld.Interface()
		inew := fieldNew.Interface()
		if fieldOld.Kind() == reflect.Slice {
			if !reflect.DeepEqual(iold, inew) {
				names = append(names, name)
			}
			continue
		}
		if told, ok := iold.(time.Time); ok {
			if tnew, ok := inew.(time.Time); ok {
				if !told.Equal(tnew) {
					names = append(names, name)
				}
			}
		} else {
			if iold != inew {
				names = append(names, name)
			}
		}
	}
	return names
}

func getTag(f reflect.StructField, tag string) string {
	name, ok := f.Tag.Lookup(tag)
	if ok {
		s := strings.Split(name, ",")
		name = s[0]
	}
	return name
}
