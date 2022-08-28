package ulog

import (
	"fmt"
	"moon/fs"
	"os"
	"path"
	"path/filepath"
	"time"
)

func Init(conf Conf) {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()
	conf.init()
	if ctx.conf.File != conf.File || ctx.conf.Dir != conf.Dir {
		if ctx.file != nil {
			ctx.file.Close()
			ctx.file = nil
		}
		filename := conf.File
		if filename == "" && conf.Dir != "" {
			filename = GenFilename(conf.Dir, AppName())
		}
		if filename != "" {
			ctx.file = OpenFile(filename, conf.Append)
		}
	}
	ctx.conf = conf
	ctx.maxid = 2
	ctx.mapid = make(map[int]bool)
	ctx.inited = time.Now()
}

func InitConsole() {
	Init(Conf{
		Console: true,
	})
}

func InitFile(filename string) {
	Init(Conf{
		File: filename,
	})
}

func InitDir(dirname string) {
	Init(Conf{
		Dir: dirname,
	})
}

func Close() {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()
	if ctx.file != nil {
		ctx.file.Close()
		ctx.file = nil
	}
}

func GenFilename(dir, app string) string {
	const ext = ".log"
	subDir := time.Now().Format("2006-01-02")
	logDir := path.Join(dir, subDir)
	basename := time.Now().Format("15-04-05") + "@" + app
	filename := path.Join(logDir, basename+ext)
	if fs.Exist(filename) {
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

func init() {
	Init(Conf{})
}
