package ptr

func Of[T any](v T) *T {
	return &v
}

func To[T any](p *T) T {
	var v T
	return ToDef(p, v)
}

func ToDef[T any](p *T, v T) T {
	if p != nil {
		v = *p
	}
	return v
}

func Equal[T comparable](a, b *T) bool {
	if a != nil && b != nil {
		return *a == *b
	}
	return a == b
}
