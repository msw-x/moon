package app

import "time"

type Tick struct {
	fn       func()
	enabled  bool
	interval time.Duration
	last     time.Time
}

func NewTick(fn func(), interval time.Duration) *Tick {
	o := new(Tick)
	o.fn = fn
	o.enabled = true
	o.interval = interval
	return o
}

func (o *Tick) WithEnable(v bool) *Tick {
	o.enabled = v
	return o
}

func (o *Tick) Do() {
	if o.enabled && time.Since(o.last) > o.interval {
		o.last = time.Now()
		o.fn()
	}
}

func (o *Tick) Reset() {
	o.last = time.Time{}
}

func (o *Tick) LastTime() time.Time {
	return o.last
}

func (o *Tick) LastDuration() time.Duration {
	return time.Since(o.last)
}
