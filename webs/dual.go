package webs

import "net/http"

type DualServer struct {
	s   *Server
	tls *Server
}

func NewDual() *DualServer {
	return &DualServer{
		s:   New(),
		tls: New(),
	}
}

func (this *DualServer) WithSecret(certFile, keyFile string) *DualServer {
	this.tls.WithSecret(certFile, keyFile)
	return this
}

func (this *DualServer) WithSecretDir(dir string) *DualServer {
	this.tls.WithSecretDir(dir)
	return this
}

func (this *DualServer) WithAutoSecret(dir string, domains ...string) *DualServer {
	this.tls.WithAutoSecret(dir, domains...)
	return this
}

func (this *DualServer) Run(addr string, addrTls string, handler http.Handler) {
	if !this.tls.IsTls() {
		panic("dual-server: tls secret not defined")
	}
	this.s.Run(addr, handler)
	this.tls.Run(addrTls, handler)
}

func (this *DualServer) Shutdown() {
	this.s.Shutdown()
	this.tls.Shutdown()
}
