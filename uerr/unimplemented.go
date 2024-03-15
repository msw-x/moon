package uerr

import "errors"

func Unimplemented() error {
	return errors.New("unimplemented")
}
