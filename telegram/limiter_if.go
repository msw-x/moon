package telegram

import (
	"time"
)

type LimiterIf[T any] interface {
	SetInterval(time.Duration)
	SetQueue(QueueIf[T])
	Queue() QueueIf[T]
	Close()
	Push(T) bool
	Discard(int) int
}
