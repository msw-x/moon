package webs

import (
	"context"
	"crypto/tls"
	"net/http"
	"strings"
	"time"

	"github.com/msw-x/moon/app"
	"github.com/msw-x/moon/secret"
	"github.com/msw-x/moon/ufmt"
	"github.com/msw-x/moon/ulog"
	"github.com/msw-x/moon/usync"
	"golang.org/x/crypto/acme/autocert"
)

type Server struct {
	log            *ulog.Log
	s              *http.Server
	do             *usync.Do
	timeout        Timeout
	certFile       string
	keyFile        string
	domains        []string
	tlsman         *autocert.Manager
	xRemoteAddress string
	logRequests    bool
	logErrors      *ulog.Level
}

func New() *Server {
	return &Server{
		timeout: Timeout{
			Write: 15 * time.Second,
			Read:  15 * time.Second,
			Idle:  60 * time.Second,
			Close: 5 * time.Second,
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
	o.domains = splitDomains(domains[:])
	if len(o.domains) > 0 && o.domains[0] != "" {
		o.tlsman = &autocert.Manager{
			Cache:      autocert.DirCache(dir),
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(o.domains...),
		}
	} else {
		o.WithSecretDir(dir)
	}
	return o
}

func (o *Server) WithLogRequests(use bool) *Server {
	o.logRequests = use
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

func (o *Server) WithXRemoteAddress(s string) *Server {
	o.xRemoteAddress = s
	return o
}

func (o *Server) WithTimeout(timeout Timeout) *Server {
	o.timeout = timeout
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
	if o.logRequests {
		handler = o.LogRequest(handler)
	}
	o.s = &http.Server{
		Addr:         addr,
		Handler:      handler,
		WriteTimeout: o.timeout.Write,
		ReadTimeout:  o.timeout.Read,
		IdleTimeout:  o.timeout.Idle,
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
	o.do = usync.NewDo()
	app.Go(func() {
		defer func() {
			o.log.Info("stopped")
			o.do.Notify()
		}()
		o.log.Info("listen")
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

func (o *Server) Shutdown() {
	if o.s != nil {
		ctx, cancel := context.WithTimeout(context.Background(), o.timeout.Close)
		defer cancel()
		o.log.Info("shutdown")
		s := o.s
		o.s = nil
		s.Shutdown(ctx)
		o.do.Stop()
		o.log.Info("shutdown completed")
	}
}

func (o *Server) IsTls() bool {
	return o.certFile != "" && o.keyFile != "" || o.tlsman != nil
}

func (o *Server) LogRequest(mux http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := RequestNameX(r, o.xRemoteAddress)
		o.log.Debug(name)
		tm := time.Now()
		mux.ServeHTTP(w, r)
		o.log.Debug(name, time.Since(tm).Truncate(time.Millisecond))
	})
}

func splitDomains(l []string) (r []string) {
	for _, s := range l {
		r = append(r, strings.Split(s, " ")...)
	}
	return
}
