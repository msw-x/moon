package parse

import (
	"time"

	"github.com/msw-x/moon"
)

func Time(format string, s string) (time.Time, error) {
	return time.Parse(format, s)
}

func TimeStrict(format string, s string) time.Time {
	t, err := Time(format, s)
	moon.Strict(err, "time parse")
	return t
}
