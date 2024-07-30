package uhttp

import "github.com/msw-x/moon/ulog"

type Tracer[T TraceResponce] struct {
	log         *ulog.Log
	format      Format
	formatError Format
	validate    func(T) bool
	filter      func(T) bool
}

type TraceResponce interface {
	Ok() bool
	Format(Format) string
}

func NewTracer[T TraceResponce](log *ulog.Log) *Tracer[T] {
	o := new(Tracer[T])
	o.log = log
	return o
}

func (o *Tracer[T]) WithFormat(f Format) *Tracer[T] {
	o.format = f
	o.formatError = f
	return o
}

func (o *Tracer[T]) WithFormatError(f Format) *Tracer[T] {
	o.formatError = f
	return o
}

func (o *Tracer[T]) WithValidate(f func(T) bool) *Tracer[T] {
	o.validate = f
	return o
}

func (o *Tracer[T]) WithFilter(f func(T) bool) *Tracer[T] {
	o.filter = f
	return o
}

func (o *Tracer[T]) Trace(r T) {
	if o.filter != nil {
		if !o.filter(r) {
			return
		}
	}
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
