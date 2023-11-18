package app

import (
	"flag"

	"github.com/msw-x/moon/ufmt"
	"github.com/msw-x/moon/ufs"
)

func ParseCmdLine(version string) (confFile string, ok bool) {
	defConf := ufs.ReplaceExt(Executable(), ".conf")
	showVersion := flag.Bool("v", false, "show version")
	flag.StringVar(&confFile, "c", defConf, "config file")
	flag.Parse()
	if *showVersion {
		ufmt.Print("version:", version)
	}
	ok = !*showVersion
	return
}
