package syn

import (
	"time"
)

type Chan chan bool

func NewChan() Chan {
	return make(Chan, 1)
}

func (this *Chan) Notify() {
	select {
	case *this <- true:
	default:
	}
}

func (this *Chan) Cancel() {
	select {
	case *this <- false:
	default:
	}
}

func (this *Chan) Wait() bool {
	return <-*this
}

func (this *Chan) WaitTimeout(timeout time.Duration) bool {
	select {
	case cause := <-*this:
		return cause
	case <-time.After(timeout):
		return false
	}
	return false
}
