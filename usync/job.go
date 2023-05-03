package usync

import (
	"time"
)

type Job struct {
	do bool
	ch Await
}

func NewJob() *Job {
	return &Job{
		do: true,
		ch: NewAwait(),
	}
}

func (o *Job) Do() bool {
	return o.do
}

func (o *Job) Notify() {
	o.Cancel()
	o.ch.Notify()
}

func (o *Job) Cancel() {
	o.do = false
}

func (o *Job) Stop() {
	if o.do {
		o.Cancel()
		o.ch.Wait()
	}
}

func (o *Job) Sleep(timeout time.Duration) {
	if o.Do() {
		const maxSleepTime = time.Millisecond * 10
		if timeout > maxSleepTime {
			count := int(timeout / maxSleepTime)
			for n := 0; n != count; n++ {
				if !o.Do() {
					return
				}
				time.Sleep(maxSleepTime)
			}
		} else {
			time.Sleep(timeout)
		}
	}
}
