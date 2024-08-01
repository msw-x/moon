package uhttp

import (
	"errors"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/msw-x/moon/ulog"
)

type ReverseProxy struct {
	log       *ulog.Log
	proxy     *httputil.ReverseProxy
	tracer    *Tracer[ReverseProxyResponse]
	target    string
	onRequest func(http.ResponseWriter, *http.Request) error
	onRewrite func(*httputil.ProxyRequest)
	onTarget  func(*httputil.ProxyRequest) string
}

func NewReverseProxy() *ReverseProxy {
	o := new(ReverseProxy)
	o.log = ulog.Empty()
	return o
}

func (o *ReverseProxy) WithLog(log *ulog.Log) *ReverseProxy {
	o.log = log
	return o
}

func (o *ReverseProxy) WithTrace(tracer *Tracer[ReverseProxyResponse]) *ReverseProxy {
	o.tracer = tracer
	return o
}

func (o *ReverseProxy) WithTarget(s string) *ReverseProxy {
	o.target = s
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

func (o *ReverseProxy) Pure() *ReverseProxy {
	o.OnRewrite(func(*httputil.ProxyRequest) {})
	return o
}

func (o *ReverseProxy) Init() {
	if o.onTarget == nil {
		if t, ok := o.targetUrl(o.target); ok {
			if o.onRewrite == nil {
				o.proxy = httputil.NewSingleHostReverseProxy(t)
			} else {
				o.proxy = &httputil.ReverseProxy{
					Rewrite: func(r *httputil.ProxyRequest) {
						o.onRewrite(r)
						r.SetURL(t)
					},
				}
			}
		} else {
			return
		}
	} else {
		o.proxy = &httputil.ReverseProxy{
			Rewrite: func(r *httputil.ProxyRequest) {
				if t, ok := o.targetUrl(o.onTarget(r)); ok {
					r.SetURL(t)
				}
			},
		}
	}
	o.proxy.ModifyResponse = func(v *http.Response) error {
		return nil
	}
	o.proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		if v, ok := w.(*ReverseProxyResponse); ok {
			v.SetError(err)
		}
		w.WriteHeader(http.StatusBadGateway)
	}
}

func (o *ReverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resp := NewReverseProxyResponse(r, w, o.tracer)
	if o.onRequest != nil {
		if resp.ErrorFree() {
			resp.SetError(o.onRequest(resp, r))
		}
	}
	if resp.ErrorFree() {
		o.proxy.ServeHTTP(resp, r)
	}
	resp.Close()
}

func (o *ReverseProxy) Log() *ulog.Log {
	return o.log
}

func (o *ReverseProxy) Tracer() *Tracer[ReverseProxyResponse] {
	return o.tracer
}

func (o *ReverseProxy) Target() string {
	return o.target
}

func (o *ReverseProxy) targetUrl(s string) (u *url.URL, ok bool) {
	var err error
	if s == "" {
		err = errors.New("url is empty")
	} else {
		u, err = url.Parse(s)
	}
	ok = err == nil
	if !ok {
		o.log.Errorf("url[%s]: %v", s, err)
	}
	return
}
