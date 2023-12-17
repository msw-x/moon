package urest

import (
	"errors"
	"net/http"
)

func AuthKey[Account any, Session any](name string, f func(string) (Account, error)) func(http.Header) (Account, Session, bool, error) {
	return func(h http.Header) (account Account, session Session, present bool, err error) {
		key := h.Get(name)
		if key == "" {
			err = errors.New("key is empty")
			return
		}
		present = true
		account, err = f(key)
		return
	}
}
