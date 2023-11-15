package utime

import (
	"fmt"
	"strings"
	"time"
)

func TimeToDuration(t time.Time) time.Duration {
	return time.Hour*time.Duration(t.Hour()) +
		time.Minute*time.Duration(t.Minute()) +
		time.Second*time.Duration(t.Second())
}

func ParseDuration(s string) (r time.Duration, err error) {
	if strings.Contains(s, ":") {
		var t time.Time
		t, err = time.Parse("15:04:05", s)
		if err == nil {
			r = TimeToDuration(t)
		}
		return
	}
	return time.ParseDuration(s)
}

func FixedZone(offset time.Duration) *time.Location {
	return time.FixedZone(fmt.Sprintf("UTC%+d", int(offset.Hours())), int(offset.Seconds()))
}

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
	return t
}
