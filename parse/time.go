package parse

import (
	"time"

	"github.com/msw-x/moon"
)

func Time(format string, s string) time.Time {
	t, err := time.Parse(format, s)
	moon.Check(err, "time parse")
	return t
}
