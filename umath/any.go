package umath

type AnyInt interface {
	int | int32 | uint32 | int64 | uint64
}

type AnyFloat interface {
	float32 | float64
}

type AnyNumber interface {
	AnyInt | AnyFloat
}
