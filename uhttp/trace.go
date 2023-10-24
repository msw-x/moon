package uhttp

import "github.com/msw-x/moon/ulog"

type Tracer struct {
	log         *ulog.Log
	format      Format
	formatError Format
	validate    func(Responce) bool
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

func (o *Tracer) WithValidate(f func(Responce) bool) *Tracer {
	o.validate = f
	return o
}

func (o *Tracer) Trace(r Responce) {
	ok := r.Ok()
	if o.validate != nil {
		ok = o.validate(r)
	}
	if ok {
		o.log.Debug(r.Format(o.format))
	} else {
		o.log.Error(r.Format(o.formatError))
	}
}
