package ulog

import (
	"fmt"
	"strings"
	"time"
)

type Message struct {
	Time  time.Time
	GoID  int
	Level Level
	Text  string

	ctx     *context
	message string
}

func NewMessage(ctx *context, level Level, v ...any) Message {
	text := ""
	if ctx.opts.splitArgs {
		s := []string{}
		for _, a := range v {
			s = append(s, fmt.Sprint(a))
		}
		text = strings.Join(s, " ")
	} else {
		text = fmt.Sprint(v...)
	}
	return Message{
		Time:  time.Now(),
		GoID:  ctx.goroutineID(),
		Level: level,
		Text:  text,
		ctx:   ctx,
	}
}

func (o *Message) Format() string {
	if o.message == "" {
		l := o.Level.Laconic()
		if o.ctx.opts.GoID {
			l = fmt.Sprintf("%s|%s", o.ctx.fmtGoroutineID(o.GoID), l)
		}
		o.message = fmt.Sprintf(
			"%s [%s] %s\n",
			o.ctx.fmtTime(o.Time),
			l,
			o.Text,
		)
	}
	return o.message
}

func (o *Message) Size() int {
	return len(o.Format())
}
