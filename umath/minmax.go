package umath

import "math"

func Minf[A Number, B Number](a A, b B) float64 {
	return math.Min(float64(a), float64(b))
}

func Maxf[A Number, B Number](a A, b B) float64 {
	return math.Max(float64(a), float64(b))
}

func Min[T Number](a, b T) T {
	return T(Minf(a, b))
}

func Max[T Number](a, b T) T {
	return T(Maxf(a, b))
}
