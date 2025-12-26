package rt

import (
	"fmt"
	"runtime"
	"strings"
)

func StackTrace(n int) string {
	n++ // for ignore this function
	return stackTrace(n, "%s\n  %s:%d\n")
}

func stackTrace(n int, f string) (s string) {
	n++ // for ignore this function
	for i := 0; ; i++ {
		if i < n {
			continue
		}
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		fn := runtime.FuncForPC(pc)
		s += fmt.Sprintf(f, fn.Name(), file, line)
	}
	s = strings.TrimSuffix(s, "\n")
	return
}
