package fs

import (
	"os"
	"path/filepath"
	"time"
)

func Exist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func Remove(path string) error {
	return os.RemoveAll(path)
}

func Rename(path, newName string) error {
	return os.Rename(path, newName)
}

func IsDir(path string) (bool, error) {
	fi, err := os.Stat(path)
	return fi.Mode().IsDir(), err
}

func MakeDir(dir string) error {
	return os.MkdirAll(dir, os.ModePerm)
}

func FileSize(path string) (int64, error) {
	fi, err := os.Stat(path)
	return fi.Size(), err
}

func FileModifyTime(path string) (time.Time, error) {
	fi, err := os.Stat(path)
	return fi.ModTime(), err
}

func DirSize(path string) (size int64, err error) {
	err = filepath.Walk(path,
		func(_ string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				size += info.Size()
			}
			return err
		})
	return
}

func Create(filename string) (f *os.File, err error) {
	err = MakeDir(filepath.Dir(filename))
	if err == nil {
		f, err = os.Create(filename)
	}
	return
}

func Chmod(name string, mode os.FileMode) error {
	return os.Chmod(name, mode)
}
