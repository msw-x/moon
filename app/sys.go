package app

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/msw-x/moon/uerr"
	"github.com/msw-x/moon/ufs"
)

func Executable() string {
	p, err := os.Executable()
	uerr.Strict(err, "app executable")
	return path.Clean(p)
}

func Name() string {
	return ufs.RemoveExt(filepath.Base(Executable()))
}

func Dir() string {
	return path.Dir(Executable())
}

func Pwd() string {
	s, _ := os.Getwd()
	return s
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
