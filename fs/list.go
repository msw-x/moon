package fs

import (
	"os"
	"path"

	"github.com/msw-x/moon/uerr"
)

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

func (o List) Get(dir string) (result []string, err error) {
	content, err := ReadDir(dir, o.IgnorAccessDenied)
	filter := func(c os.FileInfo) bool {
		if o.Folders {
			return c.IsDir()
		} else if o.Files {
			if c.Mode().IsRegular() {
				if o.Extension == "" {
					if len(o.ExtensionList) == 0 {
						return true
					} else {
						for _, ext := range o.ExtensionList {
							if EqualExt(Ext(c.Name()), ext) {
								return true
							}
						}
						return false
					}
				} else {
					return EqualExt(Ext(c.Name()), o.Extension)
				}
			}
		}
		return false
	}
	for _, c := range content {
		if filter(c) {
			result = append(result, c.Name())
		}
		if o.Recursive && c.IsDir() {
			var subResult []string
			subResult, err = o.sub().Get(path.Join(dir, c.Name()))
			if err != nil {
				return
			}
			for n, sub := range subResult {
				subResult[n] = path.Join(c.Name(), sub)
			}
			result = append(result, subResult...)
		}
	}
	if !o.RelativePath && o.level == 0 {
		for n, i := range result {
			result[n] = path.Join(dir, i)
		}
	}
	return
}

func (o List) GetStrict(dir string) []string {
	r, err := o.Get(dir)
	uerr.Strict(err, "fs list")
	return r
}

func (o List) sub() List {
	o.level++
	return o
}
