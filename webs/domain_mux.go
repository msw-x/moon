package webs

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/msw-x/moon/ulog"
)

type DomainMux struct {
	m map[string]http.Handler
}

func NewDomainMux() *DomainMux {
	o := new(DomainMux)
	o.m = make(map[string]http.Handler)
	return o
}

func (o *DomainMux) ReverseProxy(domain string, target string) error {
	targetUrl, err := url.Parse(target)
	if err == nil {
		o.m[domain] = httputil.NewSingleHostReverseProxy(targetUrl)
	}
	return err
}

func (o *DomainMux) PureReverseProxy(domain string, target string) error {
	log := ulog.New("reverse-proxy").WithID(domain)
	targetUrl, err := url.Parse(target)
	if err == nil {
		o.m[domain] = &httputil.ReverseProxy{
			Rewrite: func(r *httputil.ProxyRequest) {
				r.SetURL(targetUrl)
			},
			ModifyResponse: func(v *http.Response) error {
				r := v.Request
				if r == nil {
					log.Error("request is nil")
				} else {
					if v.StatusCode == http.StatusOK {
						log.Debug(r.Method, r.URL, v.StatusCode)
					} else {
						log.Error(r.Method, r.URL, v.StatusCode)
					}
				}
				return nil
			},
			ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
				log.Errorf("%s %v: %v", r.Method, r.URL, err)
			},
		}
	}
	return err
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
