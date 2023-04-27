package ujson

import (
	"time"

	"github.com/msw-x/moon/ufmt"
	"github.com/msw-x/moon/ustring"
)

type Duration time.Duration

func (o Duration) Std() time.Duration {
	return time.Duration(o)
}

func (o *Duration) UnmarshalJSON(b []byte) error {
	s := string(b)
	s = ustring.TrimQuotes(s)
	duration, err := time.ParseDuration(s)
	ufmt.Print(duration)
	*o = Duration(duration)
	return err
}
