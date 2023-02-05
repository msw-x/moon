package ptr

func From[T any](v T) *T {
	return &v
}

func To[T any](p *T) T {
	var v T
	return ToValueDef(p, v)
}

func ToDef[T any](p *T, v T) T {
	if p != nil {
		v = *p
	}
	return v
}

func EqualValues[T any](a, b *T) bool {
	if a == b {
		return true
	}
	return ToValue(a) == ToValue(b)
}
