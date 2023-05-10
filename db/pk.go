package db

import (
	"errors"
	"reflect"
	"strings"

	"github.com/serenize/snaker"
)

func PkName(model any) (name string, err error) {
	fail := func(reason string) {
		err = errors.New("pk name fail: " + reason)
	}
	rt := reflect.TypeOf(model)
	if rt.Kind() != reflect.Pointer {
		fail("model is not pointer")
		return
	}
	rt = rt.Elem()
	if rt.Kind() != reflect.Struct {
		fail("model is not struct")
		return
	}
	name = pkName(rt)
	if name == "" {
		fail("it not found")
	}
	return
}

func pkName(rt reflect.Type) (name string) {
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		if f.Anonymous && f.Type.Kind() == reflect.Struct {
			name = pkName(f.Type)
			if name != "" {
				return
			}
			continue
		}
		s := strings.Split(f.Tag.Get("bun"), ",")
		if len(s) > 1 {
			name = s[0]
			s = s[1:]
			for _, tag := range s {
				if tag == "pk" {
					if name == "" {
						name = snaker.CamelToSnake(f.Name)
					}
					return
				}
			}
		}
	}
	return
}
