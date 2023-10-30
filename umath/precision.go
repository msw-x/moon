package umath

import "math"

func Precision(f float64) int {
	if f == 0 {
		return 0
	}
	return int(math.Round(math.Log10(math.Abs(f))))
}
