package ujson

import (
	"time"

	"github.com/msw-x/moon/ustring"
)

type Duration time.Duration

func (o Duration) Std() time.Duration {
	return time.Duration(o)
}

func (o Duration) String() string {
	return time.Duration(o).String()
}

func (o Duration) MarshalJSON() ([]byte, error) {
	return quote(o.String()), nil
}

func (o *Duration) UnmarshalJSON(b []byte) error {
	s := string(b)
	s = ustring.TrimQuotes(s)
	duration, err := time.ParseDuration(s)
	*o = Duration(duration)
	return err
}
