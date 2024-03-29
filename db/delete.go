package db

import "github.com/uptrace/bun"

func Delete[T any](o *Db, fn func(*bun.DeleteQuery)) (int64, error) {
	return o.Delete((*T)(nil), fn)
}

func DeleteAll[T any](o *Db) (int64, error) {
	return o.DeleteAll((*T)(nil))
}
