package usync

import (
	"sync"
	"sync/atomic"
	"time"
)

// timerPool holds stopped or expired timers with drained channels.
var timerPool sync.Pool

// Do is created with NewDo. The zero value is an already canceled Do:
// its methods return immediately.
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
	o.closeDone()
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
	// Timers are pooled: a timer and its channel per call would cost
	// three allocations on every Sleep.
	timer, _ := timerPool.Get().(*time.Timer)
	if timer == nil {
		timer = time.NewTimer(timeout)
	} else {
		timer.Reset(timeout)
	}
	select {
	case <-timer.C:
		// Expired and drained: safe to reuse.
	case <-o.cancel:
		if !timer.Stop() {
			// Fired concurrently: before Go 1.23 its channel holds (or is
			// about to get) a value, so the timer must not be reused.
			return
		}
	}
	timerPool.Put(timer)
}

// The channels are nil in the zero value of Do.
func (o *Do) closeCancel() {
	o.cancelOnce.Do(func() {
		if o.cancel != nil {
			close(o.cancel)
		}
	})
}

func (o *Do) closeDone() {
	o.doneOnce.Do(func() {
		if o.done != nil {
			close(o.done)
		}
	})
}
