package uhttp

import (
	"fmt"
	"time"

	"github.com/msw-x/moon/utime"
)

type Timeouts struct {
	Write time.Duration
	Read  time.Duration
	Idle  time.Duration
	Close time.Duration
}

func (o *Timeouts) Set(t Timeouts) {
	if t.Write != 0 {
		o.Write = t.Write
	}
	if t.Read != 0 {
		o.Read = t.Read
	}
	if t.Idle != 0 {
		o.Idle = t.Idle
	}
	if t.Close != 0 {
		o.Close = t.Close
	}
}

func (o Timeouts) String() string {
	return fmt.Sprintf("write[%s] read[%s] idle[%s] close[%s]",
		utime.Pretty(o.Write), utime.Pretty(o.Read), utime.Pretty(o.Idle), utime.Pretty(o.Close))
}
