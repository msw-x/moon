package app

import (
	"sync"
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

func GoGroup(n int, fn func()) {
	var wg sync.WaitGroup
	for i := 0; i != n; i++ {
		wg.Add(1)
		Go(func() {
			defer wg.Done()
			fn()
		})
	}
	wg.Wait()
}
