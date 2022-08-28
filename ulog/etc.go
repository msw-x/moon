package ulog

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func Panicf(format string, v ...any) {
	panic(fmt.Sprintf("ulog: "+format, v...))
}

func AppName() string {
	path, err := os.Executable()
	if err == nil {
		base := filepath.Base(path)
		return strings.TrimSuffix(base, filepath.Ext(base))
	}
	return ""
}
