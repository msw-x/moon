package uhttp

import (
	"errors"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type ReverseProxy struct {
	path      string
	onRequest func(http.ResponseWriter, *http.Request) error
	onRewrite func(*httputil.ProxyRequest)
	onTarget  func(*httputil.ProxyRequest) string
	target    string
	tracer    *Tracer[ReverseProxyResponse]
}

func NewReverseProxy(s string) *ReverseProxy {
	o := new(ReverseProxy)
	o.path = s
	return o
}

func (o *ReverseProxy) WithTrace(tracer *Tracer[ReverseProxyResponse]) *ReverseProxy {
	o.tracer = tracer
	return o
}

func (o *ReverseProxy) OnRequest(f func(http.ResponseWriter, *http.Request) error) *ReverseProxy {
	o.onRequest = f
	return o
}

func (o *ReverseProxy) OnRewrite(f func(*httputil.ProxyRequest)) *ReverseProxy {
	o.onRewrite = f
	return o
}

func (o *ReverseProxy) OnTarget(f func(*httputil.ProxyRequest) string) *ReverseProxy {
	o.onTarget = f
	return o
}

func (o *ReverseProxy) Target(s string) *ReverseProxy {
	o.target = s
	return o
}

func (o *ReverseProxy) Pure() *ReverseProxy {
	o.OnRewrite(func(*httputil.ProxyRequest) {})
	return o
}

func (o *ReverseProxy) Connect(router *Router) {
	log := router.log
	path := router.nextPath(o.path)
	target := func(s string) (u *url.URL, ok bool) {
		var err error
		if s == "" {
			err = errors.New("target is empty")
		}
		u, err = url.Parse(s)
		ok = err == nil
		if !ok {
			log.Errorf("proxy url[%s]: %v", s, err)
		}
		return
	}
	var proxy *httputil.ReverseProxy
	if o.onTarget == nil {
		if t, ok := target(o.target); ok {
			log.Debugf("proxy:%s->%s", path, t)
			if o.onRewrite == nil {
				proxy = httputil.NewSingleHostReverseProxy(t)
			} else {
				proxy = &httputil.ReverseProxy{
					Rewrite: func(r *httputil.ProxyRequest) {
						o.onRewrite(r)
						r.SetURL(t)
					},
				}
			}
		}
	} else {
		log.Debugf("proxy:%s", path)
		proxy = &httputil.ReverseProxy{
			Rewrite: func(r *httputil.ProxyRequest) {
				if t, ok := target(o.onTarget(r)); ok {
					r.SetURL(t)
				}
			},
		}
	}
	proxy.ModifyResponse = func(v *http.Response) error {
		return nil
	}
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Trace("ErrorHandler 1")
		if v, ok := w.(*ReverseProxyResponse); ok {
			v.SetError(err)
			log.Trace("ErrorHandler 2")
		}
		log.Trace("ErrorHandler 3")
		w.WriteHeader(http.StatusBadGateway)
	}
	router.HandleFunc(path+"{path:.*}", func(w http.ResponseWriter, r *http.Request) {
		resp := NewReverseProxyResponse(r, w, o.tracer, router)
		if o.onRequest != nil {
			if resp.ErrorFree() {
				resp.SetError(o.onRequest(resp, r))
			}
		}
		if resp.ErrorFree() {
			proxy.ServeHTTP(resp, r)
		}
		resp.Close()
	})
}

func (o *ReverseProxy) Tracer() *Tracer[ReverseProxyResponse] {
	return o.tracer
}
