package app

import (
	"sync"

	"github.com/msw-x/moon/ulog"
)

func Go(fn func()) {
	go func() {
		defer ulog.Recover()
		fn()
	}()
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
