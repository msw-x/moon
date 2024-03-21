package utime

import "time"

func Now() time.Time {
	return time.Now().Truncate(time.Millisecond)
}
