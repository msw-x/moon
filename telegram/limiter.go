package telegram

import (
	"time"

	"github.com/msw-x/moon/app"
)

type Limiter struct {
	job      *app.Job
	send     func(string)
	messages chan string
	interval time.Duration
	prepop   func(func() (string, bool)) (string, bool)
}

func NewLimiter(send func(string)) *Limiter {
	o := &Limiter{
		job:  app.NewJob(),
		send: send,
	}
	o.WithInterval(time.Second*3 + time.Millisecond*100)
	o.WithCapacity(4000)
	o.WithPrePop(func(pop func() (string, bool)) (string, bool) {
		return pop()
	})
	o.job.RunTicks(o.tick, 10*time.Millisecond)
	return o
}

func (o *Limiter) WithInterval(interval time.Duration) *Limiter {
	o.interval = interval
	return o
}

func (o *Limiter) WithCapacity(capacity int) *Limiter {
	o.messages = make(chan string, capacity)
	return o
}

func (o *Limiter) WithPrePop(prepop func(func() (string, bool)) (string, bool)) *Limiter {
	o.prepop = prepop
	return o
}

func (o *Limiter) Close() {
	o.job.Stop()
}

func (o *Limiter) Push(text string) bool {
	select {
	case o.messages <- text:
	default:
		return false
	}
	return true
}

func (o *Limiter) pop() (text string, ok bool) {
	select {
	case text = <-o.messages:
		ok = true
	default:
	}
	return
}

func (o *Limiter) tick() {
	if text, ok := o.prepop(o.pop); ok {
		o.send(text)
		o.job.Sleep(o.interval)
	}
}
