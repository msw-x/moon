package ulog

import (
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"time"
)

type dir struct {
	path  string
	files []fs.FileInfo
}

func (o dir) empty() bool {
	return len(o.files) == 0
}

func (o dir) time() time.Time {
	return o.files[0].ModTime()
}

func (o dir) size() (n int64) {
	for _, f := range o.files {
		n += f.Size()
	}
	return
}

func (o dir) sort() {
	sort.Slice(o.files, func(i, j int) bool {
		return o.files[i].ModTime().Before(o.files[j].ModTime())
	})
}

type dirs []dir

func (o dirs) sort() {
	for _, d := range o {
		d.sort()
	}
	sort.Slice(o, func(i, j int) bool {
		return o[i].time().Before(o[j].time())
	})
}

func (o dirs) size() (n int64) {
	for _, d := range o {
		n += d.size()
	}
	return
}

func (o dirs) count() int {
	return len(o)
}

func (o dirs) removeByCount(n int) {
	for _, d := range o[0:n] {
		os.RemoveAll(d.path)
	}
}

func (o dirs) removeBySize(n int64) {
	for _, d := range o {
		if n <= 0 {
			break
		}
		n -= d.size()
		os.RemoveAll(d.path)
	}
}

func scanDirs(p string) (l dirs) {
	folders, _ := os.ReadDir(p)
	for _, j := range folders {
		if j.IsDir() {
			var d dir
			d.path = filepath.Join(p, j.Name())
			files, _ := os.ReadDir(d.path)
			for _, k := range files {
				if !k.IsDir() {
					if i, err := k.Info(); err == nil {
						d.files = append(d.files, i)
					}
				}
			}
			if !d.empty() {
				l = append(l, d)
			}
		}
	}
	l.sort()
	return
}
