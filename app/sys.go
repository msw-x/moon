package app

import (
	"fmt"
	"moon"
	"moon/fs"
	"os"
	"path"
	"path/filepath"
	"runtime"
)

func Executable() string {
	p, err := os.Executable()
	moon.Check(err, "executable")
	return path.Clean(p)
}

func Name() string {
	return fs.RemoveExt(filepath.Base(Executable()))
}

func Dir() string {
	return path.Dir(Executable())
}

func Pid() int {
	return os.Getpid()
}

func OS() string {
	return runtime.GOOS
}

func Arch() string {
	architecture := runtime.GOARCH
	addressModel := ""
	if runtime.GOARCH == "amd64" {
		architecture = "x86"
		addressModel = "64"
	} else if runtime.GOARCH == "arm64" {
		architecture = "arm"
		addressModel = "64"
	} else if runtime.GOARCH == "arm" {
		architecture = "arm"
		addressModel = "32"
	} else {
		return architecture
	}
	return fmt.Sprintf("%s-%s", architecture, addressModel)
}
