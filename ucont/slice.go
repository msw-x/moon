package ucont

import (
	"sort"

	"github.com/msw-x/moon/umath"
)

type Slice[T any] []T

func NewSlice[T any]() Slice[T] {
	return Slice[T]{}
}

func NewSliceWithSize[T any](size int) Slice[T] {
	return Slice[T](make([]T, size))
}

func NewSliceWithCapacity[T any](capacity int) Slice[T] {
	return Slice[T](make([]T, 0, capacity))
}

func (this Slice[T]) Get(index int) T {
	return this[index]
}

func (this Slice[T]) Set(index int, v T) {
	this[index] = v
}

func (this Slice[T]) Data() []T {
	return this
}

func (this Slice[T]) Empty() bool {
	return this.Size() == 0
}

func (this Slice[T]) Size() int {
	return len(this)
}

func (this Slice[T]) Capacity() int {
	return cap(this)
}

func (this Slice[T]) Front() T {
	return this[0]
}

func (this Slice[T]) Back() T {
	return this[this.Size()-1]
}

func (this Slice[T]) Head(count int) Slice[T] {
	return this.FromTo(0, count)
}

func (this Slice[T]) Tail(count int) Slice[T] {
	return this.FromTo(this.Size()-count, this.Size())
}

func (this Slice[T]) HeadMax(count int) Slice[T] {
	count = int(umath.Min(count, this.Size()))
	return this.Head(count)
}

func (this Slice[T]) TailMax(count int) Slice[T] {
	count = int(umath.Min(count, this.Size()))
	return this.Tail(count)
}

func (this Slice[T]) FromTo(from, to int) Slice[T] {
	return this[from:to]
}

func (this Slice[T]) Equal(o Slice[T], fn func(T, T) bool) bool {
	if this.Size() != o.Size() {
		return false
	}
	for n, v := range this {
		if !fn(o[n], v) {
			return false
		}
	}
	return true
}

func (this Slice[T]) Sort(fn func(T, T) bool) {
	sort.Slice(this, func(i, j int) bool {
		return fn(this[i], this[j])
	})
}

func (this Slice[T]) Reverse() {
	Reverse(this)
}

func (this Slice[T]) Tansform(fn func(v T) T) {
	for n, v := range this {
		this[n] = fn(v)
	}
}

func (this *Slice[T]) SetData(data []T) {
	*this = data
}

func (this *Slice[T]) Clear() {
	this.Resize(0)
}

func (this *Slice[T]) Resize(size int) {
	if size > this.Size() {
		o := NewSliceWithSize[T](size)
		o.CopyFrom(*this)
		*this = o
	} else if size < this.Size() {
		*this = this.Head(size)
	}
}

func (this *Slice[T]) CopyFrom(o Slice[T]) {
	count := int(umath.Min(this.Size(), o.Size()))
	for n, v := range o[0:count] {
		(*this)[n] = v
	}
}

func (this *Slice[T]) Insert(index int, v T) {
	*this = Insert(*this, index, v)
}

func (this *Slice[T]) PushBack(v T) {
	*this = append(*this, v)
}

func (this *Slice[T]) PushFront(v T) {
	this.Insert(0, v)
}

func (this *Slice[T]) Erase(index int) {
	*this = Remove(*this, index)
}

func (this *Slice[T]) EraseIf(fn func(T) bool) {
	o := NewSlice[T]()
	for _, v := range *this {
		if !fn(v) {
			o.PushBack(v)
		}
	}
	*this = o
}
