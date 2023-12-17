package urest

import (
	"fmt"
	"net/http"

	"github.com/msw-x/moon/ujwt"
)

func AuthJwt[Account any, Session any](jwt *ujwt.Jwt[Account, Session], validate func(Session, string) error) func(http.Header) (Account, Session, bool, error) {
	return func(h http.Header) (account Account, session Session, present bool, err error) {
		var token string
		token, present, err = AuthBearer(h)
		if err == nil {
			var claims ujwt.AccessClaims[Account, Session]
			if claims, err = jwt.ParseAccess(token); err == nil {
				account = claims.Account
				session = claims.Session
				if validate != nil {
					err = validate(session, token)
				}
			} else {
				err = fmt.Errorf("jwt is invalid: %v", err)
			}
		}
		return
	}
}
