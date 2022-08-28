package syn

import (
	"time"
)

type Do struct {
	do bool
	ch Chan
}

func NewDo() *Do {
	return &Do{
		do: true,
		ch: NewChan(),
	}
}

func (this *Do) Do() bool {
	return this.do
}

func (this *Do) Notify() {
	this.ch.Notify()
}

func (this *Do) Cancel() {
	this.do = false
}

func (this *Do) Stop() {
	if this.do {
		this.Cancel()
		this.ch.Wait()
	}
}

func (this *Do) Sleep(timeout time.Duration) {
	if this.Do() {
		const maxSleepTime = time.Millisecond * 10
		if timeout > maxSleepTime {
			count := int(timeout / maxSleepTime)
			for n := 0; n != count; n++ {
				if !this.Do() {
					return
				}
				time.Sleep(maxSleepTime)
			}
		} else {
			time.Sleep(timeout)
		}
	}
}
