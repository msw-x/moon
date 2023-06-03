package ujson

import (
	"strconv"
	"time"
)

type TimeMs time.Time

func (o *TimeMs) UnmarshalJSON(b []byte) error {
	s := unquote(b)
	if s == "" {
		*o = TimeMs{}
		return nil
	}
	i, err := strconv.ParseInt(s, 10, 64)
	t := time.Unix(0, i*int64(time.Millisecond))
	*o = TimeMs(t)
	return err
}

func (o TimeMs) Std() time.Time {
	return time.Time(o)
}
