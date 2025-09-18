package ulog

import "time"

type Queue[T any] struct {
	m chan T
}

func NewQueue[T any]() *Queue[T] {
	o := new(Queue[T])
	o.SetCapacity(10000)
	return o
}

func (o *Queue[T]) SetCapacity(capacity int) {
	o.m = make(chan T, capacity)
}

func (o *Queue[T]) Capacity() int {
	return cap(o.m)
}

func (o *Queue[T]) Size() int {
	return len(o.m)
}

func (o *Queue[T]) Reset() {
	o.SetCapacity(cap(o.m))
}

func (o *Queue[T]) Push(v T) bool {
	select {
	case o.m <- v:
	default:
		return false
	}
	return true
}

func (o *Queue[T]) Pop() (v T, ok bool) {
	select {
	case v = <-o.m:
		ok = true
	case <-time.After(100 * time.Millisecond):
		// timeout
	}
	return
}
