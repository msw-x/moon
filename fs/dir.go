package fs

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/msw-x/moon"
)

func ReadDir(dir string, ignorAccessDenied bool) []os.FileInfo {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		if !(os.IsPermission(err) && ignorAccessDenied) {
			moon.Check(err, "read dir")
		}
	}
	return files
}

type List struct {
	Files             bool
	Folders           bool
	Extension         string
	ExtensionList     []string
	Recursive         bool
	IgnorAccessDenied bool
	RelativePath      bool

	level int
}

func (this List) Get(dir string) []string {
	content := ReadDir(dir, this.IgnorAccessDenied)
	filter := func(c os.FileInfo) bool {
		if this.Folders {
			return c.IsDir()
		} else if this.Files {
			if c.Mode().IsRegular() {
				if this.Extension == "" {
					if len(this.ExtensionList) == 0 {
						return true
					} else {
						for _, ext := range this.ExtensionList {
							if EqualExt(Ext(c.Name()), ext) {
								return true
							}
						}
						return false
					}
				} else {
					return EqualExt(Ext(c.Name()), this.Extension)
				}
			}
		}
		return false
	}
	var result []string
	for _, c := range content {
		if filter(c) {
			result = append(result, c.Name())
		}
		if this.Recursive && c.IsDir() {
			subResult := this.sub().Get(path.Join(dir, c.Name()))
			for n, sub := range subResult {
				subResult[n] = path.Join(c.Name(), sub)
			}
			result = append(result, subResult...)
		}
	}
	if !this.RelativePath && this.level == 0 {
		for n, i := range result {
			result[n] = path.Join(dir, i)
		}
	}
	return result
}

func (this List) sub() List {
	this.level++
	return this
}
