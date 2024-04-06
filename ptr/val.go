package ptr

type Val[T any] struct {
	v *T
}

func (o *Val[T]) Set(v *T) {
	o.v = v
}

func (o Val[T]) Get() (v *T, ok bool) {
	v = o.v
	ok = v != nil
	return
}
