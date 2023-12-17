package uhttp

import (
	"context"
	"crypto/tls"
	"net/http"
	"time"

	"github.com/msw-x/moon/app"
	"github.com/msw-x/moon/secret"
	"github.com/msw-x/moon/ufmt"
	"github.com/msw-x/moon/ulog"
	"golang.org/x/crypto/acme/autocert"
)

type Server struct {
	log       *ulog.Log
	s         *http.Server
	job       *app.Job
	timeouts  timeouts
	certFile  string
	keyFile   string
	domains   []string
	tlsman    *autocert.Manager
	logErrors *ulog.Level
}

func NewServer() *Server {
	return &Server{
		timeouts: timeouts{
			write: 15 * time.Second,
			read:  15 * time.Second,
			idle:  60 * time.Second,
			close: 5 * time.Second,
		},
	}
}

func (o *Server) WithSecret(certFile, keyFile string) *Server {
	o.certFile = certFile
	o.keyFile = keyFile
	return o
}

func (o *Server) WithSecretDir(dir string) *Server {
	o.certFile, o.keyFile = secret.FileNames(dir)
	return o
}

func (o *Server) WithAutoSecret(dir string, domains ...string) *Server {
	o.certFile = ""
	o.keyFile = ""
	o.domains = domains[:]
	if len(o.domains) > 0 && o.domains[0] != "" {
		o.tlsman = &autocert.Manager{
			Cache:      autocert.DirCache(dir),
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(domains...),
		}
	} else {
		o.WithSecretDir(dir)
	}
	return o
}

func (o *Server) WithLogErrors(use bool) *Server {
	if use {
		if o.logErrors == nil {
			o.WithLogErrorsLevel(ulog.LevelError)
		}
	} else {
		o.logErrors = nil
	}
	return o
}

func (o *Server) WithLogErrorsLevel(level ulog.Level) *Server {
	o.logErrors = &level
	return o
}

func (o *Server) Run(addr string, handler http.Handler) {
	if addr == "" {
		return
	}
	name := "http"
	if o.IsTls() {
		name = "https"
		if o.tlsman == nil {
			secret.Ensure(o.certFile, o.keyFile)
		}
	}
	o.log = ulog.New(name).WithID(addr)
	if o.IsTls() {
		if o.tlsman == nil {
			o.log.Info("cert:", o.certFile)
			o.log.Info("key:", o.keyFile)
		} else {
			o.log.Info("domains:", ufmt.JoinSlice(o.domains))
		}
	}
	o.s = &http.Server{
		Addr:         addr,
		Handler:      handler,
		WriteTimeout: o.timeouts.write,
		ReadTimeout:  o.timeouts.read,
		IdleTimeout:  o.timeouts.idle,
		ErrorLog: ulog.StdBridge(func(m string) {
			if o.logErrors != nil {
				level := *o.logErrors
				o.log.Print(level, m)
			}
		}),
	}
	if o.tlsman != nil {
		o.s.TLSConfig = o.tlsman.TLSConfig()
		o.s.TLSConfig.GetCertificate = func(hello *tls.ClientHelloInfo) (cert *tls.Certificate, err error) {
			timeout := 10 * time.Second
			timer := time.AfterFunc(timeout, func() {
				o.log.Warning("getting the certificate takes more than", timeout)
			})
			defer timer.Stop()
			return o.tlsman.GetCertificate(hello)
		}
	}
	o.job = app.NewJob().WithLog(o.log)
	o.job.Run(func() {
		var err error
		if o.IsTls() {
			err = o.s.ListenAndServeTLS(o.certFile, o.keyFile)
		} else {
			err = o.s.ListenAndServe()
		}
		if o.s != nil && err != nil {
			o.log.Error(err)
		}
	})
}

func (o *Server) Close() {
	if o.s != nil {
		ctx, cancel := context.WithTimeout(context.Background(), o.timeouts.close)
		defer cancel()
		s := o.s
		o.s = nil
		s.Shutdown(ctx)
		o.job.Stop()
	}
}

func (o *Server) IsTls() bool {
	return o.certFile != "" && o.keyFile != "" || o.tlsman != nil
}

type timeouts struct {
	write time.Duration
	read  time.Duration
	idle  time.Duration
	close time.Duration
}
