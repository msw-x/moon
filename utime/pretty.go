package utime

import "time"

func PrettyTruncate(t time.Duration) time.Duration {
	if t > time.Second*10 {
		return t.Truncate(time.Second)
	}
	if t > time.Second {
		return t.Truncate(time.Millisecond * 100)
	}
	if t > time.Millisecond {
		return t.Truncate(time.Millisecond)
	}
	if t > time.Microsecond {
		return t.Truncate(time.Microsecond)
	}
	if t > time.Nanosecond {
		return t.Truncate(time.Nanosecond)
	}
	return t
}
