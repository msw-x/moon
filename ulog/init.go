package ulog

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/msw-x/moon/ufs"
)

func Init(opts Options) {
	ctx.init(opts)
}

func InitConsole() {
	Init(Options{
		Console: true,
	})
}

func InitFile(filename string) {
	Init(Options{
		File: filename,
	})
}

func InitDir(dirname string) {
	Init(Options{
		Dir: dirname,
	})
}

func Close() {
	ctx.close()
}

func GenFilename(ts time.Time, dir, app string) string {
	const ext = ".log"
	subDir := ts.Format("2006-01-02")
	logDir := path.Join(dir, subDir)
	basename := ts.Format("2006-01-02--15-04-05") + "@" + app
	filename := path.Join(logDir, basename+ext)
	if ufs.Exist(filename) {
		filename = path.Join(logDir, fmt.Sprintf("%s.%d%s", basename, os.Getpid(), ext))
	}
	return filename
}

func OpenFile(filename string, append bool) *os.File {
	dir := filepath.Dir(filename)
	if dir == "." {
		dir = ""
	}
	if dir != "" {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			Panicf("make log directory: %v", err)
		}
	}
	flag := os.O_WRONLY | os.O_CREATE
	if append {
		flag |= os.O_APPEND
	} else {
		flag |= os.O_TRUNC
	}
	file, err := os.OpenFile(filename, flag, 0600)
	if err != nil {
		Panicf("open log file: %v", err)
	}
	return file
}

func SetHook(hook func(Message)) {
	ctx.hook = hook
}

func init() {
	Init(Options{
		CrtStdErr: true,
	})
}
