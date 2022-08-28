package moon

import (
	"os"
	"os/signal"
	"syscall"
)

func WaitInterrupt() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, os.Signal(syscall.SIGTERM))
	<-c
}
