package app

import (
	"time"

	"github.com/msw-x/moon/ulog"
	"github.com/msw-x/moon/usync"
)

type Do struct {
	do *usync.Do
}

func NewDo() *Do {
	return &Do{
		do: usync.NewDo(),
	}
}

func (o *Do) Do() bool {
	return o.do.Do()
}

func (o *Do) Wait() {
	for o.Do() {
		o.Sleep(time.Millisecond)
	}
}

func (o *Do) Stop() {
	o.do.Stop()
}

func (o *Do) Cancel() {
	o.do.Cancel()
}

func (o *Do) Sleep(timeout time.Duration) {
	o.do.Sleep(timeout)
}

func (o *Do) Run(fn func()) {
	go func() {
		defer o.do.Cancel()
		defer ulog.Recover()
		fn()
	}()
}

func (o *Do) RunLoop(fn func()) {
	o.Run(func() {
		for o.do.Do() {
			fn()
		}
	})
}

func (o *Do) Ticks(fn func(), interval time.Duration) {
	o.RunLoop(func() {
		fn()
		o.do.Sleep(interval)
	})
}
