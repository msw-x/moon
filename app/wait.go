package app

import (
	"time"

	"github.com/msw-x/moon/ulog"
	"github.com/msw-x/moon/utime"
)

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
