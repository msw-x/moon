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

func (this *Log) Init(opts Options) *Log {
	c := &context{}
	c.init(opts)
	this.ctx = c
	return this
}

func (this *Log) Close() {
	if this.lifetimeLevel != nil {
		this.Print(*this.lifetimeLevel, "~")
	}
	if !this.IsGloabl() {
		this.ctx.close()
	}
}

func (this *Log) IsGloabl() bool {
	return this.ctx == &ctx
}

func (this *Log) WithID(id any) *Log {
	if id != nil {
		i := fmt.Sprint(id)
		if i != "" {
			this.prefix = fmt.Sprintf("%s[%s]", this.prefix, i)
		}
	}
	return this
}

func (this *Log) WithLifetime() *Log {
	l := LevelInfo
	this.Print(l, "+")
	this.lifetimeLevel = &l
	return this
}

func (this *Log) WithLifetimeDebug() *Log {
	l := LevelDebug
	this.Print(l, "+")
	this.lifetimeLevel = &l
	return this
}

func (this *Log) WithLevel(level Level) *Log {
	this.level = level
	return this
}

func (this *Log) Enable(enable bool) *Log {
	this.enable = enable
	return this
}

func (this *Log) Branch(prefix string) *Log {
	return New(this.prefix + "." + prefix)
}

func (this *Log) Print(level Level, v ...any) {
	if this.enable && level >= this.level {
		if this.prefix != "" {
			space := ""
			if !ctx.opts.splitArgs {
				space = " "
			}
			v = append([]any{fmt.Sprintf("<%s>%s", this.prefix, space)}, v...)
		}
		print(this.ctx, level, v...)
	}
}

func (this *Log) Printf(level Level, format string, v ...any) {
	this.Print(level, fmt.Sprintf(format, v...))
}

func (this *Log) Trace(v ...any) {
	this.Print(LevelTrace, v...)
}

func (this *Log) Tracef(format string, v ...any) {
	this.Printf(LevelTrace, format, v...)
}

func (this *Log) Debug(v ...any) {
	this.Print(LevelDebug, v...)
}

func (this *Log) Debugf(format string, v ...any) {
	this.Printf(LevelDebug, format, v...)
}

func (this *Log) Info(v ...any) {
	this.Print(LevelInfo, v...)
}

func (this *Log) Infof(format string, v ...any) {
	this.Printf(LevelInfo, format, v...)
}

func (this *Log) Warning(v ...any) {
	this.Print(LevelWarning, v...)
}

func (this *Log) Warningf(format string, v ...any) {
	this.Printf(LevelWarning, format, v...)
}

func (this *Log) Error(v ...any) {
	this.Print(LevelError, v...)
}

func (this *Log) Errorf(format string, v ...any) {
	this.Printf(LevelError, format, v...)
}

func (this *Log) Critical(v ...any) {
	this.Print(LevelCritical, v...)
}

func (this *Log) Criticalf(format string, v ...any) {
	this.Printf(LevelCritical, format, v...)
}

func (this *Log) Stat() {
	this.Info(this.ctx.statistics())
}

func (this *Log) Recover() {
	if r := recover(); r != nil {
		this.Critical(r)
	}
}

func (this *Log) Query(f Filter) (lines []string, err error) {
	return this.ctx.query(f)
}
