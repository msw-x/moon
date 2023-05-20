package ulog

import (
	"reflect"

	"github.com/msw-x/moon/parse"
	"github.com/msw-x/moon/uerr"
)

type Options struct {
	Level         any
	Console       bool
	File          string
	Dir           string
	Append        bool
	AppName       string
	GoID          bool
	CrtStdErr     bool
	SplitArgs     any
	FileSizeLimit uint64

	level     Level
	splitArgs bool
}

func (o *Options) init() {
	o.level = initLevel(o.Level)
	o.splitArgs = initSplitArgs(o.SplitArgs)
}

func (o *Options) useDir() bool {
	return o.File == "" && o.Dir != ""
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
	uerr.Panic("invalid ulog level type:", reflect.TypeOf(a))
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
	uerr.Panic("invalid ulog level type:", reflect.TypeOf(a))
	return false
}
