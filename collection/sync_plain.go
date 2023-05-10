package collection

import (
	"github.com/msw-x/moon/db"
	"github.com/msw-x/moon/ulog"
	"github.com/uptrace/bun"
	"golang.org/x/exp/constraints"
)

type SyncPlain[Id constraints.Ordered, Item any] struct {
	c Sync[Id, Item, Item]
}

func (o *SyncPlain[Id, Item]) Open(name string, log *ulog.Log, db *db.Db, ItemId func(Item) Id) {
	o.c.Open(name, log, db, ItemId)
	o.c.OnConvert(plainConvert[Item], plainConvert[Item])
}

func (o *SyncPlain[Id, Item]) Close() {
	o.c.Close()
}

func (o *SyncPlain[Id, Item]) OnSelect(onSelect func(*bun.SelectQuery)) {
	o.c.OnSelect(onSelect)
}

func (o *SyncPlain[Id, Item]) Log() *ulog.Log {
	return o.c.Log()
}

func (o *SyncPlain[Id, Item]) LogUpdate(yes bool) {
	o.c.LogUpdate(yes)
}

func (o *SyncPlain[Id, Item]) Inited() bool {
	return o.c.Inited()
}

func (o *SyncPlain[Id, Item]) InitError() error {
	return o.c.InitError()
}

func (o *SyncPlain[Id, Item]) Count() int {
	return o.c.Count()
}

func (o *SyncPlain[Id, Item]) Empty() bool {
	return o.c.Empty()
}

func (o *SyncPlain[Id, Item]) Init() bool {
	return o.c.Init()
}

func (o *SyncPlain[Id, Item]) Add(e Item) (Id, error) {
	return o.c.Add(e)
}

func (o *SyncPlain[Id, Item]) AddNamed(name string, e Item) (Id, error) {
	return o.c.AddNamed(name, e)
}

func (o *SyncPlain[Id, Item]) Delete(id Id) error {
	return o.c.Delete(id)
}

func (o *SyncPlain[Id, Item]) Remove(id Id, fn func(e Item) Item) error {
	return o.c.Remove(id, fn)
}

func (o *SyncPlain[Id, Item]) Update(id Id, fn func(e Item) Item) error {
	return o.c.Update(id, fn)
}

func (o *SyncPlain[Id, Item]) Replace(e Item) error {
	return o.c.Replace(e)
}

func (o *SyncPlain[Id, Item]) Upsert(e Item) error {
	id := o.c.fn.dbItemId(e)
	if o.Exist(id) {
		return o.Replace(e)
	}
	_, err := o.Add(e)
	return err
}

func (o *SyncPlain[Id, Item]) UpsertNamed(name string, e Item) error {
	id := o.c.fn.dbItemId(e)
	if o.Exist(id) {
		return o.Replace(e)
	}
	_, err := o.AddNamed(name, e)
	return err
}

func (o *SyncPlain[Id, Item]) ForEach(fn func(Item)) {
	o.c.ForEach(fn)
}

func (o *SyncPlain[Id, Item]) Walk(fn func(Item) bool) bool {
	return o.c.Walk(fn)
}

func (o *SyncPlain[Id, Item]) List() []Item {
	return o.c.List()
}

func (o *SyncPlain[Id, Item]) Exist(id Id) bool {
	return o.c.Exist(id)
}

func (o *SyncPlain[Id, Item]) Get(id Id) (Item, error) {
	return o.c.Get(id)
}
