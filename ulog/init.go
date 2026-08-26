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

func GenFilename(ts time.Time, dir, app string, fileLink bool, dayLink string) (fdir, filename, flink, dlink string) {
	const ext = ".log"
	fdir = path.Join(dir, ts.Format("2006-01-02"))
	base := ts.Format("2006-01-02--15-04-05") + "@" + app
	filename = path.Join(fdir, base+ext)
	if ufs.Exist(filename) {
		filename = path.Join(fdir, fmt.Sprintf("%s.%d%s", base, os.Getpid(), ext))
	}
	if fileLink {
		flink = path.Join(dir, app+ext)
	}
	if dayLink != "" {
		dlink = path.Join(dir, dayLink)
	}
	return
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
