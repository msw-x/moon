package app

import (
	"os"

	"github.com/msw-x/moon"
	"github.com/msw-x/moon/uerr"
	"github.com/msw-x/moon/ulog"
)

func Run[UserConf any](version string, fn func(UserConf)) {
	defer exit()
	defer ulog.Close()
	defer uerr.Recover(fatal)
	if confFile, ok := ParseCmdLine(version); ok {
		conf := LoadConf[Conf](confFile)
		run(version, conf.Log(), func(log *ulog.Log) {
			log.Info("conf:", confFile)
			userConf := LoadConf[UserConf](confFile)
			fn(userConf)
		})
	}
}

func RunJust(version string, opts ulog.Options, fn func()) {
	defer exit()
	defer ulog.Close()
	defer uerr.Recover(fatal)
	run(version, opts, fn)
}

func WaitInterrupt() {
	s := moon.WaitInterrupt()
	log.Info("signal:", s)
}

var log *ulog.Log
var exitCode int

func fatal(s string) {
	exitCode = 1
	ulog.Critical(s)
}

func exit() {
	if exitCode > 0 {
		os.Exit(1)
	}
}

func run(version string, opts ulog.Options, fn any) {
	ulog.Init(opts)
	log = ulog.New("app")
	defer func() {
		log.Info(ulog.Stat())
		log.Info("shutdown")
		log.Close()
	}()
	defer uerr.Recover(func(e string) {
		exitCode = 1
		log.Critical(e)
	})
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
}
