package ulog

import (
	"github.com/msw-x/moon/fs"
)

type Filter struct {
	Level Level
	Last  int
}

func QueryFromFile(filename string, f Filter) (ret []string, err error) {
	var concatenate bool
	err = fs.ForEachLine(filename, func(line string) {
		l := selectLevel(line)
		if l == "" {
			if concatenate {
				n := len(ret) - 1
				ret[n] = ret[n] + "\n" + line
			}
		} else {
			if concatenate = f.Level.Laconic() == l; concatenate {
				ret = append(ret, line)
			}
		}
	})
	count := len(ret)
	overflow := count - f.Last - 1
	if f.Last > 0 && overflow > 0 {
		ret = ret[overflow:len(ret)]
	}
	return ret, err
}

func Query(f Filter) (lines []string, err error) {
	return ctx.query(f)
}

func selectLevel(line string) string {
	const fi = 25
	const ei = 29
	if len(line) > ei && line[fi] == '[' && line[ei] == ']' {
		return line[fi+1 : ei]
	}
	return ""
}
