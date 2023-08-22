package app

import (
	"time"

	"github.com/msw-x/moon/ulog"
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
		time.Sleep(time.Millisecond * 100)
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
	t := time.Now().Sub(ts)
	if t > time.Second*10 {
		return t.Truncate(time.Second)
	}
	if t > time.Second {
		return t.Truncate(time.Millisecond * 100)
	}
	if t > time.Millisecond*10 {
		return t.Truncate(time.Millisecond)
	}
	return t
}
