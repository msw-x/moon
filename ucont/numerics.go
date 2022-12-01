package ucont

import "golang.org/x/exp/constraints"

// Numeric slice
type Numerics[T constraints.Ordered] struct {
	Slice[T]
}

func NewNumerics[T constraints.Ordered]() Numerics[T] {
	return Numerics[T]{}
}

func NewNumericsWithSize[T constraints.Ordered](size int) Numerics[T] {
	return Numerics[T]{Slice: make([]T, size)}
}

func NewNumericsWithCapacity[T constraints.Ordered](capacity int) Numerics[T] {
	return Numerics[T]{Slice: make([]T, 0, capacity)}
}

func (this Numerics[T]) Equal() bool {
	return false
}

func (this Numerics[T]) Sort() {
	this.Slice.Sort(func(a, b T) bool {
		return a < b
	})
}

func (this Numerics[T]) Find(w T) int {
	for n, v := range this.Data() {
		if v == w {
			return n
		}
	}
	return -1
}

func (this Numerics[T]) Includes(v T) bool {
	return this.Find(v) != -1
}

func (this *Numerics[T]) EraseAll(w T) {
	this.EraseIf(func(v T) bool {
		return v == w
	})
}
