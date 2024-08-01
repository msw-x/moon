package uhttp

import (
	"net/http"
)

type DomainMux struct {
	m map[string]http.Handler
}

func NewDomainMux() *DomainMux {
	o := new(DomainMux)
	o.m = make(map[string]http.Handler)
	return o
}

func (o *DomainMux) Handler(domain string, h http.Handler) {
	o.m[domain] = h
}

func (o *DomainMux) DefaultHandler(h http.Handler) {
	o.Handler("", h)
}

func (o *DomainMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	host := r.Host
	if o.serve(host, w, r) {
		return
	}
	if o.serve("", w, r) {
		return
	}
	w.WriteHeader(http.StatusForbidden)
	w.Write([]byte("Host forbidden: " + host))
}

func (o *DomainMux) serve(host string, w http.ResponseWriter, r *http.Request) bool {
	handler, ok := o.m[host]
	if ok {
		handler.ServeHTTP(w, r)
	}
	return ok
}
