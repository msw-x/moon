package uhttp

import (
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/msw-x/moon/ulog"
)

type DualServer struct {
	s           *Server
	tls         *Server
	tlsRedirect string
}

func NewDualServer() *DualServer {
	return &DualServer{
		s:   NewServer(),
		tls: NewServer(),
	}
}

func (o *DualServer) WithSecret(certFile, keyFile string) *DualServer {
	o.tls.WithSecret(certFile, keyFile)
	return o
}

func (o *DualServer) WithSecretDir(dir string) *DualServer {
	o.tls.WithSecretDir(dir)
	return o
}

func (o *DualServer) WithAutoSecret(dir string, domains ...string) *DualServer {
	o.tls.WithAutoSecret(dir, domains...)
	return o
}

func (o *DualServer) WithLogErrors(use bool) *DualServer {
	o.s.WithLogErrors(use)
	o.tls.WithLogErrors(use)
	return o
}

func (o *DualServer) WithLogErrorsLevel(level ulog.Level) *DualServer {
	o.s.WithLogErrorsLevel(level)
	o.tls.WithLogErrorsLevel(level)
	return o
}

func (o *DualServer) WithRedirectToTls(use string) *DualServer {
	o.tlsRedirect = use
	return o
}

func (o *DualServer) Run(addr string, addrTls string, handler http.Handler) error {
	if !o.tls.IsTls() {
		return errors.New("dual-server: tls secret not defined")
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

func (o *DualServer) Close() {
	o.s.Close()
	o.tls.Close()
}
