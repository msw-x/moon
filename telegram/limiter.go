package telegram

import (
	"time"

	"github.com/msw-x/moon/app"
)

type Limiter[T any] struct {
	job      *app.Job
	send     func(T)
	queue    QueueIf[T]
	interval time.Duration
}

func NewLimiter[T any](send func(T)) *Limiter[T] {
	o := &Limiter[T]{
		job:  app.NewJob(),
		send: send,
	}
	o.SetInterval(time.Second*3 + time.Millisecond*100)
	o.SetQueue(NewQueue[T]())
	o.job.RunTicks(o.tick, 10*time.Millisecond)
	return o
}

func (o *Limiter[T]) SetInterval(interval time.Duration) {
	o.interval = interval
}

func (o *Limiter[T]) SetQueue(queue QueueIf[T]) {
	o.queue = queue
}

func (o *Limiter[T]) Queue() QueueIf[T] {
	return o.queue
}

func (o *Limiter[T]) Close() {
	o.job.Stop()
}

func (o *Limiter[T]) Push(m T) bool {
	return o.queue.Push(m)
}

func (o *Limiter[T]) tick() {
	if text, ok := o.queue.Pop(); ok {
		o.send(text)
		o.job.Sleep(o.interval)
	}
}
