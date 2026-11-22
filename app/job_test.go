package app

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// immediate is the upper bound for "returns right away".
const immediate = 50 * time.Millisecond

// within runs fn and fails the test if it does not return within timeout.
// It returns how long fn took.
func within(t *testing.T, timeout time.Duration, what string, fn func()) time.Duration {
	t.Helper()
	done := make(chan struct{})
	start := time.Now()
	go func() {
		fn()
		close(done)
	}()
	select {
	case <-done:
		return time.Since(start)
	case <-time.After(timeout):
		t.Fatalf("%s: not returned within %v", what, timeout)
		return 0
	}
}

// waitNotRunning polls Running() until false or fails on timeout.
func waitNotRunning(t *testing.T, j *Job, timeout time.Duration) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for j.Running() {
		if time.Now().After(deadline) {
			t.Fatalf("job still running after %v", timeout)
		}
		time.Sleep(time.Millisecond)
	}
}

// TS-9 (docs/moon.md TestTickStop)
func TestJobTickStop(t *testing.T) {
	j := NewJob()
	j.Tick(func() {}, 30*time.Second)
	time.Sleep(20 * time.Millisecond)
	if d := within(t, time.Second, "Stop", j.Stop); d > immediate {
		t.Fatalf("Stop took %v, want <= %v", d, immediate)
	}
	if j.Running() {
		t.Fatal("Running() == true after Stop")
	}
}

// TS-10 (docs/moon.md TestRunningAfterStop)
func TestJobRunningAfterStop(t *testing.T) {
	j := NewJob()
	j.Tick(func() {}, 30*time.Second)
	time.Sleep(20 * time.Millisecond)
	within(t, time.Second, "Stop", j.Stop)
	if j.Running() {
		t.Fatal("Running() == true after Stop")
	}
}

// TS-11 (docs/moon.md TestRunSelfFinish)
func TestJobRunSelfFinish(t *testing.T) {
	j := NewJob()
	j.Run(func() { time.Sleep(5 * time.Millisecond) })
	waitNotRunning(t, j, time.Second)
	within(t, immediate, "Stop after self finish", j.Stop)
}

// TS-12
func TestJobTickInterval(t *testing.T) {
	const interval = 150 * time.Millisecond
	var mu sync.Mutex
	var calls []time.Time
	j := NewJob()
	j.Tick(func() {
		mu.Lock()
		calls = append(calls, time.Now())
		mu.Unlock()
	}, interval)
	time.Sleep(700 * time.Millisecond)
	within(t, time.Second, "Stop", j.Stop)

	mu.Lock()
	got := append([]time.Time(nil), calls...)
	mu.Unlock()
	if len(got) < 4 || len(got) > 5 {
		t.Fatalf("calls = %d, want 4..5", len(got))
	}
	for i := 1; i < len(got); i++ {
		if gap := got[i].Sub(got[i-1]); gap < interval {
			t.Fatalf("gap %d = %v, want >= %v", i, gap, interval)
		}
	}

	time.Sleep(2 * interval)
	mu.Lock()
	after := len(calls)
	mu.Unlock()
	if after != len(got) {
		t.Fatalf("calls after Stop: %d -> %d", len(got), after)
	}
}

// TS-13
func TestJobWaitWakesOnCancel(t *testing.T) {
	j := NewJob()
	j.Loop(func() { j.Sleep(time.Hour) })
	defer j.Stop()

	waited := make(chan time.Time, 1)
	go func() {
		j.Wait()
		waited <- time.Now()
	}()
	time.Sleep(20 * time.Millisecond)
	canceledAt := time.Now()
	j.Cancel()
	select {
	case at := <-waited:
		if lag := at.Sub(canceledAt); lag > immediate {
			t.Fatalf("Wait returned %v after Cancel, want <= %v", lag, immediate)
		}
	case <-time.After(time.Second):
		t.Fatal("Wait did not return after Cancel")
	}
}

// TS-14
func TestJobStopNotStartedAndSelfCancel(t *testing.T) {
	t.Run("not started", func(t *testing.T) {
		j := NewJob()
		within(t, immediate, "Stop of not started job", j.Stop)
	})
	t.Run("cancel inside job", func(t *testing.T) {
		j := NewJob()
		var calls atomic.Int32
		var finished atomic.Bool
		j.OnFinish(func() { finished.Store(true) })
		j.Loop(func() {
			calls.Add(1)
			j.Cancel()
		})
		waitNotRunning(t, j, time.Second)
		if n := calls.Load(); n != 1 {
			t.Fatalf("fn calls = %d, want 1", n)
		}
		if !finished.Load() {
			t.Fatal("OnFinish not called")
		}
		within(t, immediate, "Stop after self cancel", j.Stop)
	})
}
