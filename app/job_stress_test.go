package app

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// mustReturn is generous on purpose: it is reached only on failure.
const mustReturn = 5 * time.Second

// Stop, Cancel, Wait, Sleep and Running from many goroutines on a running
// job: nothing may hang, panic or race.
func TestJobStressStorm(t *testing.T) {
	iters := 300
	if testing.Short() {
		iters = 30
	}
	for it := 0; it < iters; it++ {
		j := NewJob()
		j.Tick(func() {}, time.Duration(it%3)*50*time.Microsecond)
		var wg sync.WaitGroup
		for g := 0; g < 8; g++ {
			wg.Add(1)
			go func(g int) {
				defer wg.Done()
				switch g % 4 {
				case 0:
					j.Stop()
				case 1:
					j.Cancel()
					j.Stop()
				case 2:
					j.Wait()
					j.Running()
				case 3:
					j.Sleep(time.Hour) // returns only on cancellation
					j.Stop()
				}
			}(g)
		}
		within(t, mustReturn, "job storm", wg.Wait)
		waitNotRunning(t, j, mustReturn)
		if j.Do() {
			t.Fatal("Do() == true after Stop")
		}
	}
}

// A panic in fn is recovered: the job finishes, Stop does not hang.
func TestJobPanicInFn(t *testing.T) {
	j := NewJob()
	var finished atomic.Bool
	j.OnFinish(func() { finished.Store(true) })
	j.Run(func() { panic("test panic, expected in the log") })
	waitNotRunning(t, j, mustReturn)
	within(t, mustReturn, "Stop after panic", j.Stop)
	if j.Do() {
		t.Fatal("Do() == true after panic")
	}
	if !finished.Load() {
		t.Fatal("OnFinish not called")
	}
}

func TestJobInitError(t *testing.T) {
	j := NewJob()
	var ran, finished, notRun atomic.Bool
	j.OnInit(func() error { return errors.New("init failed") })
	j.OnFinish(func() { finished.Store(true) })
	j.OnNotRun(func() { notRun.Store(true) })
	j.Run(func() { ran.Store(true) })
	waitNotRunning(t, j, mustReturn)
	within(t, mustReturn, "Stop after init error", j.Stop)
	if ran.Load() || finished.Load() || notRun.Load() {
		t.Fatalf("ran=%v finished=%v notRun=%v, want all false", ran.Load(), finished.Load(), notRun.Load())
	}
}

func TestJobCancelBeforeRun(t *testing.T) {
	j := NewJob()
	var ran, started, finished, notRun atomic.Bool
	j.OnStart(func() { started.Store(true) })
	j.OnFinish(func() { finished.Store(true) })
	j.OnNotRun(func() { notRun.Store(true) })
	j.Cancel()
	j.Run(func() { ran.Store(true) })
	waitNotRunning(t, j, mustReturn)
	within(t, mustReturn, "Stop", j.Stop)
	if ran.Load() {
		t.Fatal("fn called after Cancel")
	}
	if !started.Load() || !notRun.Load() || !finished.Load() {
		t.Fatalf("started=%v notRun=%v finished=%v, want all true", started.Load(), notRun.Load(), finished.Load())
	}
}

// uws.Client.Run pattern: Cancel in OnStart.
func TestJobCancelInOnStart(t *testing.T) {
	j := NewJob()
	var calls atomic.Int32
	var notRun atomic.Bool
	j.OnStart(func() { j.Cancel() })
	j.OnNotRun(func() { notRun.Store(true) })
	j.Loop(func() { calls.Add(1) })
	waitNotRunning(t, j, mustReturn)
	within(t, mustReturn, "Stop", j.Stop)
	if n := calls.Load(); n != 0 {
		t.Fatalf("fn calls = %d, want 0", n)
	}
	if !notRun.Load() {
		t.Fatal("OnNotRun not called")
	}
}

// A finished job does not run again.
func TestJobRunAgain(t *testing.T) {
	j := NewJob()
	var calls, notRun atomic.Int32
	j.OnNotRun(func() { notRun.Add(1) })
	for i := 0; i < 2; i++ {
		j.Run(func() { calls.Add(1) })
		waitNotRunning(t, j, mustReturn)
	}
	within(t, mustReturn, "Stop", j.Stop)
	if calls.Load() != 1 || notRun.Load() != 1 {
		t.Fatalf("calls=%d notRun=%d, want 1 and 1", calls.Load(), notRun.Load())
	}
}

// uws.Client.Close pattern: Cancel, then Stop while fn is still busy.
// Stop does not wait for fn.
func TestJobCancelThenStopWhileBusy(t *testing.T) {
	j := NewJob()
	release := make(chan struct{})
	j.Run(func() { <-release })
	time.Sleep(10 * time.Millisecond)
	j.Cancel()
	within(t, mustReturn, "Stop after Cancel", j.Stop)
	if !j.Running() {
		t.Fatal("Running() == false while fn is busy")
	}
	close(release)
	waitNotRunning(t, j, mustReturn)
}

// Stop without Cancel waits until fn returns.
func TestJobStopWaitsForFn(t *testing.T) {
	j := NewJob()
	var finished atomic.Bool
	j.Run(func() {
		for j.Do() {
			j.Sleep(time.Hour)
		}
		time.Sleep(50 * time.Millisecond)
		finished.Store(true)
	})
	time.Sleep(10 * time.Millisecond)
	within(t, mustReturn, "Stop", j.Stop)
	if !finished.Load() {
		t.Fatal("Stop returned before fn finished")
	}
	if j.Running() {
		t.Fatal("Running() == true after Stop")
	}
}

// db.Db.Wait pattern: Sleep in a goroutine that is not the job's own.
func TestJobForeignSleepWakesOnStop(t *testing.T) {
	j := NewJob()
	j.Tick(func() {}, time.Second)
	woke := make(chan struct{})
	go func() {
		j.Sleep(time.Hour)
		close(woke)
	}()
	time.Sleep(20 * time.Millisecond)
	within(t, mustReturn, "Stop", j.Stop)
	within(t, mustReturn, "foreign Sleep", func() { <-woke })
}

func TestJobRepeatedStopAndCancel(t *testing.T) {
	j := NewJob()
	j.Tick(func() {}, time.Second)
	time.Sleep(10 * time.Millisecond)
	within(t, mustReturn, "repeated Stop and Cancel", func() {
		for i := 0; i < 100; i++ {
			j.Stop()
			j.Cancel()
			j.Stop()
		}
	})
	if j.Running() || j.Do() {
		t.Fatalf("Running()=%v Do()=%v, want false", j.Running(), j.Do())
	}
}
