package ufs

import (
	"os"
	"time"

	"github.com/msw-x/moon/uerr"
)

func RemoveStrict(path string) {
	uerr.Strict(Remove(path), "fs remove")
}

func RenameStrict(path, newName string) {
	uerr.Strict(Rename(path, newName), "fs rename")
}

func IsDirStrict(path string) bool {
	yes, err := IsDir(path)
	uerr.Strict(err, "fs is dir")
	return yes
}

func MakeDirStrict(dir string) {
	uerr.Strict(MakeDir(dir), "fs make dir")
}

func FileSizeStrict(path string) int64 {
	s, err := FileSize(path)
	uerr.Strict(err, "file size")
	return s
}

func FileModifyTimeStrict(path string) time.Time {
	tm, err := FileModifyTime(path)
	uerr.Strict(err, "file modify time")
	return tm
}

func DirSizeStrict(path string) int64 {
	s, err := DirSize(path)
	uerr.Strict(err, "dir size")
	return s
}

func CreateStrict(filename string) *os.File {
	f, err := Create(filename)
	uerr.Strict(err, "fs create")
	return f
}

func ChmodStrict(name string, mode os.FileMode) {
	uerr.Strict(Chmod(name, mode), "chmod")
}
