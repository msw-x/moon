package db

func Count[T any](o *Db) (int, error) {
	return o.Count((*T)(nil))
}
