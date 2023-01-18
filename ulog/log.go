package ulog

import (
	"fmt"
)

type Log struct {
	ctx           *context
	enable        bool
	level         Level
	prefix        string
	lifetimeLevel *Level
}

func New(prefix string) *Log {
	return &Log{
		ctx:    &ctx,
		enable: true,
		prefix: prefix,
	}
}

func Empty() *Log {
	return New("").Enable(false)
}

func (o *Log) Init(opts Options) *Log {
	c := &context{}
	c.init(opts)
	o.ctx = c
	return o
}

func (o *Log) Close() {
	if o.lifetimeLevel != nil {
		o.Print(*o.lifetimeLevel, "~")
	}
	if !o.IsGloabl() {
		o.ctx.close()
	}
}

func (o *Log) IsGloabl() bool {
	return o.ctx == &ctx
}

func (o *Log) WithID(id any) *Log {
	if id != nil {
		i := fmt.Sprint(id)
		if i != "" {
			o.prefix = fmt.Sprintf("%s[%s]", o.prefix, i)
		}
	}
	return o
}

func (o *Log) WithLifetime() *Log {
	l := LevelInfo
	o.Print(l, "+")
	o.lifetimeLevel = &l
	return o
}

func (o *Log) WithLifetimeDebug() *Log {
	l := LevelDebug
	o.Print(l, "+")
	o.lifetimeLevel = &l
	return o
}

func (o *Log) WithLevel(level Level) *Log {
	o.level = level
	return o
}

func (o *Log) Enable(enable bool) *Log {
	o.enable = enable
	return o
}

func (o *Log) Branch(prefix string) *Log {
	return New(o.prefix + "." + prefix)
}

func (o *Log) Print(level Level, v ...any) {
	if o.enable && level >= o.level {
		if o.prefix != "" {
			space := ""
			if !ctx.opts.splitArgs {
				space = " "
			}
			v = append([]any{fmt.Sprintf("<%s>%s", o.prefix, space)}, v...)
		}
		print(o.ctx, level, v...)
	}
}

func (o *Log) Printf(level Level, format string, v ...any) {
	o.Print(level, fmt.Sprintf(format, v...))
}

func (o *Log) Trace(v ...any) {
	o.Print(LevelTrace, v...)
}

func (o *Log) Tracef(format string, v ...any) {
	o.Printf(LevelTrace, format, v...)
}

func (o *Log) Debug(v ...any) {
	o.Print(LevelDebug, v...)
}

func (o *Log) Debugf(format string, v ...any) {
	o.Printf(LevelDebug, format, v...)
}

func (o *Log) Info(v ...any) {
	o.Print(LevelInfo, v...)
}

func (o *Log) Infof(format string, v ...any) {
	o.Printf(LevelInfo, format, v...)
}

func (o *Log) Warning(v ...any) {
	o.Print(LevelWarning, v...)
}

func (o *Log) Warningf(format string, v ...any) {
	o.Printf(LevelWarning, format, v...)
}

func (o *Log) Error(v ...any) {
	o.Print(LevelError, v...)
}

func (o *Log) Errorf(format string, v ...any) {
	o.Printf(LevelError, format, v...)
}

func (o *Log) Critical(v ...any) {
	o.Print(LevelCritical, v...)
}

func (o *Log) Criticalf(format string, v ...any) {
	o.Printf(LevelCritical, format, v...)
}

func (o *Log) Stat() {
	o.Info(o.ctx.statistics())
}

func (o *Log) Recover() {
	if r := recover(); r != nil {
		o.Critical(r)
	}
}

func (o *Log) Query(f Filter) (lines []string, err error) {
	return o.ctx.query(f)
}
