package uhttp

import (
	"context"
	"errors"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

type ReverseProxy struct {
	path      string
	onRequest func(http.ResponseWriter, *http.Request) bool
	onRewrite func(*httputil.ProxyRequest)
	onTarget  func(*httputil.ProxyRequest) string
	target    string
	tracer    *Tracer[ReverseProxyResponce]
}

func NewReverseProxy(s string) *ReverseProxy {
	o := new(ReverseProxy)
	o.path = s
	return o
}

func (o *ReverseProxy) WithTrace(tracer *Tracer[ReverseProxyResponce]) *ReverseProxy {
	o.tracer = tracer
	return o
}

func (o *ReverseProxy) OnRequest(f func(http.ResponseWriter, *http.Request) bool) *ReverseProxy {
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
			log.Debugf("PROXY:%s->%s", path, t)
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
		log.Debugf("PROXY:%s", path)
		proxy = &httputil.ReverseProxy{
			Rewrite: func(r *httputil.ProxyRequest) {
				if t, ok := target(o.onTarget(r)); ok {
					r.SetURL(t)
				}
			},
		}
	}
	toContextRequest := func(r *http.Request, key string, val any) {
		if o.tracer != nil {
			r.WithContext(context.WithValue(r.Context(), key, val))
		}
	}
	proxy.ModifyResponse = func(v *http.Response) error {
		toContextRequest(v.Request, "status", v.Status)
		toContextRequest(v.Request, "status-code", v.StatusCode)
		return nil
	}
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		toContextRequest(r, "error", err)
		w.WriteHeader(http.StatusBadGateway)
	}
	router.HandleFunc(path+"{path:.*}", func(w http.ResponseWriter, r *http.Request) {
		ts := time.Now()
		ok := true
		if o.onRequest != nil {
			ok = o.onRequest(w, r)
		}
		if ok {
			proxy.ServeHTTP(w, r)
		}
		if o.tracer != nil {
			o.trace(router, w, r, time.Since(ts))
		}
	})
}

func (o *ReverseProxy) Tracer() *Tracer[ReverseProxyResponce] {
	return o.tracer
}

func (o *ReverseProxy) trace(router *Router, w http.ResponseWriter, r *http.Request, tm time.Duration) {
	var trace ReverseProxyResponce
	trace.Request = r
	trace.Header = w.Header()
	trace.Time = tm
	trace.router = router
	if v := r.Context().Value("status"); v != nil {
		trace.Status = v.(string)
	}
	if v := r.Context().Value("status-code"); v != nil {
		trace.StatusCode = v.(int)
	}
	if v := r.Context().Value("error"); v != nil {
		trace.Error = v.(error)
	}
	o.tracer.Trace(trace)
}
