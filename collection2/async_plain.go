package collection2

import (
	"time"

	"github.com/msw-x/moon/db"
	"github.com/msw-x/moon/ulog"
	"github.com/uptrace/bun"
	"golang.org/x/exp/constraints"
)

type AsyncPlain[Id constraints.Ordered, Item any] struct {
	c Async[Id, Item, Item]
}

func (o *AsyncPlain[Id, Item]) Open(name string, log *ulog.Log, db *db.Db, ItemId func(Item) Id, interval time.Duration, handle func()) {
	o.c.OnConvert(plainConvert[Item], plainConvert[Item])
	o.c.Open(name, log, db, ItemId, interval, handle)
}

func (o *AsyncPlain[Id, Item]) Close() {
	o.c.Close()
}

func (o *AsyncPlain[Id, Item]) TwoStageClose() func() {
	return o.c.TwoStageClose()
}

func (o *AsyncPlain[Id, Item]) OnSelect(onSelect func(*bun.SelectQuery)) {
	o.c.OnSelect(onSelect)
}

func (o *AsyncPlain[Id, Item]) OnDelete(onDelete func(Item, *bun.DeleteQuery)) {
	o.c.OnDelete(onDelete)
}

func (o *AsyncPlain[Id, Item]) ExcludeMutex() {
	o.c.ExcludeMutex()
}

func (o *AsyncPlain[Id, Item]) HandleImmediate() {
	o.c.HandleImmediate()
}

func (o *AsyncPlain[Id, Item]) Db() *db.Db {
	return o.c.Db()
}

func (o *AsyncPlain[Id, Item]) Log() *ulog.Log {
	return o.c.Log()
}

func (o *AsyncPlain[Id, Item]) LogUpdate(yes bool) {
	o.c.LogUpdate(yes)
}

func (o *AsyncPlain[Id, Item]) Inited() bool {
	return o.c.Inited()
}

func (o *AsyncPlain[Id, Item]) InitError() error {
	return o.c.InitError()
}

func (o *AsyncPlain[Id, Item]) Count() int {
	return o.c.Count()
}

func (o *AsyncPlain[Id, Item]) Empty() bool {
	return o.c.Empty()
}

func (o *AsyncPlain[Id, Item]) Init() {
	o.c.Init()
}

func (o *AsyncPlain[Id, Item]) Add(e Item) (Id, error) {
	return o.c.Add(e)
}

func (o *AsyncPlain[Id, Item]) AddNamed(name string, e Item) (Id, error) {
	return o.c.AddNamed(name, e)
}

func (o *AsyncPlain[Id, Item]) Delete(id Id) error {
	return o.c.Delete(id)
}

func (o *AsyncPlain[Id, Item]) DeleteAll() error {
	return o.c.DeleteAll()
}

func (o *AsyncPlain[Id, Item]) SoftDelete(id Id, fn func(e Item) Item) error {
	return o.c.SoftDelete(id, fn)
}

func (o *AsyncPlain[Id, Item]) Remove(id Id, fn func(e Item) Item) error {
	return o.c.Remove(id, fn)
}

func (o *AsyncPlain[Id, Item]) Update(id Id, fn func(e Item) Item) error {
	return o.c.Update(id, fn)
}

func (o *AsyncPlain[Id, Item]) SoftUpdate(id Id, fn func(e Item) Item) error {
	return o.c.SoftUpdate(id, fn)
}

func (o *AsyncPlain[Id, Item]) Replace(e Item) error {
	return o.c.Replace(e)
}

func (o *AsyncPlain[Id, Item]) SoftReplace(e Item) error {
	return o.c.SoftReplace(e)
}

func (o *AsyncPlain[Id, Item]) ForEach(fn func(Item)) {
	o.c.ForEach(fn)
}

func (o *AsyncPlain[Id, Item]) ForEachSwarm(fn func(Item)) {
	o.c.ForEachSwarm(fn)
}

func (o *AsyncPlain[Id, Item]) ForEachSwarmPool(fn func(Item), limit int) {
	o.c.ForEachSwarmPool(fn, limit)
}

func (o *AsyncPlain[Id, Item]) Walk(fn func(Item) bool) bool {
	return o.c.Walk(fn)
}

func (o *AsyncPlain[Id, Item]) Keys() []Id {
	return o.c.Keys()
}

func (o *AsyncPlain[Id, Item]) List() []Item {
	return o.c.List()
}

func (o *AsyncPlain[Id, Item]) Exist(id Id) bool {
	return o.c.Exist(id)
}

func (o *AsyncPlain[Id, Item]) Get(id Id) (Item, error) {
	return o.c.Get(id)
}

func (o *AsyncPlain[Id, Item]) GetIfExists(id Id) (Item, bool) {
	return o.c.GetIfExists(id)
}
