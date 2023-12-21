package urest

import (
	"errors"
	"net/http"

	"github.com/msw-x/moon/uhttp"
)

type AuthContext[Account any, Session any] struct {
	Base *Context
	auth []func(h http.Header) (Account, Session, bool, error)
}

func NewAuthContext[Account any, Session any](base *Context, auth func(h http.Header) (Account, Session, bool, error)) *AuthContext[Account, Session] {
	return (&AuthContext[Account, Session]{
		Base: base,
	}).AppendAuth(auth)
}

func (o AuthContext[Account, Session]) Router() *uhttp.Router {
	return o.Base.Router()
}

func (o AuthContext[Account, Session]) Branch(name string) *AuthContext[Account, Session] {
	o.Base = o.Base.Branch(name)
	return &o
}

func (o *AuthContext[Account, Session]) AppendAuth(auth func(h http.Header) (Account, Session, bool, error)) *AuthContext[Account, Session] {
	o.auth = append(o.auth, auth)
	return o
}

func (o *AuthContext[Account, Session]) Auth(h http.Header) (account Account, session Session, err error) {
	if len(o.auth) == 0 {
		err = errors.New("authorization is unavailable")
	} else {
		for _, f := range o.auth {
			var present bool
			account, session, present, err = f(h)
			if present {
				return
			}
		}
		err = errors.New("authorization is missing")
	}
	return
}
