package ulog

import "github.com/msw-x/moon/uerr"

type Level int

const (
	LevelTrace Level = iota + 1
	LevelDebug
	LevelInfo
	LevelWarning
	LevelError
	LevelCritical
)

const LevelDefault = LevelInfo

func (l Level) Laconic() string {
	switch l {
	case LevelTrace:
		return "trc"
	case LevelDebug:
		return "dbg"
	case LevelInfo:
		return "inf"
	case LevelWarning:
		return "wrn"
	case LevelError:
		return "err"
	case LevelCritical:
		return "crt"
	}
	return ""
}

func (l Level) String() string {
	switch l {
	case LevelTrace:
		return "trace"
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelWarning:
		return "warning"
	case LevelError:
		return "error"
	case LevelCritical:
		return "critical"
	}
	return ""
}

func ParseLevel(s string) Level {
	switch s {
	case "":
		return LevelDefault
	case "trace", "trc":
		return LevelTrace
	case "debug", "dbg":
		return LevelDebug
	case "info", "inf":
		return LevelInfo
	case "warning", "wrn":
		return LevelWarning
	case "error", "err":
		return LevelError
	case "critical", "crt":
		return LevelCritical
	}
	uerr.Panic("unknown log level:", s)
	return -1
}
