package ulog

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/msw-x/moon/rt"
)

func Format(level Level, v ...any) string {
	return NewMessage(level, v...).Format()
}

type Message struct {
	Time    string
	GoID    string
	Level   string
	Text    string
	message string
}

func NewMessage(level Level, v ...any) *Message {
	text := ""
	if ctx.conf.splitArgs {
		s := []string{}
		for _, a := range v {
			s = append(s, fmt.Sprint(a))
		}
		text = strings.Join(s, " ")
	} else {
		text = fmt.Sprint(v...)
	}
	return &Message{
		Time:  fmtTime(),
		GoID:  fmtGoroutineID(),
		Level: level.Laconic(),
		Text:  text,
	}
}

func (this *Message) Format() string {
	if this.message == "" {
		l := this.Level
		if ctx.conf.GoID {
			l = fmt.Sprintf("%s|%s", this.GoID, this.Level)
		}
		this.message = fmt.Sprintf(
			"%s [%s] %s\n",
			this.Time,
			l,
			this.Text,
		)
	}
	return this.message
}

func fmtTime() string {
	ts := time.Now()
	ms := ts.Sub(ts.Truncate(time.Second)).Milliseconds()
	return fmt.Sprintf("%s.%03d", ts.Format("2006-Jan-02 15:04:05"), ms)
}

func fmtGoroutineID() string {
	if !ctx.conf.GoID {
		return ""
	}
	id := rt.GoroutineID()
	sid := strconv.Itoa(id)
	if len(sid) > ctx.maxid {
		ctx.maxid = len(sid)
	}
	for n := ctx.maxid - len(sid); n != 0; n-- {
		sid = " " + sid
	}
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()
	ctx.mapid[id] = true
	return sid
}
