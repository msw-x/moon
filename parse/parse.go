package parse

import (
	"strconv"

	"github.com/msw-x/moon"
)

func Int(s string) int {
	i, err := strconv.Atoi(s)
	moon.Check(err, "parse int")
	return i
}

func Int64(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	moon.Check(err, "parse int")
	return i
}

func Uint64(s string) uint64 {
	i, err := strconv.ParseUint(s, 10, 64)
	moon.Check(err, "parse int")
	return i
}

func Float32(s string) float32 {
	f, err := strconv.ParseFloat(s, 32)
	moon.Check(err, "parse float32")
	return float32(f)
}

func Float64(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	moon.Check(err, "parse float64")
	return f
}
