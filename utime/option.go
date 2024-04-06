package utime

import "time"

func Option(v time.Time) *time.Time {
	if v.IsZero() {
		return nil
	}
	return &v
}
