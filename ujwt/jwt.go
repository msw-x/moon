package ujwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Jwt[Account any, Session any] struct {
	opts Options
}

func New[Account any, Session any](opts Options) *Jwt[Account, Session] {
	return &Jwt[Account, Session]{
		opts: opts,
	}
}

func (o *Jwt[Account, Session]) NewAccess(account Account, session Session) (string, error) {
	return o.new(AccessClaims[Account, Session]{
		RegisteredClaims: registeredClaims(o.opts.ExpirationTime),
		Account:          account,
		Session:          session,
	})
}

func (o *Jwt[Account, Session]) NewRefresh(session Session) (string, error) {
	return o.new(RefreshClaims[Session]{
		RegisteredClaims: registeredClaims(o.opts.RefreshExpirationTime),
		Session:          session,
	})
}

func (o *Jwt[Account, Session]) New(account Account, session Session) (token string, refreshToken string, err error) {
	token, err = o.NewAccess(account, session)
	if err == nil {
		refreshToken, err = o.NewRefresh(session)
	}
	return
}

func (o *Jwt[Account, Session]) ParseAccess(token string) (claims AccessClaims[Account, Session], err error) {
	err = o.parse(token, &claims)
	return
}

func (o *Jwt[Account, Session]) ParseRefresh(token string) (claims RefreshClaims[Session], err error) {
	err = o.parse(token, &claims)
	return
}

func (o *Jwt[Account, Session]) ParseRefreshSession(token string) (Session, error) {
	claims, err := o.ParseRefresh(token)
	return claims.Session, err
}

func (o *Jwt[Account, Session]) Options() Options {
	return o.opts
}

func (o *Jwt[Account, Session]) new(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(o.opts.KeyBytes())
}

func (o *Jwt[Account, Session]) parse(token string, claims jwt.Claims) (err error) {
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
