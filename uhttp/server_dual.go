package uhttp

import (
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/msw-x/moon/ulog"
)

type ServerDual struct {
	s           *Server
	tls         *Server
	tlsRedirect string
}

func ServerNewDual() *ServerDual {
	return &ServerDual{
		s:   NewServer(),
		tls: NewServer(),
	}
}

func (o *ServerDual) WithSecret(certFile, keyFile string) *ServerDual {
	o.tls.WithSecret(certFile, keyFile)
	return o
}

func (o *ServerDual) WithSecretDir(dir string) *ServerDual {
	o.tls.WithSecretDir(dir)
	return o
}

func (o *ServerDual) WithAutoSecret(dir string, domains ...string) *ServerDual {
	o.tls.WithAutoSecret(dir, domains...)
	return o
}

func (o *ServerDual) WithLogErrors(use bool) *ServerDual {
	o.s.WithLogErrors(use)
	o.tls.WithLogErrors(use)
	return o
}

func (o *ServerDual) WithLogErrorsLevel(level ulog.Level) *ServerDual {
	o.s.WithLogErrorsLevel(level)
	o.tls.WithLogErrorsLevel(level)
	return o
}

func (o *ServerDual) WithRedirectToTls(use string) *ServerDual {
	o.tlsRedirect = use
	return o
}

func (o *ServerDual) Run(addr string, addrTls string, handler http.Handler) error {
	if !o.tls.IsTls() {
		return errors.New("server-dual: tls secret not defined")
	}
	if o.tlsRedirect == "" {
		o.s.Run(addr, handler)
	} else {
		_, port, err := net.SplitHostPort(addrTls)
		if err != nil {
			return nil
		}
		o.s.Run(addr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			url := fmt.Sprintf("https://%s:%s%s", o.tlsRedirect, port, r.RequestURI)
			http.Redirect(w, r, url, http.StatusMovedPermanently)
		}))
	}
	o.tls.Run(addrTls, handler)
	return nil
}

func (o *ServerDual) Close() {
	o.s.Close()
	o.tls.Close()
}
