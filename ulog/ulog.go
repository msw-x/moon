package ulog

import (
	"fmt"
	"os"
)

func Print(level Level, v ...any) {
	print(&ctx, level, v...)
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
	return ctx.statistics()
}

func Recover() {
	if r := recover(); r != nil {
		Critical(r)
	}
}

func print(ctx *context, level Level, v ...any) {
	if level >= ctx.opts.level {
		m := NewMessage(ctx, level, v...)
		printMessage(ctx, level, m)
		if ctx.hook != nil {
			ctx.hook(m)
		}
	}
}

func printMessage(ctx *context, level Level, m Message) {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()
	ctx.stat.Push(level, m.Size())
	if ctx.opts.Console || level == LevelCritical {
		if level >= LevelError {
			if ctx.opts.Console {
				printStdErr(m.Format())
			} else {
				if ctx.opts.CrtStdErr {
					printStdErr(m.Text)
				}
			}
		} else {
			fmt.Print(m.Format())
		}
	}
	if ctx.file != nil {
		ctx.rotate(m.Size())
		ctx.fileSize += uint64(m.Size())
		ctx.file.WriteString(m.Format())
	}
}

func printStdErr(text string) {
	fmt.Fprint(os.Stderr, text)
}
