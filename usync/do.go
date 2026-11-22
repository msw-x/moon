package usync

import (
	"sync"
	"sync/atomic"
	"time"
)

type Do struct {
	do         atomic.Bool
	cancel     chan struct{} // closed on first cancellation, wakes all sleepers
	cancelOnce sync.Once
	done       chan struct{} // closed by Notify
	doneOnce   sync.Once
}

func NewDo() *Do {
	o := &Do{
		cancel: make(chan struct{}),
		done:   make(chan struct{}),
	}
	o.do.Store(true)
	return o
}

func (o *Do) Do() bool {
	return o.do.Load()
}

func (o *Do) Notify() {
	o.Cancel()
	o.doneOnce.Do(func() { close(o.done) })
}

func (o *Do) Cancel() {
	// Store before close: a sleeper woken by the channel must see Do() == false.
	o.do.Store(false)
	o.closeCancel()
}

// Stop cancels and waits for Notify. Only the call that switches the flag
// from true to false waits; repeated, concurrent calls and calls after
// Cancel/Notify return immediately (Cancel followed by Stop must not block).
func (o *Do) Stop() {
	if o.do.CompareAndSwap(true, false) {
		o.closeCancel()
		<-o.done
	}
}

func (o *Do) Sleep(timeout time.Duration) {
	if timeout <= 0 || !o.Do() {
		return
	}
	timer := time.NewTimer(timeout)
	// Before Go 1.23 an unstopped timer lives until it fires.
	defer timer.Stop()
	select {
	case <-timer.C:
	case <-o.cancel:
	}
}

func (o *Do) closeCancel() {
	o.cancelOnce.Do(func() { close(o.cancel) })
}
