package ujwt

import "github.com/golang-jwt/jwt/v4"

type AccessClaims[Account any, Session any] struct {
	Account Account `json:"account,omitempty"`
	Session Session `json:"session,omitempty"`
	jwt.RegisteredClaims
}

type RefreshClaims[Session any] struct {
	Session Session `json:"session,omitempty"`
	jwt.RegisteredClaims
}
