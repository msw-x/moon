package ulog

import (
	"fmt"
	"os"
	"time"

	"github.com/msw-x/moon"
	"github.com/msw-x/moon/ufmt"
)

func Print(level Level, v ...any) {
	if level >= ctx.conf.level {
		m := NewMessage(level, v...)
		ctx.mutex.Lock()
		defer ctx.mutex.Unlock()
		ctx.stat.Push(level, m.Size())
		if ctx.conf.Console || level == LevelCritical {
			if level >= LevelError {
				if ctx.conf.Console {
					fmt.Fprint(os.Stderr, m.Format())
				} else {
					fmt.Fprint(os.Stderr, m.Text)
				}
			} else {
				fmt.Print(m.Format())
			}
		}
		if ctx.file != nil {
			ctx.file.WriteString(m.Format())
		}
	}
}

func Printf(level Level, format string, v ...any) {
	Print(level, fmt.Sprintf(format, v...))
}

func Trace(v ...any) {
	Print(LevelTrace, v...)
}

func Tracef(format string, v ...any) {
	Printf(LevelTrace, format, v...)
}

func Debug(v ...any) {
	Print(LevelDebug, v...)
}

func Debugf(format string, v ...any) {
	Printf(LevelDebug, format, v...)
}

func Info(v ...any) {
	Print(LevelInfo, v...)
}

func Infof(format string, v ...any) {
	Printf(LevelInfo, format, v...)
}

func Warning(v ...any) {
	Print(LevelWarning, v...)
}

func Warningf(format string, v ...any) {
	Printf(LevelWarning, format, v...)
}

func Error(v ...any) {
	Print(LevelError, v...)
}

func Errorf(format string, v ...any) {
	Printf(LevelError, format, v...)
}

func Critical(v ...any) {
	Print(LevelCritical, v...)
}

func Criticalf(format string, v ...any) {
	Printf(LevelCritical, format, v...)
}

func Stat() string {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()
	tm := time.Since(ctx.inited)
	dur := moon.DurationToTime(tm)
	var text string
	if tm < time.Second {
		text = fmt.Sprintf("%d ms", dur.Milliseconds)
	} else {
		text = dur.FormatDays()
	}
	text = fmt.Sprintf("%s | %s", text, ufmt.ByteSize(ctx.stat.Size))
	if ctx.conf.GoID {
		text = fmt.Sprintf("%s go[%s]", text, ufmt.WideInt(len(ctx.mapid)))
	}
	add := func(level Level, count uint) {
		if count > 0 {
			text = fmt.Sprintf("%s %v[%s]", text, level.Laconic(), ufmt.WideInt(count))
		}
	}
	add(LevelTrace, ctx.stat.Trace)
	add(LevelDebug, ctx.stat.Debug)
	add(LevelInfo, ctx.stat.Info)
	add(LevelWarning, ctx.stat.Warning)
	add(LevelError, ctx.stat.Error)
	add(LevelCritical, ctx.stat.Critical)
	return text
}
