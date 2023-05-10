package collection

import (
	"fmt"
	"ginats/db"
	"sort"
	"sync"

	"github.com/msw-x/moon"
	"github.com/msw-x/moon/ulog"
	"github.com/uptrace/bun"
	"golang.org/x/exp/constraints"
)

type Sync[Id constraints.Ordered, MapItem any, DbItem any] struct {
	log         *ulog.Log
	db          *db.Db
	m           map[Id]MapItem
	name        string
	dbItemId    func(DbItem) Id
	newMapItem  func(DbItem) MapItem
	mapToDbItem func(MapItem) DbItem
	onSelect    func(*bun.SelectQuery)
	mutex       sync.Mutex
	logUpdate   bool
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
	o.dbItemId = dbItemId
	o.assertFunc("db-item-id", o.dbItemId == nil)
}

func (o *Sync[Id, MapItem, DbItem]) Close() {
	o.log.Close()
}

func (o *Sync[Id, MapItem, DbItem]) OnConvert(newMapItem func(DbItem) MapItem, mapToDbItem func(MapItem) DbItem) {
	o.newMapItem = newMapItem
	o.mapToDbItem = mapToDbItem
}

func (o *Sync[Id, MapItem, DbItem]) OnSelect(onSelect func(*bun.SelectQuery)) {
	o.onSelect = onSelect
}

func (o *Sync[Id, MapItem, DbItem]) Log() *ulog.Log {
	return o.log
}

func (o *Sync[Id, MapItem, DbItem]) LogUpdate(yes bool) {
	o.logUpdate = yes
}

func (o *Sync[Id, MapItem, DbItem]) Inited() bool {
	return o.m != nil
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
	o.mutex.Lock()
	defer o.mutex.Unlock()
	o.m = make(map[Id]MapItem)
	for _, e := range list {
		o.put(e)
	}
	o.log.Info("inited. count:", o.Count())
	return true
}

func (o *Sync[Id, MapItem, DbItem]) Add(e DbItem) (Id, error) {
	return o.add("add", e)
}

func (o *Sync[Id, MapItem, DbItem]) AddNamed(name string, e DbItem) (Id, error) {
	return o.add(fmt.Sprintf("add[%s]", name), e)
}

func (o *Sync[Id, MapItem, DbItem]) Delete(id Id) error {
	e, err := o.Get(id)
	if err == nil {
		o.mutex.Lock()
		defer o.mutex.Unlock()
		v := o.mapToDbItem(e)
		err = o.db.Delete(&v, nil)
		if err == nil {
			delete(o.m, id)
		}
	}
	return err
}

func (o *Sync[Id, MapItem, DbItem]) Remove(id Id, fn func(e MapItem) MapItem) error {
	return o.update(true, id, fn)
}

func (o *Sync[Id, MapItem, DbItem]) Update(id Id, fn func(e MapItem) MapItem) error {
	return o.update(false, id, fn)
}

func (o *Sync[Id, MapItem, DbItem]) Replace(e MapItem) error {
	id := o.mapItemId(e)
	return o.Update(id, func(MapItem) MapItem {
		return e
	})
}

func (o *Sync[Id, MapItem, DbItem]) ForEach(fn func(MapItem)) {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	o.check()
	for _, e := range o.m {
		fn(e)
	}
}

func (o *Sync[Id, MapItem, DbItem]) Walk(fn func(MapItem) bool) bool {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	o.check()
	for _, e := range o.m {
		if fn(e) {
			return true
		}
	}
	return false
}

func (o *Sync[Id, MapItem, DbItem]) List() []MapItem {
	o.mutex.Lock()
	defer o.mutex.Unlock()
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
		l[n] = o.mapToDbItem(e)
	}
	return l
}

func (o *Sync[Id, MapItem, DbItem]) Exist(id Id) bool {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	o.check()
	_, ok := o.m[id]
	return ok
}

func (o *Sync[Id, MapItem, DbItem]) Get(id Id) (MapItem, error) {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	o.check()
	i, ok := o.m[id]
	var err error
	if !ok {
		err = fmt.Errorf("%s get[%d]: not found", o.name, id)
	}
	return i, err
}

func (o *Sync[Id, MapItem, DbItem]) add(action string, e DbItem) (Id, error) {
	o.log.Debug(action)
	err := o.InitError()
	if err == nil {
		err = o.db.Insert(&e)
	}
	var id Id
	if err == nil {
		id = o.dbItemId(e)
		o.log.Info(action, "id:", id)
		o.mutex.Lock()
		defer o.mutex.Unlock()
		o.check()
		o.put(e)
	} else {
		o.log.Error(action, err)
	}
	return id, err
}

func (o *Sync[Id, MapItem, DbItem]) put(e DbItem) {
	o.assertFunc("new-map-item", o.newMapItem == nil)
	o.m[o.dbItemId(e)] = o.newMapItem(e)
}

func (o *Sync[Id, MapItem, DbItem]) update(remove bool, id Id, fn func(e MapItem) MapItem) error {
	var action string
	if remove {
		action = "remove"
	} else {
		action = "update"
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
			o.assertFunc("map-to-db-item", o.mapToDbItem == nil)
			prev := o.mapToDbItem(e)
			e = fn(e)
			curr := o.mapToDbItem(e)
			_, err = db.UpdateDiff(o.db, prev, curr)
			if err == nil {
				o.mutex.Lock()
				defer o.mutex.Unlock()
				o.check()
				o.m[id] = e
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
	return o.dbItemId(o.mapToDbItem(e))
}

func (o *Sync[Id, MapItem, DbItem]) assertFunc(name string, isnil bool) {
	if isnil {
		moon.Panicf("%s func '%s' is nil", o.name, name)
	}
}

func (o *Sync[Id, MapItem, DbItem]) check() {
	if err := o.InitError(); err != nil {
		o.log.Error(err)
	}
}
