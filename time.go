package moon

import (
	"fmt"
	"time"
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

func (this *Time) SetDailyHours(h int) {
	if h >= 24 {
		Panic("hours must be less than 24 but not", h)
	}
	this.Hours = h
}

func (this *Time) SetMinutes(m int) {
	if m >= 60 {
		Panic("minutes must be less than 60 but not", m)
	}
	this.Minutes = m
}

func (this *Time) SetSeconds(s int) {
	if s >= 60 {
		Panic("seconds must be less than 60 but not", s)
	}
	this.Seconds = s
}

func (this *Time) SetTime(t time.Time) *Time {
	this.Hours = t.Hour()
	this.Minutes = t.Minute()
	this.Seconds = t.Second()
	return this
}

func (this *Time) SetNow() *Time {
	return this.SetTime(time.Now())
}

func (this Time) Days() (days, hours int) {
	days = this.Hours / 24
	hours = this.Hours % 24
	return
}

func (this Time) TotalMinutes() int {
	return this.Hours*60 + this.Minutes
}

func (this Time) TotalSeconds() int {
	return this.TotalMinutes()*60 + this.Seconds
}

func (this Time) Format() string {
	return fmt.Sprintf("%02d:%02d:%02d", this.Hours, this.Minutes, this.Seconds)
}

func (this Time) FormatMs() string {
	return fmt.Sprintf("%s.%03d", this.Format(), this.Milliseconds)
}

func (this Time) FormatDays() string {
	days, hours := this.Days()
	if days > 0 {
		plural := ""
		if days > 1 {
			plural = "s"
		}
		return fmt.Sprintf("%d day%s %02d:%02d:%02d", days, plural, hours, this.Minutes, this.Seconds)
	}
	return this.Format()
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
