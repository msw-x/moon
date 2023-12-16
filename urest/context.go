package urest

import (
	"github.com/msw-x/moon/ulog"
	"github.com/msw-x/moon/webs"
)

type Context struct {
	log       *ulog.Log
	router    *webs.Router
	allowCors bool
}

func NewContext(branch string, log *ulog.Log, router *webs.Router) Context {
	if branch != "" {
		log = log.Branch(branch)
	}
	return HandleContext{
		log:    log,
		router: router.Branch(branch),
	}
}

func (o *Context) AllowCors(v bool) Context {
	o.allowCors = v
	return *o
}

func (o Context) Branch(name string) Context {
	return NewContext(name, o.log, o.router).AllowCors(o.allowCors)
}
