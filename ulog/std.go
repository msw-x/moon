package ulog

import (
	"log"
	"strings"
)

func StdBridge(onMessage func(string)) *log.Logger {
	return log.New(NewStdWriter(onMessage), "", 0)
}

type StdWriter struct {
	onMessage func(string)
}

func NewStdWriter(onMessage func(string)) *StdWriter {
	return &StdWriter{
		onMessage: onMessage,
	}
}

func (o *StdWriter) Write(bytes []byte) (n int, err error) {
	o.onMessage(strings.TrimSuffix(string(bytes), "\n"))
	return len(bytes), nil
}
