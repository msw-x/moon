package usync

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

const (
	// mustReturn is generous on purpose: it is reached only on failure,
	// so a slow machine cannot produce a false failure.
	mustReturn = 5 * time.Second
	// mustBlock is how long a call is observed to prove that it blocks.
	mustBlock = 30 * time.Millisecond
)

// call is a Do method running in its own goroutine.
type call struct {
	op   byte
	done chan struct{}
}

func startCall(d *Do, op byte) *call {
	c := &call{op: op, done: make(chan struct{})}
	go func() {
		defer close(c.done)
		switch op {
		case 'C':
			d.Cancel()
		case 'N':
			d.Notify()
		case 'S':
			d.Stop()
		case 'Z':
			d.Sleep(time.Minute)
		}
	}()
	return c
}

func (c *call) returned(within time.Duration) bool {
	timer := time.NewTimer(within)
	defer timer.Stop()
	select {
	case <-c.done:
		return true
	case <-timer.C:
		return false
	}
}

// checkSequence runs ops one after another and compares every step with
// the contract of Do. A blocked call stays blocked in its own goroutine,
// like a call made from another goroutine. Ops: C - Cancel, N - Notify,
// S - Stop, Z - Sleep(time.Minute).
//
// Contract: Cancel and Notify never block. Stop blocks only when Do is
// still active, until Notify. Sleep blocks only while Do is active, any
// of Cancel, Notify, Stop wakes it.
func checkSequence(seq string) error {
	d := NewDo()
	active := true
	var stopper *call
	var sleepers []*call
	for i := 0; i < len(seq); i++ {
		op := seq[i]
		step := fmt.Sprintf("%s[%d]=%c", seq, i, op)
		c := startCall(d, op)
		wantBlock := active && (op == 'S' || op == 'Z')
		if wantBlock {
			if op == 'S' {
				// Stop has taken effect once the flag is cleared; a fixed
				// pause is not enough when the goroutine starts late.
				deadline := time.Now().Add(mustReturn)
				for d.Do() && time.Now().Before(deadline) {
					time.Sleep(100 * time.Microsecond)
				}
			}
			if c.returned(mustBlock) {
				return fmt.Errorf("%s: returned, want blocked", step)
			}
		} else if !c.returned(mustReturn) {
			return fmt.Errorf("%s: blocked, want returned", step)
		}
		switch op {
		case 'C':
			active = false
		case 'N':
			active = false
			if stopper != nil {
				if !stopper.returned(mustReturn) {
					return fmt.Errorf("%s: Notify did not release Stop", step)
				}
				stopper = nil
			}
		case 'S':
			if wantBlock {
				stopper = c
			}
			active = false
		case 'Z':
			if wantBlock {
				sleepers = append(sleepers, c)
			}
		}
		if !active {
			for _, s := range sleepers {
				if !s.returned(mustReturn) {
					return fmt.Errorf("%s: Sleep not woken", step)
				}
			}
			sleepers = nil
		}
		if stopper != nil && stopper.returned(0) {
			return fmt.Errorf("%s: Stop returned before Notify", step)
		}
		if d.Do() != active {
			return fmt.Errorf("%s: Do() = %v, want %v", step, d.Do(), active)
		}
	}
	d.Notify()
	for _, c := range append(sleepers, stopper) {
		if c != nil && !c.returned(mustReturn) {
			return fmt.Errorf("%s: final Notify did not release %c", seq, c.op)
		}
	}
	return nil
}

func sequences(alphabet string, maxLen int) []string {
	var out []string
	var rec func(prefix string)
	rec = func(prefix string) {
		if prefix != "" {
			out = append(out, prefix)
		}
		if len(prefix) == maxLen {
			return
		}
		for i := 0; i < len(alphabet); i++ {
			rec(prefix + string(alphabet[i]))
		}
	}
	rec("")
	return out
}

// Every order of Cancel, Notify, Stop and Sleep, repeats included, follows
// the contract: no panic, no unexpected block, no early return.
func TestDoAllSequences(t *testing.T) {
	maxLen := 5
	if testing.Short() {
		maxLen = 4
	}
	seqs := sequences("CNSZ", maxLen)
	errs := make(chan error, len(seqs))
	limit := make(chan struct{}, 256)
	var wg sync.WaitGroup
	for _, seq := range seqs {
		wg.Add(1)
		go func(seq string) {
			defer wg.Done()
			limit <- struct{}{}
			defer func() { <-limit }()
			if err := checkSequence(seq); err != nil {
				errs <- err
			}
		}(seq)
	}
	wg.Wait()
	close(errs)
	shown := 0
	for err := range errs {
		if shown++; shown > 20 {
			t.Fatal("too many failures")
		}
		t.Error(err)
	}
}

func waitGroup(t *testing.T, wg *sync.WaitGroup, what string) {
	t.Helper()
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()
	waitDone(t, done, mustReturn, what)
}

// Random concurrent mix of all methods. The owner goroutine calls Notify,
// as a job goroutine does. Nothing may hang, panic or race.
func TestDoStressRandomConcurrent(t *testing.T) {
	iters := 2000
	if testing.Short() {
		iters = 200
	}
	for it := 0; it < iters; it++ {
		d := NewDo()
		var wg sync.WaitGroup
		ownerSleep := time.Duration(it%300+1) * time.Microsecond
		wg.Add(1)
		go func() {
			defer wg.Done()
			for n := 0; d.Do() && n < 3; n++ {
				d.Sleep(ownerSleep)
			}
			d.Notify()
		}()
		for g := 0; g < 8; g++ {
			wg.Add(1)
			go func(seed int64) {
				defer wg.Done()
				r := rand.New(rand.NewSource(seed))
				for k := 0; k < 12; k++ {
					switch r.Intn(6) {
					case 0:
						d.Cancel()
					case 1:
						d.Notify()
					case 2:
						d.Stop()
					case 3:
						d.Sleep(time.Duration(r.Intn(200)) * time.Microsecond)
					case 4:
						d.Sleep(time.Hour) // returns only on cancellation
					case 5:
						d.Do()
					}
				}
			}(int64(it*100 + g))
		}
		waitGroup(t, &wg, "random mix")
		if d.Do() {
			t.Fatal("Do() == true after Notify")
		}
	}
}

// Many concurrent Stop calls and one Notify from the owner: nobody hangs,
// and the Stop that canceled Do returns only after Notify.
func TestDoStressStopStorm(t *testing.T) {
	iters := 300
	if testing.Short() {
		iters = 30
	}
	for it := 0; it < iters; it++ {
		d := NewDo()
		var notified atomic.Bool
		var afterNotify atomic.Int32
		var wg sync.WaitGroup
		start := make(chan struct{})
		for g := 0; g < 16; g++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				<-start
				d.Stop()
				if notified.Load() {
					afterNotify.Add(1)
				}
			}()
		}
		owner(d, 2*time.Millisecond, &notified)
		close(start)
		waitGroup(t, &wg, "stop storm")
		if afterNotify.Load() == 0 {
			t.Fatalf("iteration %d: no Stop waited for Notify", it)
		}
	}
}

// Without Notify exactly one Stop stays blocked: the one that canceled Do.
// Every other call returns, a later Notify releases the blocked one.
func TestDoStopWithoutNotifyBlocksOne(t *testing.T) {
	iters := 100
	if testing.Short() {
		iters = 10
	}
	const n = 12
	for it := 0; it < iters; it++ {
		d := NewDo()
		var returned atomic.Int32
		var wg sync.WaitGroup
		for g := 0; g < n; g++ {
			wg.Add(1)
			go func(g int) {
				defer wg.Done()
				if g%3 == 0 {
					d.Sleep(time.Hour)
				}
				d.Stop()
				d.Cancel()
				d.Stop()
				returned.Add(1)
			}(g)
		}
		deadline := time.Now().Add(mustReturn)
		for returned.Load() < n-1 && time.Now().Before(deadline) {
			time.Sleep(100 * time.Microsecond)
		}
		time.Sleep(2 * time.Millisecond)
		if got := returned.Load(); got != n-1 {
			t.Fatalf("iteration %d: %d goroutines returned, want %d", it, got, n-1)
		}
		d.Notify()
		d.Notify()
		waitGroup(t, &wg, "Stop released by Notify")
	}
}

// Repeated calls from one goroutine in every order that does not start
// with Stop on an active Do.
func TestDoRepeatedCalls(t *testing.T) {
	ops := map[byte]func(*Do){'C': (*Do).Cancel, 'N': (*Do).Notify, 'S': (*Do).Stop}
	for _, order := range []string{"CNS", "CSN", "NCS", "NSC"} {
		d := NewDo()
		done := make(chan struct{})
		go func() {
			defer close(done)
			for rep := 0; rep < 1000; rep++ {
				for i := 0; i < len(order); i++ {
					ops[order[i]](d)
				}
				d.Sleep(time.Hour)
			}
		}()
		waitDone(t, done, mustReturn, "order "+order)
		if d.Do() {
			t.Fatalf("order %s: Do() == true", order)
		}
	}
}

// The zero value is a canceled Do: no method panics or blocks.
func TestDoZeroValue(t *testing.T) {
	for _, order := range []string{"CNSZ", "NCSZ", "SZNC", "ZSCN"} {
		d := new(Do)
		done := make(chan struct{})
		go func() {
			defer close(done)
			for rep := 0; rep < 3; rep++ {
				for i := 0; i < len(order); i++ {
					switch order[i] {
					case 'C':
						d.Cancel()
					case 'N':
						d.Notify()
					case 'S':
						d.Stop()
					case 'Z':
						d.Sleep(time.Hour)
					}
				}
			}
		}()
		waitDone(t, done, mustReturn, "zero value, order "+order)
		if d.Do() {
			t.Fatalf("zero value, order %s: Do() == true", order)
		}
	}
}

// Timers are reused between Sleep calls. A timer taken from a Sleep that
// was canceled right at its expiry must not shorten a later Sleep.
func TestDoSleepNotShortenedByReusedTimer(t *testing.T) {
	rounds := 100
	if testing.Short() {
		rounds = 10
	}
	const probe = 5 * time.Millisecond
	for round := 0; round < rounds; round++ {
		var wg sync.WaitGroup
		// Cancel races with the expiry of short sleeps.
		for g := 0; g < 16; g++ {
			d := NewDo()
			pause := time.Duration(50+g*10) * time.Microsecond
			wg.Add(2)
			go func() {
				defer wg.Done()
				d.Sleep(pause)
			}()
			go func() {
				defer wg.Done()
				time.Sleep(pause)
				d.Cancel()
			}()
		}
		wg.Wait()
		// Sleeps on fresh Do take timers from the pool.
		var short atomic.Int64
		for g := 0; g < 16; g++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				d := NewDo()
				start := time.Now()
				d.Sleep(probe)
				if elapsed := time.Since(start); elapsed < probe {
					short.Store(int64(elapsed))
				}
			}()
		}
		wg.Wait()
		if v := short.Load(); v != 0 {
			t.Fatalf("round %d: Sleep(%v) returned after %v", round, probe, time.Duration(v))
		}
	}
}
