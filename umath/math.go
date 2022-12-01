package umath

import (
	"math"
	"math/rand"

	"golang.org/x/exp/constraints"
)

func Percent[Minor Number, Major Number](minor Minor, major Major) int {
	mn := int64(minor)
	mj := int64(major)
	if mn == 0 {
		return 0
	}
	return int(float64(mn) / float64(mj) * 100)
}

func Rand(min, max int) int {
	return rand.Intn(max-min) + min
}

func NormalFloatDegree[T constraints.Float](v T) int {
	return int(math.Floor(math.Log10(float64(v))))
}

func Min[A Number, B Number](a A, b B) float64 {
	return math.Min(float64(a), float64(b))
}

func Max[A Number, B Number](a A, b B) float64 {
	return math.Max(float64(a), float64(b))
}
