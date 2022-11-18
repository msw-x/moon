package fs

import (
	"os"
	"time"

	"github.com/msw-x/moon"
)

func RemoveStrict(path string) {
	moon.Strict(Remove(path), "fs remove")
}

func RenameStrict(path, newName string) {
	moon.Strict(Rename(path, newName), "fs rename")
}

func IsDirStrict(path string) bool {
	yes, err := IsDir(path)
	moon.Strict(err, "fs is dir")
	return yes
}

func MakeDirStrict(dir string) {
	moon.Strict(MakeDir(dir), "fs make dir")
}

func FileSizeStrict(path string) int64 {
	s, err := FileSize(path)
	moon.Strict(err, "file size")
	return s
}

func FileModifyTimeStrict(path string) time.Time {
	tm, err := FileModifyTime(path)
	moon.Strict(err, "file modify time")
	return tm
}

func DirSizeStrict(path string) int64 {
	s, err := DirSize(path)
	moon.Strict(err, "dir size")
	return s
}

func CreateStrict(filename string) *os.File {
	f, err := Create(filename)
	moon.Strict(err, "fs create")
	return f
}

func ChmodStrict(name string, mode os.FileMode) {
	moon.Strict(Chmod(name, mode), "chmod")
}
