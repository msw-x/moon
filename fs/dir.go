package fs

import (
	"io/ioutil"
	"moon"
	"os"
	"path"
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

type ListCtx struct {
	Files             bool
	Folders           bool
	Extension         string
	ExtensionList     []string
	Recursive         bool
	IgnorAccessDenied bool
	RelativePath      bool

	level int
}

func List(dir string, ctx ListCtx) []string {
	content := ReadDir(dir, ctx.IgnorAccessDenied)
	filter := func(c os.FileInfo) bool {
		if ctx.Folders {
			return c.IsDir()
		} else if ctx.Files {
			if c.Mode().IsRegular() {
				if ctx.Extension == "" {
					if len(ctx.ExtensionList) == 0 {
						return true
					} else {
						for _, ext := range ctx.ExtensionList {
							if EqualExt(Ext(c.Name()), ext) {
								return true
							}
						}
						return false
					}
				} else {
					return EqualExt(Ext(c.Name()), ctx.Extension)
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
		if ctx.Recursive && c.IsDir() {
			nextCtx := ctx
			nextCtx.level++
			subResult := List(path.Join(dir, c.Name()), nextCtx)
			for n, sub := range subResult {
				subResult[n] = path.Join(c.Name(), sub)
			}
			result = append(result, subResult...)
		}
	}
	if !ctx.RelativePath && ctx.level == 0 {
		for n, i := range result {
			result[n] = path.Join(dir, i)
		}
	}
	return result
}
