package app

import (
	"path"
	"reflect"
	"strings"

	"github.com/msw-x/moon"
	"github.com/msw-x/moon/parse"
	"github.com/msw-x/moon/ulog"

	"github.com/BurntSushi/toml"
)

func LoadConf[Conf any](filename string) Conf {
	var conf Conf
	_, err := toml.DecodeFile(filename, &conf)
	moon.Strict(err, "load conf")
	rv := reflect.ValueOf(&conf).Elem()
	for i := 0; i < rv.NumField(); i++ {
		name := rv.Type().Field(i).Name
		if strings.HasSuffix(name, "Dir") || strings.HasSuffix(name, "File") {
			val := rv.Field(i)
			if val.Kind() == reflect.String {
				s := val.String()
				if s != "" {
					s = path.Clean(s)
					if !path.IsAbs(s) {
						s = path.Join(Dir(), s)
					}
					if val.CanSet() {
						val.SetString(s)
					}
				}
			}
		}
	}
	return conf
}

type Conf struct {
	LogLevel    string
	LogDir      string
	LogFile     string
	LogConsole  bool
	LogGoID     bool
	LogFileSize string
}

func (o *Conf) Log() ulog.Options {
	opts := ulog.Options{
		Level:   ulog.ParseLevel(o.LogLevel),
		Console: o.LogConsole,
		File:    o.LogFile,
		Dir:     o.LogDir,
		GoID:    o.LogGoID,
	}
	if o.LogFileSize != "" {
		opts.FileSizeLimit = parse.BytesCountStrict(o.LogFileSize)
	}
	return opts
}
