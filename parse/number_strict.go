package parse

import "github.com/msw-x/moon"

func IntStrict(s string) int {
	i, err := Int(s)
	moon.Strict(err, "parse int")
	return i
}

func Int64Strict(s string) int64 {
	i, err := Int64(s)
	moon.Strict(err, "parse int64")
	return i
}

func Uint64Strict(s string) uint64 {
	i, err := Uint64(s)
	moon.Strict(err, "parse uint64")
	return i
}

func Float32Strict(s string) float32 {
	i, err := Float32(s)
	moon.Strict(err, "parse float32")
	return i
}

func Float64Strict(s string) float64 {
	i, err := Float64(s)
	moon.Strict(err, "parse float64")
	return i
}
