package app

import (
	"time"

	"github.com/msw-x/moon/ulog"
	"github.com/msw-x/moon/usync"
)

func Go(fn func()) GoDo {
	do := usync.NewDo()
	go func() {
		defer do.Notify()
		defer ulog.Recover()
		fn()
	}()
	return do
}

type GoDo interface {
	Do() bool
	Stop()
	Cancel()
	Sleep(timeout time.Duration)
}
