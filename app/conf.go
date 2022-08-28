package app

import (
	"github.com/msw-x/moon"
	"github.com/msw-x/moon/ulog"

	"github.com/BurntSushi/toml"
)

func LoadConf[Conf any](filename string) Conf {
	var conf Conf
	_, err := toml.DecodeFile(filename, &conf)
	moon.Check(err, "load conf")
	return conf
}

type Conf struct {
	LogLevel   string
	LogDir     string
	LogFile    string
	LogConsole bool
}

func (this *Conf) Log() ulog.Conf {
	return ulog.Conf{
		Level:   ulog.ParseLevel(this.LogLevel),
		Console: this.LogConsole,
		File:    this.LogFile,
		Dir:     this.LogDir,
	}
}
