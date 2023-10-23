package uhttp

import "github.com/msw-x/moon/ulog"

type Tracer struct {
	log         *ulog.Log
	format      Format
	formatError Format
}

func NewTracer(log *ulog.Log) *Tracer {
	o := new(Tracer)
	o.log = log
	return o
}

func (o *Tracer) WithFormat(f Format) *Tracer {
	o.format = f
	o.formatError = f
	return o
}

func (o *Tracer) WithFormatError(f Format) *Tracer {
	o.formatError = f
	return o
}

func (o *Tracer) Trace(r Responce) {
	if r.Ok() {
		o.log.Debug(r.Format(o.format))
	} else {
		o.log.Error(r.Format(o.formatError))
	}
}

func TraceFormat(log *ulog.Log, format Format) func(Responce) {
	t := NewTracer(log).WithFormat(format)
	return func(r Responce) {
		t.Trace(r)
	}
}

func TraceTwinFormat(log *ulog.Log, format Format, formatError Format) func(Responce) {
	t := NewTracer(log).WithFormat(format).WithFormatError(formatError)
	return func(r Responce) {
		t.Trace(r)
	}
}

func Trace(log *ulog.Log) func(Responce) {
	return TraceFormat(log, Format{})
}
