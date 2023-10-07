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

func GoGroupWithLog(n int, log *ulog.Log, fn func()) {
	GoGroup(n, func() {
		defer log.Recover()
		fn()
	})
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

func (o *GoSwarm) Add(fn func()) *GoSwarm {
	o.wg.Add(1)
	Go(func() {
		defer o.log.Recover()
		defer o.wg.Done()
		fn()
	})
	return o
}

func (o *GoSwarm) Wait() {
	o.wg.Wait()
}

type GoSwarmLimit struct {
	log *ulog.Log
	fns []func()
}

func NewGoSwarmLimit() *GoSwarmLimit {
	o := new(GoSwarmLimit)
	o.log = ulog.New("")
	return o
}

func (o *GoSwarmLimit) WithLog(log *ulog.Log) *GoSwarmLimit {
	o.log = log
	return o
}

func (o *GoSwarmLimit) Add(fn func()) *GoSwarmLimit {
	o.fns = append(o.fns, fn)
	return o
}

func (o *GoSwarmLimit) Execute(limit int) {
	fns := make(chan func(), len(o.fns)+1)
	for _, fn := range o.fns {
		fns <- fn
	}
	GoGroupWithLog(limit, o.log, func() {
		for {
			select {
			case fn := <-fns:
				fn()
			default:
				return
			}
		}
	})
	close(fns)
}
