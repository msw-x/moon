package moon

import (
	"fmt"
	"time"

	"github.com/msw-x/moon/uerr"
)

type Time struct {
	Hours        int
	Minutes      int
	Seconds      int
	Milliseconds int
}

func NewTime() *Time {
	return &Time{}
}

func Now() *Time {
	return NewTime().SetNow()
}

func (o *Time) SetDailyHours(h int) {
	if h >= 24 {
		uerr.Panic("hours must be less than 24 but not", h)
	}
	o.Hours = h
}

func (o *Time) SetMinutes(m int) {
	if m >= 60 {
		uerr.Panic("minutes must be less than 60 but not", m)
	}
	o.Minutes = m
}

func (o *Time) SetSeconds(s int) {
	if s >= 60 {
		uerr.Panic("seconds must be less than 60 but not", s)
	}
	o.Seconds = s
}

func (o *Time) SetTime(t time.Time) *Time {
	o.Hours = t.Hour()
	o.Minutes = t.Minute()
	o.Seconds = t.Second()
	return o
}

func (o *Time) SetNow() *Time {
	return o.SetTime(time.Now())
}

func (o Time) Days() (days, hours int) {
	days = o.Hours / 24
	hours = o.Hours % 24
	return
}

func (o Time) TotalMinutes() int {
	return o.Hours*60 + o.Minutes
}

func (o Time) TotalSeconds() int {
	return o.TotalMinutes()*60 + o.Seconds
}

func (o Time) Format() string {
	return fmt.Sprintf("%02d:%02d:%02d", o.Hours, o.Minutes, o.Seconds)
}

func (o Time) FormatMs() string {
	return fmt.Sprintf("%s.%03d", o.Format(), o.Milliseconds)
}

func (o Time) FormatDays() string {
	days, hours := o.Days()
	if days > 0 {
		plural := ""
		if days > 1 {
			plural = "s"
		}
		return fmt.Sprintf("%d day%s %02d:%02d:%02d", days, plural, hours, o.Minutes, o.Seconds)
	}
	return o.Format()
}

func DurationToTime(d time.Duration) Time {
	d = d.Round(time.Millisecond)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	d -= s * time.Second
	ms := d / time.Millisecond
	return Time{int(h), int(m), int(s), int(ms)}
}

func FormatTime(t time.Time) string {
	return NewTime().SetTime(t).Format()
}

func FormatDuration(d time.Duration) string {
	return DurationToTime(d).Format()
}

func FormatDurationMs(d time.Duration) string {
	return DurationToTime(d).FormatMs()
}

func FormatDurationDays(d time.Duration) string {
	return DurationToTime(d).FormatDays()
}
