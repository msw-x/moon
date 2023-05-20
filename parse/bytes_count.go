package parse

import (
	"math"

	"github.com/msw-x/moon/uerr"
)

func BytesCount(s string) (v uint64, err error) {
	var multiplier uint64
	if len(s) > 0 {
		var level int
		suf := s[len(s)-1:]
		switch suf {
		case "K":
			level = 1
		case "M":
			level = 2
		case "G":
			level = 3
		case "T":
			level = 4
		}
		if level > 0 {
			s = s[:len(s)-1]
		}
		multiplier = uint64(math.Pow(1024, float64(level)))
	}
	v, err = Uint64(s)
	v *= multiplier
	return
}

func BytesCountStrict(s string) uint64 {
	n, err := BytesCount(s)
	uerr.Strict(err, "parse bytes count")
	return n
}
