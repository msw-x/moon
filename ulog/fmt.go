package ulog

import (
	"fmt"
	"moon/rt"
	"strconv"
	"strings"
	"time"
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
		this.message = fmt.Sprintf(
			"%s [%s|%s] %s\n",
			this.Time,
			this.GoID,
			this.Level,
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
