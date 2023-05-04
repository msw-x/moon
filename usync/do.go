package usync

import (
	"time"
)

type Do struct {
	do bool
	ch Await
}

func NewDo() *Do {
	return &Do{
		do: true,
		ch: NewAwait(),
	}
}

func (o *Do) Do() bool {
	return o.do
}

func (o *Do) Notify() {
	o.Cancel()
	o.ch.Notify()
}

func (o *Do) Cancel() {
	o.do = false
}

func (o *Do) Stop() {
	if o.do {
		o.Cancel()
		o.ch.Wait()
	}
}

func (o *Do) Sleep(timeout time.Duration) {
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
