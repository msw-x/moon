package ptr

func From[T any](v T) *T {
	return &v
}

func Value[T any](p *T) T {
	var v T
	return ValueDef(p, v)
}

func ValueDef[T any](p *T, v T) T {
	if p != nil {
		v = *p
	}
	return v
}

func Equal[T any](a, b *T) bool {
	if a == b {
		return true
	}
	return Value(a) == Value(b)
}
