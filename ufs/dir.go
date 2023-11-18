package ufs

import (
	"io/ioutil"
	"os"

	"github.com/msw-x/moon/uerr"
)

func ReadDir(dir string, ignorAccessDenied bool) (files []os.FileInfo, err error) {
	files, err = ioutil.ReadDir(dir)
	if err != nil && os.IsPermission(err) && ignorAccessDenied {
		err = nil
	}
	return
}

func ReadDirStrict(dir string, ignorAccessDenied bool) []os.FileInfo {
	files, err := ReadDir(dir, ignorAccessDenied)
	uerr.Strict(err, "read dir")
	return files
}
