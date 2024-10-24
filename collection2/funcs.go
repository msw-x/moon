package collection2

import "golang.org/x/exp/constraints"

type Funcs[Id constraints.Ordered, MapItem any, DbItem any] struct {
	dbItemId    func(DbItem) Id
	newMapItem  func(DbItem) MapItem
	mapToDbItem func(MapItem) DbItem
}
