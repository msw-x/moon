package webs

import (
	"context"
	"net/http"
	"time"

	"github.com/msw-x/moon/app"
	"github.com/msw-x/moon/secret"
	"github.com/msw-x/moon/ulog"
	"github.com/msw-x/moon/usync"
	"golang.org/x/crypto/acme/autocert"
)

type Server struct {
	log      *ulog.Log
	s        *http.Server
	do       *usync.Do
	timeout  timeout
	certFile string
	keyFile  string
	tlsman   *autocert.Manager
}

func New() *Server {
	return &Server{
		timeout: timeout{
			write: 15 * time.Second,
			read:  15 * time.Second,
			idle:  60 * time.Second,
			close: 5 * time.Second,
		},
	}
}

func (this *Server) WithSecret(certFile, keyFile string) *Server {
	this.certFile = certFile
	this.keyFile = keyFile
	return this
}

func (this *Server) WithSecretDir(dir string) *Server {
	this.certFile, this.keyFile = secret.FileNames(dir)
	return this
}

func (this *Server) WithAutoSecret(dir string, domains ...string) *Server {
	this.certFile = ""
	this.keyFile = ""
	this.tlsman = &autocert.Manager{
		Cache:      autocert.DirCache(dir),
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(domains...),
	}
	return this
}

func (this *Server) Run(addr string, handler http.Handler) {
	if addr == "" {
		return
	}
	name := "http"
	if this.IsTls() {
		name = "https"
		if this.tlsman == nil {
			secret.Ensure(this.certFile, this.keyFile)
		}
	}
	this.log = ulog.New(name).WithID(addr)
	if this.IsTls() {
		if this.tlsman == nil {
			this.log.Info("cert:", this.certFile)
			this.log.Info("key:", this.keyFile)
		}
	}
	this.s = &http.Server{
		Addr:         addr,
		Handler:      handler,
		WriteTimeout: this.timeout.write,
		ReadTimeout:  this.timeout.read,
		IdleTimeout:  this.timeout.idle,
	}
	if this.tlsman != nil {
		this.s.TLSConfig = this.tlsman.TLSConfig()
	}
	this.do = usync.NewDo()
	app.Go(func() {
		defer func() {
			this.log.Info("stopped")
			this.do.Notify()
		}()
		this.log.Info("listen")
		var err error
		if this.IsTls() {
			err = this.s.ListenAndServeTLS(this.certFile, this.keyFile)
		} else {
			err = this.s.ListenAndServe()
		}
		if this.s != nil && err != nil {
			this.log.Error(err)
		}
	})
}

func (this *Server) Shutdown() {
	if this.s != nil {
		ctx, cancel := context.WithTimeout(context.Background(), this.timeout.close)
		defer cancel()
		this.log.Info("shutdown")
		s := this.s
		this.s = nil
		s.Shutdown(ctx)
		this.do.Stop()
		this.log.Info("shutdown completed")
	}
}

func (this *Server) IsTls() bool {
	return this.certFile != "" && this.keyFile != "" || this.tlsman != nil
}

type timeout struct {
	write time.Duration
	read  time.Duration
	idle  time.Duration
	close time.Duration
}
