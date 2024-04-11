package app

import (
	"fmt"
	"os"
	"time"

	"github.com/msw-x/moon"
	"github.com/msw-x/moon/hw"
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
	hws := hw.GetStatus()
	log.Info("cpu:", hws.Cpu())
	log.Info("ram:", hws.Ram())
	log.Info("disk:", hws.Disk())
	logTimezone(log)
	log.Info("name:", Name())
	dir := Dir()
	pwd := Pwd()
	log.Info("path:", dir)
	if dir != pwd && pwd != "" {
		log.Info("pwd:", pwd)
	}
	switch f := fn.(type) {
	case func():
		f()
	case func(*ulog.Log):
		f(log)
	default:
		panic("app run: invalid callback")
	}
}

func logTimezone(log *ulog.Log) {
	tzName, tzOffset := time.Now().Zone()
	logTzName, logTzOffset := log.Timezone()
	var tz string
	if tzOffset != logTzOffset {
		tz = fmt.Sprintf(" => %s (%+d)", logTzName, logTzOffset/(60*60))
	}
	log.Infof("timezone: %s (%+d)%s", tzName, tzOffset/(60*60), tz)
}
