package db

func Equal[T any](a, b T) bool {
	return len(Diff(a, b)) == 0
}
