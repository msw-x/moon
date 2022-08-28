package app

import (
	"github.com/msw-x/moon"
	"github.com/msw-x/moon/ulog"
	"os"
)

func Run[UserConf any](version string, run func(UserConf)) {
	defer ulog.Close()
	defer moon.Recover(Fatal)
	if confFile, ok := ParseCmdLine(version); ok {
		conf := LoadConf[Conf](confFile)
		ulog.Init(conf.Log())
		log := ulog.NewLog("app")
		defer log.Close()
		log.Info("startup")
		log.Info("version:", version)
		log.Info("pid:", Pid())
		log.Info("os:", OS())
		log.Info("arch:", Arch())
		log.Info("name:", Name())
		log.Info("path:", Dir())
		log.Info("conf:", confFile)
		userConf := LoadConf[UserConf](confFile)
		run(userConf)
		log.Info("shutdown")
	}
}

func Go(fn func()) {
	go func() {
		defer moon.Recover(Fatal)
		fn()
	}()
}

func Fatal(s string) {
	ulog.Critical(s)
	os.Exit(1)
}
