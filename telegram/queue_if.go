package telegram

type QueueIf[T any] interface {
	SetCapacity(int)
	Capacity() int
	Size() int
	Reset()
	Push(T) bool
	Pop() (T, bool)
}
