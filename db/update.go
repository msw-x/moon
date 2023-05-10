package db

import (
	"github.com/serenize/snaker"
	"github.com/uptrace/bun"
)

func UpdateDiff[T any](db *Db, prev, curr T) (columns []string, err error) {
	columns = Diff(prev, curr)
	if len(columns) > 0 {
		for n, s := range columns {
			columns[n] = snaker.CamelToSnake(s)
		}
		err = db.Update(&curr, func(q *bun.UpdateQuery) {
			q.WherePK().Column(columns...)
		})
	}
	return
}

type Updater[T any] struct {
	db   *Db
	last T
}

func NewUpdater[T any](db *Db, v T) *Updater[T] {
	return &Updater[T]{
		db:   db,
		last: v,
	}
}

func (o *Updater[T]) Update(v T) (columns []string, err error) {
	columns, err = UpdateDiff(o.db, o.last, v)
	if len(columns) > 0 && err == nil {
		o.last = v
	}
	return
}

func (o *Updater[T]) Equal(v T) bool {
	return Equal(o.last, v)
}
