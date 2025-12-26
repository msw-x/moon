package collection2

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"sync"

	"github.com/msw-x/moon/app"
	"github.com/msw-x/moon/db"
	"github.com/msw-x/moon/uerr"
	"github.com/msw-x/moon/ulog"
	"github.com/uptrace/bun"
	"golang.org/x/exp/constraints"
)

type Sync[Id constraints.Ordered, MapItem any, DbItem any] struct {
	log          *ulog.Log
	db           *db.Db
	m            map[Id]MapItem
	name         string
	fn           Funcs[Id, MapItem, DbItem]
	onSelect     func(*bun.SelectQuery)
	onDelete     func(MapItem, *bun.DeleteQuery)
	mutex        sync.RWMutex
	excludeMutex bool
	inited       bool
	logUpdate    bool
	dbro         bool
}

func (o *Sync[Id, MapItem, DbItem]) Open(name string, log *ulog.Log, db *db.Db, dbItemId func(DbItem) Id) {
	if log == nil {
		log = ulog.New(name)
	} else {
		log = log.Branch(name)
	}
	o.log = log.WithLifetime()
	o.db = db
	o.name = name
	o.fn.dbItemId = dbItemId
	o.assertFn("db-item-id", o.fn.dbItemId)
}

func (o *Sync[Id, MapItem, DbItem]) Close() {
	o.log.Close()
}

func (o *Sync[Id, MapItem, DbItem]) OnConvert(newMapItem func(DbItem) MapItem, mapToDbItem func(MapItem) DbItem) {
	o.fn.newMapItem = newMapItem
	o.fn.mapToDbItem = mapToDbItem
}

func (o *Sync[Id, MapItem, DbItem]) OnSelect(onSelect func(*bun.SelectQuery)) {
	o.onSelect = onSelect
}

func (o *Sync[Id, MapItem, DbItem]) OnDelete(onDelete func(MapItem, *bun.DeleteQuery)) {
	o.onDelete = onDelete
}

func (o *Sync[Id, MapItem, DbItem]) ExcludeMutex() {
	o.excludeMutex = true
}

func (o *Sync[Id, MapItem, DbItem]) DbReadonly() {
	o.dbro = true
}

func (o *Sync[Id, MapItem, DbItem]) Db() *db.Db {
	return o.db
}

func (o *Sync[Id, MapItem, DbItem]) Log() *ulog.Log {
	return o.log
}

func (o *Sync[Id, MapItem, DbItem]) LogUpdate(yes bool) {
	o.logUpdate = yes
}

func (o *Sync[Id, MapItem, DbItem]) Inited() bool {
	return o.inited
}

func (o *Sync[Id, MapItem, DbItem]) InitError() error {
	if !o.Inited() {
		return fmt.Errorf("%s not inited", o.name)
	}
	return nil
}

func (o *Sync[Id, MapItem, DbItem]) Count() int {
	if o.Inited() {
		return len(o.m)
	}
	return 0
}

func (o *Sync[Id, MapItem, DbItem]) Empty() bool {
	return o.Count() == 0
}

func (o *Sync[Id, MapItem, DbItem]) Init() bool {
	if !o.db.Ok() {
		return false
	}
	var list []DbItem
	if o.db.Select(&list, o.onSelect) != nil {
		return false
	}
	if !o.excludeMutex {
		o.mutex.Lock()
		defer o.mutex.Unlock()
	}
	o.m = make(map[Id]MapItem)
	for _, e := range list {
		o.put(e)
	}
	o.inited = true
	o.log.Info("inited. count:", o.Count())
	return true
}

func (o *Sync[Id, MapItem, DbItem]) Add(e DbItem) (Id, error) {
	return o.add("add", e)
}

func (o *Sync[Id, MapItem, DbItem]) AddNamed(name string, e DbItem) (Id, error) {
	return o.add(fmt.Sprintf("add[%s]", name), e)
}

func (o *Sync[Id, MapItem, DbItem]) Attach(e DbItem) {
	o.ins("attach", e)
}

func (o *Sync[Id, MapItem, DbItem]) AttachNamed(name string, e DbItem) {
	o.ins(fmt.Sprintf("attach[%s]", name), e)
}

func (o *Sync[Id, MapItem, DbItem]) Delete(id Id) error {
	o.log.Debugf("delete[%v]", id)
	e, err := o.Get(id)
	if err == nil {
		if !o.excludeMutex {
			o.mutex.Lock()
			defer o.mutex.Unlock()
		}
		v := o.fn.mapToDbItem(e)
		var onDelete func(*bun.DeleteQuery)
		if o.onDelete != nil {
			onDelete = func(q *bun.DeleteQuery) {
				o.onDelete(e, q)
			}
		}
		if !o.dbro {
			_, err = o.db.Delete(&v, onDelete)
		}
		if err == nil {
			delete(o.m, id)
			o.log.Infof("delete[%v] completed", id)
		}
	}
	return err
}

func (o *Sync[Id, MapItem, DbItem]) DeleteAll() (err error) {
	o.log.Debug("delete all")
	if o.onDelete == nil {
		if !o.excludeMutex {
			o.mutex.Lock()
			defer o.mutex.Unlock()
		}
		if !o.dbro {
			_, err = db.DeleteAll[DbItem](o.db)
		}
		if err == nil {
			o.m = make(map[Id]MapItem)
		}
	} else {
		err = errors.New("collection: delete all not implemented if onDelete defined")
	}
	return
}

func (o *Sync[Id, MapItem, DbItem]) SoftDelete(id Id, fn func(e MapItem) MapItem) error {
	return o.update(updateDelete, id, fn)
}

func (o *Sync[Id, MapItem, DbItem]) Remove(id Id, fn func(e MapItem) MapItem) error {
	return o.update(updateRemove, id, fn)
}

func (o *Sync[Id, MapItem, DbItem]) Update(id Id, fn func(e MapItem) MapItem) error {
	return o.update(updatePure, id, fn)
}

func (o *Sync[Id, MapItem, DbItem]) SoftUpdate(id Id, fn func(e MapItem) MapItem) error {
	return o.update(updateSoft, id, fn)
}

func (o *Sync[Id, MapItem, DbItem]) Replace(e MapItem) error {
	id := o.mapItemId(e)
	return o.Update(id, func(MapItem) MapItem {
		return e
	})
}

func (o *Sync[Id, MapItem, DbItem]) SoftReplace(e MapItem) error {
	id := o.mapItemId(e)
	return o.SoftUpdate(id, func(MapItem) MapItem {
		return e
	})
}

func (o *Sync[Id, MapItem, DbItem]) ForEach(fn func(MapItem)) {
	if !o.excludeMutex {
		o.mutex.RLock()
		defer o.mutex.RUnlock()
	}
	o.check()
	for _, e := range o.m {
		fn(e)
	}
}

func (o *Sync[Id, MapItem, DbItem]) ForEachSwarm(fn func(MapItem)) {
	if !o.excludeMutex {
		o.mutex.RLock()
		defer o.mutex.RUnlock()
	}
	o.check()
	swarm := app.NewGoSwarm().WithLog(o.log)
	for _, e := range o.m {
		v := e
		swarm.Add(func() {
			fn(v)
		})
	}
	swarm.Wait()
}

func (o *Sync[Id, MapItem, DbItem]) ForEachSwarmPool(fn func(MapItem), limit int) {
	if !o.excludeMutex {
		o.mutex.RLock()
		defer o.mutex.RUnlock()
	}
	o.check()
	swarm := app.NewGoSwarmLimit().WithLog(o.log)
	for _, e := range o.m {
		v := e
		swarm.Add(func() {
			fn(v)
		})
	}
	swarm.Execute(limit)
}

func (o *Sync[Id, MapItem, DbItem]) Walk(fn func(MapItem) bool) bool {
	if !o.excludeMutex {
		o.mutex.RLock()
		defer o.mutex.RUnlock()
	}
	o.check()
	for _, e := range o.m {
		if fn(e) {
			return true
		}
	}
	return false
}

func (o *Sync[Id, MapItem, DbItem]) Keys() []Id {
	if !o.excludeMutex {
		o.mutex.RLock()
		defer o.mutex.RUnlock()
	}
	var l []Id
	for id := range o.m {
		l = append(l, id)
	}
	return l
}

func (o *Sync[Id, MapItem, DbItem]) List() []MapItem {
	if !o.excludeMutex {
		o.mutex.RLock()
		defer o.mutex.RUnlock()
	}
	o.check()
	var l []MapItem
	for _, e := range o.m {
		l = append(l, e)
	}
	sort.Slice(l, func(i, j int) bool {
		return o.mapItemId(l[i]) < o.mapItemId(l[j])
	})
	return l
}

func (o *Sync[Id, MapItem, DbItem]) DbList() []DbItem {
	lst := o.List()
	l := make([]DbItem, len(lst))
	for n, e := range lst {
		l[n] = o.fn.mapToDbItem(e)
	}
	return l
}

func (o *Sync[Id, MapItem, DbItem]) Exist(id Id) bool {
	if !o.excludeMutex {
		o.mutex.RLock()
		defer o.mutex.RUnlock()
	}
	o.check()
	_, ok := o.m[id]
	return ok
}

func (o *Sync[Id, MapItem, DbItem]) Get(id Id) (MapItem, error) {
	if !o.excludeMutex {
		o.mutex.RLock()
		defer o.mutex.RUnlock()
	}
	o.check()
	i, ok := o.m[id]
	var err error
	if !ok {
		err = fmt.Errorf("%s get[%v]: not found", o.name, id)
	}
	return i, err
}

func (o *Sync[Id, MapItem, DbItem]) GetIfExists(id Id) (MapItem, bool) {
	if !o.excludeMutex {
		o.mutex.RLock()
		defer o.mutex.RUnlock()
	}
	o.check()
	i, ok := o.m[id]
	return i, ok
}

func (o *Sync[Id, MapItem, DbItem]) add(action string, e DbItem) (Id, error) {
	o.log.Debug(action)
	err := o.InitError()
	if err == nil {
		if !o.dbro {
			err = o.db.Insert(&e)
		}
	}
	var id Id
	if err == nil {
		id = o.ins(action, e)
	} else {
		o.log.Error(action, err)
	}
	return id, err
}

func (o *Sync[Id, MapItem, DbItem]) ins(action string, e DbItem) (id Id) {
	id = o.fn.dbItemId(e)
	o.log.Info(action, "id:", id)
	if !o.excludeMutex {
		o.mutex.Lock()
		defer o.mutex.Unlock()
	}
	o.check()
	o.put(e)
	return
}

func (o *Sync[Id, MapItem, DbItem]) put(e DbItem) {
	o.assertFn("new-map-item", o.fn.newMapItem)
	o.m[o.fn.dbItemId(e)] = o.fn.newMapItem(e)
}

func (o *Sync[Id, MapItem, DbItem]) update(mode updateMode, id Id, fn func(e MapItem) MapItem) error {
	var action string
	switch mode {
	case updatePure:
		action = "update"
	case updateSoft:
		action = "update (soft)"
	case updateRemove:
		action = "remove"
	case updateDelete:
		action = "delete"
	}
	action = fmt.Sprintf("%s[%v]", action, id)
	if o.logUpdate {
		o.log.Debug(action)
	}
	err := o.InitError()
	if err == nil {
		var e MapItem
		e, err = o.Get(id)
		if err == nil {
			o.assertFn("map-to-db-item", o.fn.mapToDbItem)
			prev := o.fn.mapToDbItem(e)
			e = fn(e)
			if mode != updateSoft {
				if !o.dbro {
					curr := o.fn.mapToDbItem(e)
					_, err = db.UpdateDiff(o.db, prev, curr)
				}
			}
			if err == nil {
				if !o.excludeMutex {
					o.mutex.Lock()
					defer o.mutex.Unlock()
				}
				o.check()
				if mode == updateDelete {
					delete(o.m, id)
				} else {
					o.m[id] = e
				}
				if o.logUpdate {
					o.log.Info(action, "completed")
				}
			}
		}
	}
	if err != nil {
		o.log.Error(action, err)
	}
	return err
}

func (o *Sync[Id, MapItem, DbItem]) mapItemId(e MapItem) Id {
	return o.fn.dbItemId(o.fn.mapToDbItem(e))
}

func (o *Sync[Id, MapItem, DbItem]) assertFn(name string, fn any) {
	if reflect.ValueOf(fn).IsNil() {
		uerr.Panicf("%s func '%s' is nil", o.name, name)
	}
}

func (o *Sync[Id, MapItem, DbItem]) check() {
	if err := o.InitError(); err != nil {
		o.log.Error(err)
	}
}
