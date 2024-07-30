package urest

import (
	"github.com/msw-x/moon/uhttp"
	"github.com/msw-x/moon/ulog"
)

type Context struct {
	router      *uhttp.Router
	allowCors   bool
	format      uhttp.Format
	formatError uhttp.Format
}

func NewContext(router *uhttp.Router) *Context {
	return &Context{
		router: router,
	}
}

func (o *Context) Router() *uhttp.Router {
	return o.router
}

func (o *Context) ReverseProxy(proxy *uhttp.ReverseProxy) {
	if proxy.Tracer() == nil {
		proxy.WithTrace(uhttp.NewTracer[uhttp.ReverseProxyResponce](o.Router().Log().Branch("proxy")).
			WithFormat(o.format).
			WithFormatError(o.formatError))
	}
	o.Router().ReverseProxy(proxy)
}

func (o Context) Branch(name string) *Context {
	o.router = o.router.Branch(name)
	return &o
}

func (o *Context) AllowCors(v bool) *Context {
	o.allowCors = v
	return o
}

func (o *Context) WithFormat(f uhttp.Format) *Context {
	o.format = f
	return o.WithFormatError(f)
}

func (o *Context) WithFormatError(f uhttp.Format) *Context {
	o.formatError = f
	return o
}

func (o *Context) Log() *ulog.Log {
	return o.router.Log()
}

func (o *Context) Trace(v uhttp.FormatProvider, ok bool) {
	if ok {
		o.Log().Debug(v.Format(o.format))
	} else {
		o.Log().Error(v.Format(o.formatError))
	}
}
