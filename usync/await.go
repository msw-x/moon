package usync

import (
	"time"
)

type Await chan bool

func NewAwait() Await {
	return make(Await, 1)
}

func (o *Await) Notify() {
	select {
	case *o <- true:
	default:
	}
}

func (o *Await) Cancel() {
	select {
	case *o <- false:
	default:
	}
}

func (o *Await) Wait() bool {
	return <-*o
}

func (o *Await) WaitTimeout(timeout time.Duration) bool {
	select {
	case cause := <-*o:
		return cause
	case <-time.After(timeout):
		return false
	}
}
