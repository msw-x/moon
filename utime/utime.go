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

func SetLocation(v time.Time, loc *time.Location) time.Time {
	_, src := v.Zone()
	_, dst := time.Now().In(loc).Zone()
	return v.In(loc).Add(time.Second * time.Duration(src-dst))
}
