package umath

import "math"

func RestrictFrac[T AnyFloat](x T, fracLen int, round func(float64) float64) T {
	accuracy := 1 / math.Pow10(fracLen)
	return ReduceFrac(x, accuracy, round)
}

func RestrictFracRound[T AnyFloat](x T, fracLen int) T {
	return RestrictFrac(x, fracLen, math.Round)
}

func RestrictFracTrunc[T AnyFloat](x T, fracLen int) T {
	return RestrictFrac(x, fracLen, math.Trunc)
}

func RestrictFracFloor[T AnyFloat](x T, fracLen int) T {
	return RestrictFrac(x, fracLen, math.Floor)
}

func RestrictFracCeil[T AnyFloat](x T, fracLen int) T {
	return RestrictFrac(x, fracLen, math.Ceil)
}

func ReduceFrac[T, A AnyFloat](x T, accuracy A, round func(float64) float64) T {
	return T(round(float64(x/T(accuracy))) * float64(accuracy))
}

func ReduceFracRound[T, A AnyFloat](x T, accuracy A) T {
	return ReduceFrac(x, accuracy, math.Round)
}

func ReduceFracTrunc[T, A AnyFloat](x T, accuracy A) T {
	return ReduceFrac(x, accuracy, math.Trunc)
}

func ReduceFracFloor[T, A AnyFloat](x T, accuracy A) T {
	return ReduceFrac(x, accuracy, math.Floor)
}

func ReduceFracCeil[T, A AnyFloat](x T, accuracy A) T {
	return ReduceFrac(x, accuracy, math.Ceil)
}
