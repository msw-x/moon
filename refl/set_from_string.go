package refl

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/msw-x/moon/parse"
)

func SetFromString(v reflect.Value, s string) error {
	if !v.CanSet() {
		return errors.New("SetFromString: can't set")
	}
	switch v.Kind() {
	case reflect.Bool:
		x, err := parse.Bool(s)
		if err != nil {
			return err
		}
		v.SetBool(x)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		x, err := parse.Int64(s)
		if err != nil {
			return err
		}
		v.SetInt(x)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		x, err := parse.Uint64(s)
		if err != nil {
			return err
		}
		v.SetUint(x)
	case reflect.Float32, reflect.Float64:
		x, err := parse.Float64(s)
		if err != nil {
			return err
		}
		v.SetFloat(x)
	case reflect.String:
		v.SetString(s)
	default:
		return fmt.Errorf("SetFromString: invalid type: %s", v.Type().Name())
	}
	return nil
}
