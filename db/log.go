package db

import (
	"context"
	"time"

	"github.com/msw-x/moon/ulog"
	"github.com/msw-x/moon/utime"
	"github.com/uptrace/bun"
)

type log struct {
	log             *ulog.Log
	onlyErrors      bool
	longQueriesTime time.Duration
	warnLongQueries bool
}

func newLog(ul *ulog.Log, onlyErrors bool) *log {
	return &log{
		log:        ul,
		onlyErrors: onlyErrors,
	}
}

func (o *log) WithQueriesTime(t time.Duration, warn bool) *log {
	o.longQueriesTime = t
	o.warnLongQueries = warn
	return o
}

func (o *log) Printf(format string, v ...any) {
	o.log.Infof(format, v...)
}

func (o *log) BeforeQuery(c context.Context, q *bun.QueryEvent) context.Context {
	return c
}

func (o *log) AfterQuery(c context.Context, q *bun.QueryEvent) {
	if q.Err == nil {
		if o.longQueriesTime > 0 && time.Since(q.StartTime) > o.longQueriesTime {
			if o.warnLongQueries {
				o.log.Warning(queryTime(q), q.Query)
			} else {
				o.log.Debug(queryTime(q), q.Query)
			}
			return
		}
		if !o.onlyErrors {
			o.log.Debug(queryTime(q), q.Query)
		}
	} else {
		o.log.Error(queryTime(q), q.Query, q.Err)
	}
}

func queryTime(q *bun.QueryEvent) time.Duration {
	return utime.PrettyTruncate(time.Since(q.StartTime))
}
