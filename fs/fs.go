package fs

import (
	"moon"
	"os"
	"path/filepath"
	"time"
)

func Exist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func Remove(path string) {
	moon.Check(os.RemoveAll(path), "fs remove")
}

func Rename(path, newName string) {
	err := os.Rename(path, newName)
	moon.Check(err, "fs rename")
}

func IsDir(path string) bool {
	fi, err := os.Stat(path)
	moon.Check(err)
	return fi.Mode().IsDir()
}

func MakeDir(dir string) {
	moon.Check(os.MkdirAll(dir, os.ModePerm), "make dir")
}

func FileSize(path string) int64 {
	fi, err := os.Stat(path)
	moon.Check(err, "file size")
	return fi.Size()
}

func FileModTime(path string) time.Time {
	fi, err := os.Stat(path)
	moon.Check(err, "file modify time")
	return fi.ModTime()
}

func DirSize(path string) int64 {
	var size int64
	err := filepath.Walk(path,
		func(_ string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				size += info.Size()
			}
			return err
		})
	moon.Check(err, "dir size")
	return size
}

func CreateOpt(filename string) (*os.File, error) {
	MakeDir(filepath.Dir(filename))
	return os.Create(filename)
}

func Create(filename string) *os.File {
	f, err := CreateOpt(filename)
	moon.Check(err, "fs create")
	return f
}

func Chmod(name string, mode os.FileMode) {
	moon.Check(os.Chmod(name, mode), "chmod")
}
