package ulog

import (
	"strings"

	"github.com/msw-x/moon/fs"
)

type Filter struct {
	Level Level
	Last  int
}

func QueryFromFile(filename string, f Filter) ([]string, error) {
	s, err := fs.ReadString(filename)
	var ret []string
	if err == nil {
		lines := strings.Split(s, "\n")
		for _, line := range lines {
			i := strings.Index(line, "[")
			if i > 0 {
				s := line[i+1:]
				i = strings.Index(s, "]")
				lvl := s[:i]
				if f.Level.Laconic() == lvl {
					ret = append(ret, line)
				}
			}
		}
	}
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
