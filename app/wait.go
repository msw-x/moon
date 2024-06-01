package app

import (
	"time"

	"github.com/msw-x/moon/ulog"
	"github.com/msw-x/moon/utime"
)

/*
type Wait struct {
	log     *ulog.Log
	do      func() bool
	ts      time.Time
	timeout time.Duration
	waited  bool
}

func NewWait() *Wait {
	o := new(Wait)
	o.log = ulog.Empty()
	o.do = func() bool {
		return true
	}
	return o
}

func (o *Wait) WithLog(log *ulog.Log) *Wait {
	o.log = log
	return o
}

func (o *Wait) WithDo(do func() bool) *Wait {
	o.do = do
	return o
}

func (o *Wait) WithTimeout(timeout time.Duration) *Wait {
	o.timeout = timeout
	return o
}

func (o *Wait) Wait(ok func() bool) *Wait {
	o.log.Info("wait...")
	o.ts = time.Now()
	before := o.ts.Add(o.timeout)
	var exceed bool
	for o.do() {
		if o.timeout != 0 && time.Now().After(before) {
			if !exceed {
				o.log.Warning("wait timeout:", o.timeout)
				exceed = true
			}
		}
		if ok() {
			o.waited = true
			break
		}
		time.Sleep(time.Millisecond)
	}
	waited := o.Time()
	if exceed {
		o.log.Warning("waited:", waited)
	} else {
		o.log.Info("waited:", waited)
	}
	return o
}

func (o *Wait) Time() time.Duration {
	if o.ts.IsZero() {
		return 0
	}
	return utime.PrettyTruncate(time.Now().Sub(o.ts))
}

func (o *Wait) Waited() bool {
	return o.waited
}
*/

func Wait(log *ulog.Log, fn func() bool, timeout time.Duration) time.Duration {
	log.Info("wait...")
	ts := time.Now()
	before := ts.Add(timeout)
	var exceed bool
	for {
		if time.Now().After(before) {
			if !exceed {
				log.Warning("wait timeout:", timeout)
				exceed = true
			}
		}
		if fn() {
			break
		}
		time.Sleep(time.Millisecond)
	}
	waited := Waited(ts)
	if exceed {
		log.Warning("waited:", waited)
	} else {
		log.Info("waited:", waited)
	}
	return waited
}

func Waited(ts time.Time) time.Duration {
	return utime.PrettyTruncate(time.Now().Sub(ts))
}
