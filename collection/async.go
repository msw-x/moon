package collection

import (
	"time"

	"github.com/msw-x/moon/app"
	"github.com/msw-x/moon/db"
	"github.com/msw-x/moon/ulog"
	"github.com/uptrace/bun"
	"golang.org/x/exp/constraints"
)

type Async[Id constraints.Ordered, MapItem any, DbItem any] struct {
	c               Sync[Id, MapItem, DbItem]
	job             *app.Job
	handle          func()
	handleImmediate bool
}

func (o *Async[Id, MapItem, DbItem]) Open(name string, log *ulog.Log, db *db.Db, dbItemId func(DbItem) Id, interval time.Duration, handle func()) {
	o.c.Open(name, log, db, dbItemId)
	o.handle = handle
	o.job = app.NewJob().WithLog(o.c.Log())
	o.job.Tick(o.process, interval)
}

func (o *Async[Id, MapItem, DbItem]) Close() {
	defer o.c.Close()
	defer o.job.Stop()
}

func (o *Async[Id, MapItem, DbItem]) TwoStageClose() func() {
	o.job.Stop()
	return func() {
		o.c.Close()
	}
}

func (o *Async[Id, MapItem, DbItem]) OnConvert(newMapItem func(DbItem) MapItem, mapToDbItem func(MapItem) DbItem) {
	o.c.OnConvert(newMapItem, mapToDbItem)
}

func (o *Async[Id, MapItem, DbItem]) OnSelect(onSelect func(*bun.SelectQuery)) {
	o.c.OnSelect(onSelect)
}

func (o *Async[Id, MapItem, DbItem]) OnDelete(onDelete func(MapItem, *bun.DeleteQuery)) {
	o.c.OnDelete(onDelete)
}

func (o *Async[Id, MapItem, DbItem]) ExcludeMutex() {
	o.c.ExcludeMutex()
}

func (o *Async[Id, MapItem, DbItem]) HandleImmediate() {
	o.handleImmediate = true
}

func (o *Async[Id, MapItem, DbItem]) Db() *db.Db {
	return o.c.Db()
}

func (o *Async[Id, MapItem, DbItem]) Log() *ulog.Log {
	return o.c.Log()
}

func (o *Async[Id, MapItem, DbItem]) LogUpdate(yes bool) {
	o.c.LogUpdate(yes)
}

func (o *Async[Id, MapItem, DbItem]) Inited() bool {
	return o.c.Inited()
}

func (o *Async[Id, MapItem, DbItem]) InitError() error {
	return o.c.InitError()
}

func (o *Async[Id, MapItem, DbItem]) Count() int {
	return o.c.Count()
}

func (o *Async[Id, MapItem, DbItem]) Empty() bool {
	return o.c.Empty()
}

func (o *Async[Id, MapItem, DbItem]) Init() {
	o.c.Init()
}

func (o *Async[Id, MapItem, DbItem]) Add(e DbItem) (Id, error) {
	return o.c.Add(e)
}

func (o *Async[Id, MapItem, DbItem]) AddNamed(name string, e DbItem) (Id, error) {
	return o.c.AddNamed(name, e)
}

func (o *Async[Id, MapItem, DbItem]) Delete(id Id) error {
	return o.c.Delete(id)
}

func (o *Async[Id, MapItem, DbItem]) DeleteAll() error {
	return o.c.DeleteAll()
}

func (o *Async[Id, MapItem, DbItem]) Remove(id Id, fn func(e MapItem) MapItem) error {
	return o.c.Remove(id, fn)
}

func (o *Async[Id, MapItem, DbItem]) Update(id Id, fn func(e MapItem) MapItem) error {
	return o.c.Update(id, fn)
}

func (o *Async[Id, MapItem, DbItem]) Replace(e MapItem) error {
	return o.c.Replace(e)
}

func (o *Async[Id, MapItem, DbItem]) ForEach(fn func(MapItem)) {
	o.c.ForEach(fn)
}

func (o *Async[Id, MapItem, DbItem]) ForEachSwarm(fn func(MapItem)) {
	o.c.ForEachSwarm(fn)
}

func (o *Async[Id, MapItem, DbItem]) Walk(fn func(MapItem) bool) bool {
	return o.c.Walk(fn)
}

func (o *Async[Id, MapItem, DbItem]) Keys() []Id {
	return o.c.Keys()
}

func (o *Async[Id, MapItem, DbItem]) List() []MapItem {
	return o.c.List()
}

func (o *Async[Id, MapItem, DbItem]) DbList() []DbItem {
	return o.c.DbList()
}

func (o *Async[Id, MapItem, DbItem]) Exist(id Id) bool {
	return o.c.Exist(id)
}

func (o *Async[Id, MapItem, DbItem]) Get(id Id) (MapItem, error) {
	return o.c.Get(id)
}

func (o *Async[Id, MapItem, DbItem]) GetIfExists(id Id) (MapItem, bool) {
	return o.c.GetIfExists(id)
}

func (o *Async[Id, MapItem, DbItem]) process() {
	if o.Inited() {
		o.processHandle()
	} else {
		o.c.Init()
		if o.handleImmediate && o.Inited() {
			o.processHandle()
		}
	}
}

func (o *Async[Id, MapItem, DbItem]) processHandle() {
	if o.handle == nil {
		o.job.Cancel()
	} else {
		o.handle()
	}
}
