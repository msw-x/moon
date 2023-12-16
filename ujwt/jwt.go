package ujwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Jwt[Account any, SessionId any] struct {
	opts Options
}

func New[Account any, SessionId any](opts Options) *Jwt {
	return &Jwt[Account, SessionId]{
		opts: opts,
	}
}

func (o *Jwt[Account, SessionId]) NewAccess(account Account, sessionId SessionId) (string, error) {
	return o.new(AccessClaims{
		RegisteredClaims: registeredClaims(o.opts.ExpirationTime),
		Account:          account,
		SessionId:        sessionId,
	})
}

func (o *Jwt[Account, SessionId]) NewRefresh(sessionId SessionId) (string, error) {
	return o.new(RefreshClaims{
		RegisteredClaims: registeredClaims(o.opts.RefreshExpirationTime),
		SessionId:        sessionId,
	})
}

func (o *Jwt[Account, SessionId]) New(account Account, sessionId SessionId) (token string, refreshToken string, err error) {
	token, err = o.NewAccess(account, sessionId)
	if err == nil {
		refreshToken, err = o.NewRefresh(sessionId)
	}
	return
}

func (o *Jwt[Account, SessionId]) ParseAccess(token string) (claims AccessClaims, err error) {
	err = o.parse(token, &claims)
	return
}

func (o *Jwt[Account, SessionId]) ParseRefresh(token string) (claims RefreshClaims, err error) {
	err = o.parse(token, &claims)
	return
}

func (o *Jwt[Account, SessionId]) ParseRefreshSessionId(token string) (SessionId, error) {
	claims, err := o.ParseRefresh(token)
	return claims.SessionId, err
}

func (o *Jwt[Account, SessionId]) Options() Options {
	return o.opts
}

func (o *Jwt[Account, SessionId]) new(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(o.opts.KeyBytes())
}

func (o *Jwt[Account, SessionId]) parse(token string, claims jwt.Claims) (err error) {
	err = emptyRestrict(token)
	if err == nil {
		_, err = jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("jwt unexpected signing method: %v", t.Header["alg"])
			}
			return o.opts.KeyBytes(), nil
		})
	}
	return
}

func registeredClaims(expirationTime time.Duration) jwt.RegisteredClaims {
	return jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expirationTime)),
	}
}

func emptyRestrict(token string) error {
	if token == "" || token == "null" {
		return errors.New("jwt token is empty")
	}
	return nil
}
