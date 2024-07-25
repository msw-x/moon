package db

import (
	"context"
	"crypto/tls"
	"database/sql"
	"fmt"
	"runtime"
	"time"

	"github.com/msw-x/moon/app"
	"github.com/msw-x/moon/ulog"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type Db struct {
	log *ulog.Log
	db  *bun.DB
	job *app.Job
	ro  bool
	ok  bool
}

func New(opt Options) *Db {
	//todo: circuit breaker
	o := new(Db)
	o.log = ulog.New("db")
	h := host(opt.Host)
	o.log.Info("host:", h)
	pgopts := []pgdriver.Option{
		pgdriver.WithNetwork("tcp"),
		pgdriver.WithAddr(h),
		pgdriver.WithInsecure(opt.Insecure),
		pgdriver.WithUser(opt.User),
		pgdriver.WithPassword(opt.Pass),
		pgdriver.WithDatabase(opt.Name),
		pgdriver.WithApplicationName(app.Name()),
		pgdriver.WithDialTimeout(opt.Timeout),
		pgdriver.WithReadTimeout(opt.Timeout),
		pgdriver.WithWriteTimeout(opt.Timeout),
	}
	if !opt.Insecure {
		pgopts = append(pgopts, pgdriver.WithTLSConfig(&tls.Config{InsecureSkipVerify: true}))
	}
	pgconn := pgdriver.NewConnector(pgopts...)
	sqldb := sql.OpenDB(pgconn)
	maxOpenConns := int(opt.MaxConnFactor * float32(runtime.GOMAXPROCS(0)))
	if maxOpenConns == 0 {
		maxOpenConns = 1
	}
	if maxOpenConns < opt.MinOpenConns {
		maxOpenConns = opt.MinOpenConns
	}
	o.log.Info("max open conns:", maxOpenConns)
	sqldb.SetMaxOpenConns(maxOpenConns)
	sqldb.SetMaxIdleConns(maxOpenConns)
	if opt.Strict {
		o.db = bun.NewDB(sqldb, pgdialect.New())
	} else {
		// make app more resilient to errors during migrations
		o.db = bun.NewDB(sqldb, pgdialect.New(), bun.WithDiscardUnknownColumns())
	}
	if opt.LogErrors || opt.LogQueries || opt.LogLongQueries {
		log := newLog(o.log, opt.LogErrors && !opt.LogQueries)
		if opt.LogLongQueries {
			if opt.LongQueriesTime == 0 {
				opt.LongQueriesTime = opt.Timeout / 2
			}
			log.WithQueriesTime(opt.LongQueriesTime, opt.WarnLongQueries)
		}
		o.db.AddQueryHook(log)
	}
	if opt.ReadOnly {
		o.log.Info("readonly")
	}
	o.ro = opt.ReadOnly
	o.job = app.NewJob().WithLog(o.log).OnFinish(o.close)
	o.job.Tick(o.checkConnection, time.Second)
	return o
}

func (o *Db) Close() {
	o.job.Stop()
}

func (o *Db) Ok() bool {
	return o.ok
}

func (o *Db) Ping() bool {
	_, err := o.db.Exec("SELECT 1")
	return err == nil
}

func (o *Db) Wait(timeout time.Duration) bool {
	tick := time.Millisecond * 10
	for o.job.Do() && timeout > 0 && !o.Ok() {
		timeout -= tick
		o.job.Sleep(tick)
	}
	return o.Ok()
}

func (o *Db) Await(timeout time.Duration) time.Duration {
	return app.Wait(o.log, o.Ok, timeout)
}

/*
func (o *Db) OnReady(f func(), timeout time.Duration) {
	if app.NewWait().WithLog(o.log).WithTimeout(timeout).WithDo(o.job.Do).Wait(o.Ok).Waited() {
		f()
	}
}
*/

func (o *Db) Migrator() *Migrator {
	return NewMigrator(o)
}

func (o *Db) Format(query string, arg ...any) string {
	return o.db.Formatter().FormatQuery(query, arg...)
}

func (o *Db) Scan(model any, query string, arg ...any) error {
	return o.db.NewRaw(query, arg...).Scan(o.ctx(), model)
}

func (o *Db) Exec(query string, arg ...any) error {
	if o.ro {
		return nil
	}
	_, err := o.db.Exec(o.Format(query, arg...))
	return err
}

func (o *Db) Count(model any) (int, error) {
	return o.NewSelect(model).Count(o.ctx())
}

func (o *Db) Insert(model any) error {
	if o.ro {
		return nil
	}
	_, err := o.db.NewInsert().Model(model).Exec(o.ctx())
	return err
}

func (o *Db) NewSelect(model any) *bun.SelectQuery {
	return o.db.NewSelect().Model(model)
}

func (o *Db) Select(model any, fn func(*bun.SelectQuery)) error {
	q := o.NewSelect(model)
	if fn != nil {
		fn(q)
	}
	return q.Scan(o.ctx())
}

func (o *Db) SelectPk(model any) error {
	q := o.NewSelect(model)
	q.WherePK()
	return q.Scan(o.ctx())
}

func (o *Db) SelectIn(model any, ids any) error {
	return o.Select(model, func(q *bun.SelectQuery) {
		q.Where("id IN (?)", bun.In(ids))
	})
}

func (o *Db) SelectAll(model any) error {
	return o.Select(model, nil)
}

func (o *Db) Update(model any, fn func(*bun.UpdateQuery)) error {
	if o.ro {
		return nil
	}
	q := o.db.NewUpdate().Model(model)
	if fn == nil {
		q.WherePK()
	} else {
		fn(q)
	}
	_, err := q.Exec(o.ctx())
	return err
}

func (o *Db) UpdatePk(model any) error {
	return o.Update(model, nil)
}

func (o *Db) UpdateIn(model any, ids any) error {
	return o.Update(model, func(q *bun.UpdateQuery) {
		q.Where("id IN (?)", bun.In(ids))
	})
}

func (o *Db) UpdateAll(model any) error {
	return o.Update(model, func(q *bun.UpdateQuery) {
		q.Where("TRUE")
	})
}

func (o *Db) Upsert(model any) error {
	if o.ro {
		return nil
	}
	pk, err := PkName(model)
	if err == nil {
		on := fmt.Sprintf("CONFLICT (%s) DO UPDATE", pk)
		_, err = o.db.NewInsert().Model(model).On(on).Exec(o.ctx())
	}
	return err
}

func (o *Db) Delete(model any, fn func(*bun.DeleteQuery)) (int64, error) {
	if o.ro {
		return 0, nil
	}
	q := o.db.NewDelete().Model(model)
	if fn == nil {
		q.WherePK()
	} else {
		fn(q)
	}
	var num int64
	res, err := q.Exec(o.ctx())
	if err == nil {
		num, _ = res.RowsAffected()
	}
	return num, err
}

func (o *Db) DeleteAll(model any) (int64, error) {
	if o.ro {
		return 0, nil
	}
	return o.Delete(model, func(q *bun.DeleteQuery) {
		q.Where("TRUE")
	})
}

func (o *Db) Truncate(model any) error {
	if o.ro {
		return nil
	}
	_, err := o.db.NewTruncateTable().Model(model).Exec(o.ctx())
	return err
}

func (o *Db) Exists(model any, fn func(*bun.SelectQuery)) (bool, error) {
	q := o.NewSelect(model)
	if fn == nil {
		q.WherePK()
	} else {
		fn(q)
	}
	return q.Exists(o.ctx())
}

func (o *Db) ExistsPk(model any) (bool, error) {
	return o.Exists(model, nil)
}

func (o *Db) GetIfExists(model any, fn func(*bun.SelectQuery)) (bool, error) {
	ok, err := o.Exists(model, fn)
	if ok {
		if fn == nil {
			err = o.SelectPk(model)
		} else {
			err = o.Select(model, fn)
		}
	}
	return ok, err
}

func (o *Db) GetIfExistsPk(model any) (bool, error) {
	return o.GetIfExists(model, nil)
}

func (o *Db) ctx() context.Context {
	return context.Background()
}

func (o *Db) checkConnection() {
	ok := o.Ping()
	if ok != o.ok {
		o.ok = ok
		if o.ok {
			o.log.Debug("connected")
		} else {
			o.log.Error("disconnected")
		}
	}
}

func (o *Db) close() {
	err := o.db.Close()
	s := o.db.DBStats()
	o.log.Infof("queries[%d] errors[%d]", s.Queries, s.Errors)
	if err != nil {
		o.log.Error("close:", err)
	}
}
