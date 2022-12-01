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

func (this Numerics[T]) Get(index int) T {
	return this.ref().Get(index)
}

func (this Numerics[T]) Set(index int, v T) {
	this.ref().Set(index, v)
}

func (this Numerics[T]) Data() []T {
	return this.ref().Data()
}

func (this Numerics[T]) Empty() bool {
	return this.ref().Empty()
}

func (this Numerics[T]) Size() int {
	return this.ref().Size()
}

func (this Numerics[T]) Capacity() int {
	return this.ref().Capacity()
}

func (this Numerics[T]) Front() T {
	return this.ref().Front()
}

func (this Numerics[T]) Back() T {
	return this.ref().Back()
}

func (this Numerics[T]) Head(count int) Numerics[T] {
	return this.FromTo(0, count)
}

func (this Numerics[T]) Tail(count int) Numerics[T] {
	return this.FromTo(this.Size()-count, this.Size())
}

func (this Numerics[T]) HeadMax(count int) Numerics[T] {
	count = int(umath.Min(count, this.Size()))
	return this.Head(count)
}

func (this Numerics[T]) TailMax(count int) Numerics[T] {
	count = int(umath.Min(count, this.Size()))
	return this.Tail(count)
}

func (this Numerics[T]) FromTo(from, to int) Numerics[T] {
	return this[from:to]
}

func (this Numerics[T]) Equal(o Numerics[T]) bool {
	return this.ref().Equal(o.ref(), func(a, b T) bool {
		return a == b
	})
}

func (this Numerics[T]) Sort() {
	this.ref().Sort(func(a, b T) bool {
		return a < b
	})
}

func (this Numerics[T]) Find(w T) int {
	for n, v := range this.ref() {
		if v == w {
			return n
		}
	}
	return -1
}

func (this Numerics[T]) Includes(v T) bool {
	return this.Find(v) != -1
}

func (this Numerics[T]) Reverse() {
	this.ref().Reverse()
}

func (this Numerics[T]) Transform(fn func(v T) T) {
	this.ref().Transform(fn)
}

func (this *Numerics[T]) SetData(data []T) {
	this.ptr().SetData(data)
}

func (this *Numerics[T]) Clear() {
	this.ptr().Clear()
}

func (this *Numerics[T]) Resize(size int) {
	this.ptr().Resize(size)
}

func (this *Numerics[T]) CopyFrom(o Numerics[T]) {
	this.ptr().CopyFrom(o.ref())
}

func (this *Numerics[T]) Insert(index int, v T) {
	this.ptr().Insert(index, v)
}

func (this *Numerics[T]) PushBack(v T) {
	this.ptr().PushBack(v)
}

func (this *Numerics[T]) PushFront(v T) {
	this.ptr().PushFront(v)
}

func (this *Numerics[T]) Erase(index int) {
	this.ptr().Erase(index)
}

func (this *Numerics[T]) EraseIf(fn func(T) bool) {
	this.ptr().EraseIf(fn)
}

func (this *Numerics[T]) EraseAll(w T) {
	this.ptr().EraseIf(func(v T) bool {
		return v == w
	})
}

func (this Numerics[T]) ref() Slice[T] {
	return Slice[T](this)
}

func (this *Numerics[T]) ptr() *Slice[T] {
	return (*Slice[T])(this)
}
