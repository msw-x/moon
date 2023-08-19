package db

import "github.com/uptrace/bun"

func Delete[T any](o *Db, fn func(*bun.DeleteQuery)) error {
	return o.Delete((*T)(nil), fn)
}

func DeleteAll[T any](o *Db) error {
	return o.DeleteAll((*T)(nil))
}
