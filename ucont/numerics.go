package ucont

import (
	"github.com/msw-x/moon/umath"
	"golang.org/x/exp/constraints"
)

type Numerics[T constraints.Ordered] Slice[T]

func NewNumerics[T constraints.Ordered]() Numerics[T] {
	return Numerics[T]{}
}

func NewNumericsWithSize[T constraints.Ordered](size int) Numerics[T] {
	return Numerics[T](make([]T, size))
}

func NewNumericsWithCapacity[T constraints.Ordered](capacity int) Numerics[T] {
	return Numerics[T](make([]T, 0, capacity))
}

func (o Numerics[T]) Get(index int) T {
	return o.ref().Get(index)
}

func (o Numerics[T]) Set(index int, v T) {
	o.ref().Set(index, v)
}

func (o Numerics[T]) Data() []T {
	return o.ref().Data()
}

func (o Numerics[T]) Empty() bool {
	return o.ref().Empty()
}

func (o Numerics[T]) Size() int {
	return o.ref().Size()
}

func (o Numerics[T]) Capacity() int {
	return o.ref().Capacity()
}

func (o Numerics[T]) Front() T {
	return o.ref().Front()
}

func (o Numerics[T]) Back() T {
	return o.ref().Back()
}

func (o Numerics[T]) Head(count int) Numerics[T] {
	return o.FromTo(0, count)
}

func (o Numerics[T]) Tail(count int) Numerics[T] {
	return o.FromTo(o.Size()-count, o.Size())
}

func (o Numerics[T]) HeadMax(count int) Numerics[T] {
	count = int(umath.Min(count, o.Size()))
	return o.Head(count)
}

func (o Numerics[T]) TailMax(count int) Numerics[T] {
	count = int(umath.Min(count, o.Size()))
	return o.Tail(count)
}

func (o Numerics[T]) FromTo(from, to int) Numerics[T] {
	return o[from:to]
}

func (o Numerics[T]) Equal(v Numerics[T]) bool {
	return o.ref().Equal(v.ref(), func(a, b T) bool {
		return a == b
	})
}

func (o Numerics[T]) Sort() {
	o.ref().Sort(func(a, b T) bool {
		return a < b
	})
}

func (o Numerics[T]) Find(w T) int {
	for n, v := range o.ref() {
		if v == w {
			return n
		}
	}
	return -1
}

func (o Numerics[T]) Includes(v T) bool {
	return o.Find(v) != -1
}

func (o Numerics[T]) Reverse() {
	o.ref().Reverse()
}

func (o Numerics[T]) Transform(fn func(v T) T) {
	o.ref().Transform(fn)
}

func (o *Numerics[T]) SetData(data []T) {
	o.ptr().SetData(data)
}

func (o *Numerics[T]) Clear() {
	o.ptr().Clear()
}

func (o *Numerics[T]) Resize(size int) {
	o.ptr().Resize(size)
}

func (o *Numerics[T]) CopyFrom(v Numerics[T]) {
	o.ptr().CopyFrom(v.ref())
}

func (o *Numerics[T]) Insert(index int, v T) {
	o.ptr().Insert(index, v)
}

func (o *Numerics[T]) PushBack(v T) {
	o.ptr().PushBack(v)
}

func (o *Numerics[T]) PushFront(v T) {
	o.ptr().PushFront(v)
}

func (o *Numerics[T]) Erase(index int) {
	o.ptr().Erase(index)
}

func (o *Numerics[T]) EraseIf(fn func(T) bool) {
	o.ptr().EraseIf(fn)
}

func (o *Numerics[T]) EraseAll(w T) {
	o.ptr().EraseIf(func(v T) bool {
		return v == w
	})
}

func (o Numerics[T]) ref() Slice[T] {
	return Slice[T](o)
}

func (o *Numerics[T]) ptr() *Slice[T] {
	return (*Slice[T])(o)
}
