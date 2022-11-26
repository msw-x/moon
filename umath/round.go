package umath

import "math"

func StrictFrac[T AnyFloat](x T, fracLen int, round func(float64) float64) T {
	accuracy := 1 / math.Pow10(fracLen)
	return ReduceFrac(x, accuracy, round)
}

func StrictFracRound[T AnyFloat](x T, fracLen int) T {
	return StrictFrac(x, fracLen, math.Round)
}

func StrictFracTrunc[T AnyFloat](x T, fracLen int) T {
	return StrictFrac(x, fracLen, math.Trunc)
}

func StrictFracFloor[T AnyFloat](x T, fracLen int) T {
	return StrictFrac(x, fracLen, math.Floor)
}

func StrictFracCeil[T AnyFloat](x T, fracLen int) T {
	return StrictFrac(x, fracLen, math.Ceil)
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
