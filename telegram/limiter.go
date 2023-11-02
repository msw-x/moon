package telegram

import (
	"time"

	"github.com/msw-x/moon/app"
)

type Limiter[T any] struct {
	job      *app.Job
	send     func(T)
	messages chan T
	interval time.Duration
	prepop   func(func() (T, bool)) (T, bool)
}

func NewLimiter[T any](send func(T)) *Limiter[T] {
	o := &Limiter[T]{
		job:  app.NewJob(),
		send: send,
	}
	o.WithInterval(time.Second*3 + time.Millisecond*100)
	o.WithCapacity(4000)
	o.WithPrePop(func(pop func() (T, bool)) (T, bool) {
		return pop()
	})
	o.job.RunTicks(o.tick, 10*time.Millisecond)
	return o
}

func (o *Limiter[T]) WithInterval(interval time.Duration) *Limiter[T] {
	o.interval = interval
	return o
}

func (o *Limiter[T]) WithCapacity(capacity int) *Limiter[T] {
	o.messages = make(chan T, capacity)
	return o
}

func (o *Limiter[T]) WithPrePop(prepop func(func() (T, bool)) (T, bool)) *Limiter[T] {
	o.prepop = prepop
	return o
}

func (o *Limiter[T]) Capacity() int {
	return cap(o.messages)
}

func (o *Limiter[T]) Size() int {
	return len(o.messages)
}

func (o *Limiter[T]) Reset() {
	o.WithCapacity(cap(o.messages))
}

func (o *Limiter[T]) Close() {
	o.job.Stop()
}

func (o *Limiter[T]) Push(text T) bool {
	select {
	case o.messages <- text:
	default:
		return false
	}
	return true
}

func (o *Limiter[T]) pop() (text T, ok bool) {
	select {
	case text = <-o.messages:
		ok = true
	default:
	}
	return
}

func (o *Limiter[T]) tick() {
	if text, ok := o.prepop(o.pop); ok {
		o.send(text)
		o.job.Sleep(o.interval)
	}
}
