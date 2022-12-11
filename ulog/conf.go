package ulog

import (
	"reflect"

	"github.com/msw-x/moon"
	"github.com/msw-x/moon/parse"
)

type Conf struct {
	Level     any
	Console   bool
	File      string
	Dir       string
	Append    bool
	AppName   string
	GoID      bool
	SplitArgs any

	level     Level
	splitArgs bool
}

func (this *Conf) init() {
	this.level = initLevel(this.Level)
	this.splitArgs = initSplitArgs(this.SplitArgs)
}

func initLevel(a any) Level {
	switch v := a.(type) {
	case nil:
		return LevelDefault
	case Level:
		return v
	case string:
		return ParseLevel(v)
	}
	moon.Panic("invalid ulog level type:", reflect.TypeOf(a))
	return -1
}

func initSplitArgs(a any) bool {
	switch v := a.(type) {
	case nil:
		return true
	case bool:
		return v
	case string:
		return parse.BoolStrict(v)
	}
	moon.Panic("invalid ulog level type:", reflect.TypeOf(a))
	return false
}
