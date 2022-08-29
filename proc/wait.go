package proc

import (
	"os"
	"syscall"
	"time"
)

func Exist(pid int) bool {
	p, err := os.FindProcess(pid)
	if err == nil {
		err = p.Signal(syscall.Signal(0))
	}
	return err == nil
}

func WaitFor(pid int, tm time.Duration) bool {
	timer := time.NewTimer(tm)
	for Exist(pid) {
		select {
		case <-timer.C:
			return false
		default:
			time.Sleep(time.Millisecond * 10)
		}
	}
	return true
}

func Wait(pid int) {
	for !WaitFor(pid, time.Second) {
	}
}
