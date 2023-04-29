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
