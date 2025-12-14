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
	return o.PureReverseProxyFilter(domain, target, nil)
}

func (o *DomainMux) PureReverseProxyFilter(domain string, target string, f func(http.ResponseWriter, *http.Request, *ulog.Log) bool) error {
	log := ulog.New("reverse-proxy").WithID(domain)
	targetUrl, err := url.Parse(target)
	if err == nil {
		p := &httputil.ReverseProxy{
			Rewrite: func(r *httputil.ProxyRequest) {
				r.SetURL(targetUrl)
				var body []byte
				var err error
				//body, err = io.ReadAll(r.In.Body)
				//if err == nil {
				//	r.In.Body.Close()
				//	r.In.Body = io.NopCloser(bytes.NewReader(body))
				//}
				log.Debug("request", r.In.RemoteAddr, r.In.Method, r.In.URL, len(body), "B", "=>", r.Out.URL, err)
			},
			ModifyResponse: func(v *http.Response) error {
				r := v.Request
				if r == nil {
					log.Error("request is nil")
				} else {
					if v.StatusCode == http.StatusOK {
						var body []byte
						var err error
						//body, err = io.ReadAll(r.Body)
						//if err == nil {
						//	r.Body.Close()
						//	r.Body = io.NopCloser(bytes.NewReader(body))
						//}
						log.Debug("response", r.RemoteAddr, r.Method, r.URL, v.StatusCode, len(body), "B", err)
					} else {
						log.Error("response", r.RemoteAddr, r.Method, r.URL, v.StatusCode)
					}
				}
				return nil
			},
			ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
				log.Errorf("%s %s %v: %v", r.RemoteAddr, r.Method, r.URL, err)
				w.WriteHeader(http.StatusBadGateway)
			},
		}
		o.m[domain] = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if f != nil {
				if !f(w, r, log) {
					return
				}
			}
			p.ServeHTTP(w, r)
		})
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
