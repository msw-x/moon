package moon

import (
	"os"
	"os/signal"
	"syscall"
)

func WaitInterrupt() os.Signal {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGABRT)
	return <-c
}
