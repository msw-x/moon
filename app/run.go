package app

import (
	"os"

	"github.com/msw-x/moon"
	"github.com/msw-x/moon/ulog"
)

func Run[UserConf any](version string, fn func(UserConf)) {
	defer ulog.Close()
	defer moon.Recover(Fatal)
	if confFile, ok := ParseCmdLine(version); ok {
		conf := LoadConf[Conf](confFile)
		run(version, conf.Log(), func(log *ulog.Log) {
			log.Info("conf:", confFile)
			userConf := LoadConf[UserConf](confFile)
			fn(userConf)
		})
	}
}

func RunJust(version string, conf ulog.Conf, fn func()) {
	defer ulog.Close()
	defer moon.Recover(Fatal)
	run(version, conf, fn)
}

func Go(fn func()) {
	go func() {
		defer moon.Recover(func(e string) {
			ulog.Critical(e)
		})
		fn()
	}()
}

func Fatal(s string) {
	ulog.Critical(s)
	ulog.Close()
	os.Exit(1)
}

func run(version string, conf ulog.Conf, fn any) {
	ulog.Init(conf)
	log := ulog.New("app")
	defer log.Close()
	log.Info("startup")
	log.Info("version:", version)
	log.Info("pid:", Pid())
	log.Info("os:", OS())
	log.Info("arch:", Arch())
	log.Info("name:", Name())
	log.Info("path:", Dir())
	switch f := fn.(type) {
	case func():
		f()
	case func(*ulog.Log):
		f(log)
	default:
		panic("app run: invalid callback")
	}
	log.Info(ulog.Stat())
	log.Info("shutdown")
}
