package parse

import (
	"strconv"
)

func Int(s string) (int, error) {
	return strconv.Atoi(s)
}

func Int64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

func Uint64(s string) (uint64, error) {
	return strconv.ParseUint(s, 10, 64)
}

func Float32(s string) (float32, error) {
	f, err := strconv.ParseFloat(s, 32)
	return float32(f), err
}

func Float64(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}
