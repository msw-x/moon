package webs

import (
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/msw-x/moon/ulog"
)

type DualServer struct {
	s               *Server
	tls             *Server
	tlsRedirect     string
	tlsAutoRedirect bool
}

func NewDual() *DualServer {
	return &DualServer{
		s:   New(),
		tls: New(),
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

func (o *DualServer) WithLogRequests(use bool) *DualServer {
	o.s.WithLogRequests(use)
	o.tls.WithLogRequests(use)
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

func (o *DualServer) WithXRemoteAddress(s string) *DualServer {
	o.s.WithXRemoteAddress(s)
	o.tls.WithXRemoteAddress(s)
	return o
}

func (o *DualServer) WithRedirectToTls(use string) *DualServer {
	o.tlsRedirect = use
	return o
}

func (o *DualServer) WithAutoRedirectToTls() *DualServer {
	o.tlsAutoRedirect = true
	return o
}

func (o *DualServer) WithTimeout(t Timeout) *DualServer {
	o.s.WithTimeout(t)
	o.tls.WithTimeout(t)
	return o
}

func (o *DualServer) Run(addr string, addrTls string, handler http.Handler) error {
	if !o.tls.IsTls() {
		return errors.New("dual-server: tls secret not defined")
	}
	if o.tlsRedirect == "" && !o.tlsAutoRedirect {
		o.s.Run(addr, handler)
	} else {
		_, port, err := net.SplitHostPort(addrTls)
		if err != nil {
			return nil
		}
		o.s.Run(addr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			host := o.tlsRedirect
			if o.tlsAutoRedirect {
				host = r.Host
			}
			url := fmt.Sprintf("https://%s:%s%s", host, port, r.RequestURI)
			http.Redirect(w, r, url, http.StatusMovedPermanently)
		}))
	}
	o.tls.Run(addrTls, handler)
	return nil
}

func (o *DualServer) Shutdown() {
	o.s.Shutdown()
	o.tls.Shutdown()
}
