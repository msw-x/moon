package ulog

import "fmt"

type Log struct {
	enable        bool
	level         Level
	prefix        string
	lifetimeLevel *Level
}

func New(prefix string) *Log {
	return &Log{
		enable: true,
		prefix: prefix,
	}
}

func (this *Log) Close() {
	if this.lifetimeLevel != nil {
		this.Print(*this.lifetimeLevel, "~")
	}
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
		space := ""
		if !ctx.conf.splitArgs {
			space = " "
		}
		v = append([]any{fmt.Sprintf("<%s>%s", this.prefix, space)}, v...)
		Print(level, v...)
	}
}

func (this *Log) Printf(level Level, format string, v ...any) {
	this.Print(level, fmt.Sprintf(format, v...))
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
	Print(LevelCritical, v...)
}

func (this *Log) Criticalf(format string, v ...any) {
	Printf(LevelCritical, format, v...)
}
