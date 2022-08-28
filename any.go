package moon

type AnyInt interface {
	int | int32 | uint32 | int64 | uint64
}

type AnyNumber interface {
	AnyInt | float32 | float64
}
