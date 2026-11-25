package usync

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// immediate is the upper bound for "returns right away".
const immediate = 50 * time.Millisecond

// waitDone fails the test if ch is not closed within timeout.
func waitDone(t *testing.T, ch <-chan struct{}, timeout time.Duration, what string) {
	t.Helper()
	select {
	case <-ch:
	case <-time.After(timeout):
		t.Fatalf("%s: not returned within %v", what, timeout)
	}
}

// TS-1
func TestDoSleepFullInterval(t *testing.T) {
	d := NewDo()
	const interval = 150 * time.Millisecond
	start := time.Now()
	d.Sleep(interval)
	elapsed := time.Since(start)
	if elapsed < interval || elapsed >= 250*time.Millisecond {
		t.Fatalf("Sleep(%v) took %v, want [150ms, 250ms)", interval, elapsed)
	}
}

// TS-2
func TestDoSleepWakesOnCancel(t *testing.T) {
	cases := []struct {
		name   string
		cancel func(d *Do)
	}{
		{"Cancel", func(d *Do) { d.Cancel() }},
		{"Notify", func(d *Do) { d.Notify() }},
		{"Stop", func(d *Do) {
			go func() {
				time.Sleep(5 * time.Millisecond)
				d.Notify()
			}()
			d.Stop()
		}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			d := NewDo()
			woke := make(chan time.Time, 1)
			go func() {
				d.Sleep(10 * time.Second)
				woke <- time.Now()
			}()
			time.Sleep(20 * time.Millisecond)
			canceledAt := time.Now()
			go c.cancel(d)
			select {
			case at := <-woke:
				if lag := at.Sub(canceledAt); lag > immediate {
					t.Fatalf("Sleep woke %v after cancel, want <= %v", lag, immediate)
				}
			case <-time.After(time.Second):
				t.Fatal("Sleep did not wake on cancel")
			}
			if d.Do() {
				t.Fatal("Do() == true after cancel")
			}
		})
	}
}

// TS-3
func TestDoSleepNoWait(t *testing.T) {
	canceled := NewDo()
	canceled.Cancel()
	active := NewDo()
	cases := []struct {
		name string
		fn   func()
	}{
		{"canceled", func() { canceled.Sleep(time.Second) }},
		{"zero", func() { active.Sleep(0) }},
		{"negative", func() { active.Sleep(-1) }},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			start := time.Now()
			c.fn()
			if elapsed := time.Since(start); elapsed > immediate {
				t.Fatalf("Sleep took %v, want <= %v", elapsed, immediate)
			}
		})
	}
	if !active.Do() {
		t.Fatal("Sleep changed Do()")
	}
}

// TS-4
func TestDoCancelWakesAllSleepers(t *testing.T) {
	d := NewDo()
	const n = 10
	var started, finished sync.WaitGroup
	started.Add(n)
	finished.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			started.Done()
			d.Sleep(10 * time.Second)
			finished.Done()
		}()
	}
	started.Wait()
	time.Sleep(20 * time.Millisecond)
	d.Cancel()
	all := make(chan struct{})
	go func() {
		finished.Wait()
		close(all)
	}()
	waitDone(t, all, immediate, "all sleepers")
}

// owner emulates a job goroutine: sleeps while Do() is true, works for
// `work`, sets flag and calls Notify.
func owner(d *Do, work time.Duration, flag *atomic.Bool) {
	go func() {
		for d.Do() {
			d.Sleep(time.Hour)
		}
		time.Sleep(work)
		flag.Store(true)
		d.Notify()
	}()
}

// TS-5
func TestDoStopWaitsNotify(t *testing.T) {
	d := NewDo()
	var notified atomic.Bool
	owner(d, 50*time.Millisecond, &notified)
	time.Sleep(10 * time.Millisecond)
	start := time.Now()
	stopped := make(chan struct{})
	go func() {
		d.Stop()
		close(stopped)
	}()
	waitDone(t, stopped, time.Second, "Stop")
	if !notified.Load() {
		t.Fatal("Stop returned before Notify")
	}
	if elapsed := time.Since(start); elapsed > 200*time.Millisecond {
		t.Fatalf("Stop took %v, want <= 200ms", elapsed)
	}
}

// TS-6
func TestDoStopAfterNotify(t *testing.T) {
	d := NewDo()
	d.Notify()
	for i := 0; i < 2; i++ {
		stopped := make(chan struct{})
		go func() {
			d.Stop()
			close(stopped)
		}()
		waitDone(t, stopped, immediate, "Stop after Notify")
	}
	d.Notify() // must not panic
	d.Cancel()
	if d.Do() {
		t.Fatal("Do() == true after Notify")
	}
}

// TS-7
func TestDoConcurrentStop(t *testing.T) {
	d := NewDo()
	var notified atomic.Bool
	owner(d, 50*time.Millisecond, &notified)
	time.Sleep(10 * time.Millisecond)
	const n = 10
	var wg sync.WaitGroup
	var afterNotify atomic.Int32
	start := make(chan struct{})
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			<-start
			d.Stop()
			if notified.Load() {
				afterNotify.Add(1)
			}
		}()
	}
	close(start)
	all := make(chan struct{})
	go func() {
		wg.Wait()
		close(all)
	}()
	waitDone(t, all, time.Second, "concurrent Stop")
	if afterNotify.Load() == 0 {
		t.Fatal("no Stop returned after Notify")
	}
}

// TS-8
func TestDoCancelThenStopNoWait(t *testing.T) {
	t.Run("external", func(t *testing.T) {
		d := NewDo()
		release := make(chan struct{})
		exited := make(chan struct{})
		go func() {
			<-release // owner is busy outside Sleep
			d.Notify()
			close(exited)
		}()
		d.Cancel()
		stopped := make(chan struct{})
		go func() {
			d.Stop()
			close(stopped)
		}()
		waitDone(t, stopped, immediate, "Stop after Cancel")
		close(release)
		waitDone(t, exited, time.Second, "owner exit")
	})
	t.Run("owner", func(t *testing.T) {
		d := NewDo()
		exited := make(chan struct{})
		go func() {
			d.Cancel()
			d.Stop()
			d.Notify()
			close(exited)
		}()
		waitDone(t, exited, immediate, "Cancel+Stop in owner")
	})
}
