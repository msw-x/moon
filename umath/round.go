package umath

import (
	"math"

	"golang.org/x/exp/constraints"
)

func StrictFrac[T constraints.Float](x T, fracLen int, round func(float64) float64) T {
	accuracy := 1 / math.Pow10(fracLen)
	return ReduceFrac(x, accuracy, round)
}

func StrictFracRound[T constraints.Float](x T, fracLen int) T {
	return StrictFrac(x, fracLen, math.Round)
}

func StrictFracTrunc[T constraints.Float](x T, fracLen int) T {
	return StrictFrac(x, fracLen, math.Trunc)
}

func StrictFracFloor[T constraints.Float](x T, fracLen int) T {
	return StrictFrac(x, fracLen, math.Floor)
}

func StrictFracCeil[T constraints.Float](x T, fracLen int) T {
	return StrictFrac(x, fracLen, math.Ceil)
}

func ReduceFrac[T, A constraints.Float](x T, accuracy A, round func(float64) float64) T {
	return T(round(float64(x/T(accuracy))) * float64(accuracy))
}

func ReduceFracRound[T, A constraints.Float](x T, accuracy A) T {
	return ReduceFrac(x, accuracy, math.Round)
}

func ReduceFracTrunc[T, A constraints.Float](x T, accuracy A) T {
	return ReduceFrac(x, accuracy, math.Trunc)
}

func ReduceFracFloor[T, A constraints.Float](x T, accuracy A) T {
	return ReduceFrac(x, accuracy, math.Floor)
}

func ReduceFracCeil[T, A constraints.Float](x T, accuracy A) T {
	return ReduceFrac(x, accuracy, math.Ceil)
}
