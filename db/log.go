package db

import (
	"context"
	"time"

	"github.com/msw-x/moon/ulog"
	"github.com/uptrace/bun"
)

type log struct {
	log        *ulog.Log
	onlyErrors bool
}

func newLog(ul *ulog.Log, onlyErrors bool) *log {
	return &log{
		log:        ul,
		onlyErrors: onlyErrors,
	}
}

func (o *log) Printf(format string, v ...any) {
	o.log.Infof(format, v...)
}

func (o *log) BeforeQuery(c context.Context, q *bun.QueryEvent) context.Context {
	return c
}

func (o *log) AfterQuery(c context.Context, q *bun.QueryEvent) {
	dur := time.Since(q.StartTime).Truncate(time.Microsecond)
	if q.Err == nil {
		if !o.onlyErrors {
			o.log.Debug(dur, q.Query)
		}
	} else {
		o.log.Error(dur, q.Query, q.Err)
	}
}
