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

type GoSwarm struct {
	log *ulog.Log
	wg  sync.WaitGroup
}

func NewGoSwarm() *GoSwarm {
	o := new(GoSwarm)
	o.log = ulog.New("")
	return o
}

func (o *GoSwarm) WithLog(log *ulog.Log) *GoSwarm {
	o.log = log
	return o
}

func (o *GoSwarm) Add(fn func()) {
	o.wg.Add(1)
	Go(func() {
		defer o.log.Recover()
		defer o.wg.Done()
		fn()
	})
}

func (o *GoSwarm) Wait() {
	o.wg.Wait()
}
