package urest

import (
	"fmt"
	"net/http"
	"strings"
)

func AuthToken(h http.Header, name string) (token string, present bool, err error) {
	unauthorized := func(s string) {
		err = fmt.Errorf("authorization %s", s)
	}
	auth := h.Get("Authorization")
	if auth == "" {
		unauthorized("is empty")
		return
	}
	present = true
	parts := strings.Split(auth, " ")
	if len(parts) != 2 {
		unauthorized("does not consist of 2 parts")
		return
	}
	if parts[0] != name {
		unauthorized("is not a " + name)
		return
	}
	token = parts[1]
	if token == "" {
		unauthorized("token is empty")
		return
	}
	return
}

func AuthBearer(h http.Header) (string, bool, error) {
	return AuthToken(h, "Bearer")
}
