package uhttp

import (
	"errors"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gorilla/mux"
)

type ReverseProxy struct {
	path      string
	onRequest func(http.ResponseWriter, *http.Request) bool
	onRewrite func(*httputil.ProxyRequest)
	onTarget  func(*httputil.ProxyRequest) string
	target    string
}

func NewReverseProxy(s string) *ReverseProxy {
	o := new(ReverseProxy)
	o.path = s
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
	path := router.nextPath(o.path)
	target := func(s string) (u *url.URL, ok bool) {
		var err error
		if s == "" {
			err = errors.New("target is empty")
		}
		u, err = url.Parse(s)
		ok = err == nil
		if !ok {
			router.log.Errorf("proxy url[%s]: %v", s, err)
		}
		return
	}
	var proxy *httputil.ReverseProxy
	if o.onTarget == nil {
		if t, ok := target(o.target); ok {
			router.log.Debugf("PROXY:%s->%s", path, t)
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
		router.log.Debugf("PROXY:%s", path)
		proxy = &httputil.ReverseProxy{
			Rewrite: func(r *httputil.ProxyRequest) {
				if t, ok := target(o.onTarget(r)); ok {
					r.SetURL(t)
				}
			},
		}
	}
	router.HandleFunc(path+"{path:.*}", func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = mux.Vars(r)["path"]
		ok := true
		if o.onRequest != nil {
			ok = o.onRequest(w, r)
		}
		if ok {
			proxy.ServeHTTP(w, r)
		}
	})
}
