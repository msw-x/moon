package ucont

func Insert[S ~[]T, T any](s S, index int, v T) S {
	s = append(s[:index+1], s[index:]...)
	s[index] = v
	return s
}

func Remove[S ~[]T, T any](s S, index int) S {
	return append(s[:index], s[index+1:]...)
}

func Reverse[S ~[]T, T any](s S) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func Equal[S ~[]T, T comparable](a, b S) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func Unique[S ~[]T, T comparable](s S) (r S) {
	keys := make(map[T]bool)
	for _, v := range s {
		if !keys[v] {
			keys[v] = true
			r = append(r, v)
		}
	}
	return
}

func Filter[T any](s []T, f func(T) bool) (r []T) {
	for _, v := range s {
		if f(v) {
			r = append(r, v)
		}
	}
	return
}
