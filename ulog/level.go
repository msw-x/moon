package ulog

import "github.com/msw-x/moon"

type Level int

const (
	LevelDebug Level = iota + 1
	LevelInfo
	LevelWarning
	LevelError
	LevelCritical
)

const LevelDefault = LevelInfo

func (l Level) Laconic() string {
	switch l {
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
	case "debug":
		return LevelDebug
	case "info":
		return LevelInfo
	case "warning":
		return LevelWarning
	case "error":
		return LevelError
	case "critical":
		return LevelCritical
	}
	moon.Panicf("unknown log level:", s)
	return -1
}
