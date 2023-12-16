package ujwt

import "github.com/golang-jwt/jwt/v4"

type AccessClaims[Account any, SessionId any] struct {
	Account   Account   `json:"account"`
	SessionId SessionId `json:"sessionId"`
	jwt.RegisteredClaims
}

type RefreshClaims[SessionId any] struct {
	SessionId SessionId `json:"sessionId"`
	jwt.RegisteredClaims
}
