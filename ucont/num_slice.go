package ucont

import "golang.org/x/exp/constraints"

// Numeric slice
type NumSlice[T constraints.Ordered] Slice[T]

func NewNumSliceWithSize[T constraints.Ordered](size int) NumSlice[T] {
	return NumSlice[T](make([]T, size))
}

func NewNumSliceWithCapacity[T constraints.Ordered](capacity int) NumSlice[T] {
	return NumSlice[T](make([]T, 0, capacity))
}

func (this NumSlice[T]) Equal() bool {
	return false
}

func (this *NumSlice[T]) Sort() {
	this.(Slice[T]).Sort(func(a, b T) bool {
		return a < b
	})
}
