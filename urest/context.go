package urest

import (
	"github.com/msw-x/moon/uhttp"
	"github.com/msw-x/moon/ulog"
)

type Context struct {
	router      *uhttp.Router
	allowCors   bool
	trace       bool
	traceError  bool
	format      uhttp.Format
	formatError uhttp.Format
}

func NewContext(router *uhttp.Router) *Context {
	return &Context{
		router: router,
	}
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
	o.trace = true
	o.format = f
	return o.WithFormatError(f)
}

func (o *Context) WithFormatError(f uhttp.Format) *Context {
	o.traceError = true
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
