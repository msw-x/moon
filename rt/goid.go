package rt

import (
	"bytes"
	"runtime"
	"strconv"
	"sync"
)

func GoroutineID() int {
	bp := littleBuf.Get().(*[]byte)
	defer littleBuf.Put(bp)
	b := *bp
	b = b[:runtime.Stack(b, false)]
	prefix := []byte("goroutine ")
	if bytes.HasPrefix(b, prefix) {
		b = bytes.TrimPrefix(b, prefix)
		i := bytes.IndexByte(b, ' ')
		if i > 0 {
			b = b[:i]
			n, err := strconv.Atoi(string(b))
			if err == nil {
				return n
			}
		}
	}
	return -1
}

var littleBuf = sync.Pool{
	New: func() any {
		buf := make([]byte, 64)
		return &buf
	},
}
