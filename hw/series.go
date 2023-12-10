package hw

import "github.com/msw-x/moon/umath"

type Series[T umath.Number] struct {
	i int
	l []T
}

func NewSeries[T umath.Number](n int) *Series[T] {
	o := new(Series[T])
	o.l = make([]T, 0, n)
	return o
}

func (o *Series[T]) List() []T {
	return o.l
}

func (o *Series[T]) Add(v T) {
	if o.Full() {
		o.l[o.i] = v
	} else {
		o.l = append(o.l, v)
	}
	o.i++
	if o.i == o.Capacity() {
		o.i = 0
	}
}

func (o *Series[T]) Size() int {
	return len(o.l)
}

func (o *Series[T]) Capacity() int {
	return cap(o.l)
}

func (o *Series[T]) Full() bool {
	return o.Size() == o.Capacity()
}

func (o *Series[T]) Empty() bool {
	return o.Size() == 0
}

func (o *Series[T]) Average() T {
	if o.Empty() {
		return 0
	}
	var sum float64
	for n := 0; n != o.Size(); n++ {
		sum += float64(o.l[n])
	}
	return T(sum / float64(o.Size()))
}
