package ulog

import (
	"os"
)

func link(targetName, linkName string) (err error) {
	var i os.FileInfo
	i, err = os.Lstat(linkName)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Symlink(targetName, linkName)
		}
		return
	}
	if i.Mode()&os.ModeSymlink == 0 {
		return
	}
	var target string
	target, err = os.Readlink(linkName)
	if err != nil {
		return
	}
	if target == targetName {
		return
	}
	err = os.Remove(linkName)
	if err != nil {
		return
	}
	err = os.Symlink(targetName, linkName)
	return
}
