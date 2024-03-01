package webs

import "time"

type Timeout struct {
	Write time.Duration
	Read  time.Duration
	Idle  time.Duration
	Close time.Duration
}
