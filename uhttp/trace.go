package uhttp

import "github.com/msw-x/moon/ulog"

type Tracer struct {
	log    *ulog.Log
	format Format
}

func NewTracer(log *ulog.Log) *Tracer {
	o := new(Tracer)
	o.log = log
	return o
}

func (o *Tracer) WithFormat(f Format) *Tracer {
	o.format = f
	return o
}

func (o *Tracer) Trace(r Responce) {
	m := r.Format(o.format)
	if r.Ok() {
		o.log.Debug(m)
	} else {
		o.log.Error(m)
	}
}

func TraceFormat(log *ulog.Log, format Format) func(Responce) {
	t := NewTracer(log).WithFormat(format)
	return func(r Responce) {
		t.Trace(r)
	}
}

func Trace(log *ulog.Log) func(Responce) {
	return TraceFormat(log, Format{})
}
