// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rate

import (
	"context"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	time_ "github.com/searKing/golang/go/time"
)

const (
	d = 100 * time.Millisecond
)

type allow struct {
	n  int
	ok bool
}

func run(t *testing.T, lim *BurstLimiter, allows []allow) {
	for i, allow := range allows {
		ok := lim.AllowN(allow.n)
		if ok {
			lim.PutTokenN(allow.n)
		}
		if ok != allow.ok {
			t.Errorf("step %d: lim.AllowN(%v) = %v want %v",
				i, allow.n, ok, allow.ok)
		}
	}
}

func TestLimiterBurst1(t *testing.T) {
	run(t, NewFullBurstLimiter(1), []allow{
		{1, true},
		{2, false}, // burst size is 1, so n=2 always fails
	})
}

func TestLimiterBurst3(t *testing.T) {
	run(t, NewFullBurstLimiter(3), []allow{
		{1, true},
		{2, true},
		{3, true},
		{4, false}, // burst size is 1, so n=3 always fails
	})
}

func TestSimultaneousRequests(t *testing.T) {
	const (
		burst       = 5
		numRequests = 15
	)
	var (
		wg    sync.WaitGroup
		numOK = uint32(0)
	)

	// Very slow replenishing bucket.
	lim := NewFullBurstLimiter(burst)

	// Tries to take a token, atomically updates the counter and decreases the wait
	// group counter.
	f := func() {
		defer wg.Done()
		if ok := lim.Allow(); ok {
			atomic.AddUint32(&numOK, 1)
		}
	}

	wg.Add(numRequests)
	for i := 0; i < numRequests; i++ {
		go f()
	}
	wg.Wait()
	if numOK != burst {
		t.Errorf("numOK = %d, want %d", numOK, burst)
	}
}

func TestLongRunningQPS(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in short mode")
	}
	if runtime.GOOS == "openbsd" {
		t.Skip("low resolution time.Sleep invalidates test (golang.org/issue/14183)")
		return
	}

	// The test runs for a few seconds executing many requests and then checks
	// that overall number of requests is reasonable.
	const (
		burst = 100
	)
	var numOK = int32(0)

	lim := NewFullBurstLimiter(burst)

	var wg sync.WaitGroup
	f := func() {
		if ok := lim.Allow(); ok {
			atomic.AddInt32(&numOK, 1)
		}
		wg.Done()
	}

	start := time.Now()
	end := start.Add(5 * time.Second)
	for time.Now().Before(end) {
		wg.Add(1)
		go f()

		// This will still offer ~500 requests per second, but won't consume
		// outrageous amount of CPU.
		time.Sleep(2 * time.Millisecond)
	}
	wg.Wait()
	ideal := burst

	// We should never get more requests than allowed.
	if want := int32(ideal + 1); numOK > want {
		t.Errorf("numOK = %d, want %d (ideal %d)", numOK, want, ideal)
	}
	// We should get very close to the number of requests allowed.
	if want := int32(0.999 * float64(ideal)); numOK < want {
		t.Errorf("numOK = %d, want %d (ideal %d)", numOK, want, ideal)
	}
}

type request struct {
	n  int
	ok bool
}

func runReserve(t *testing.T, lim *BurstLimiter, req request) *Reservation {
	return runReserveMax(t, lim, req, time_.InfDuration)
}

func runReserveMax(t *testing.T, lim *BurstLimiter, req request, maxReserve time.Duration) *Reservation {
	ctx, cancel := context.WithTimeout(context.Background(), maxReserve)
	defer cancel()
	r := lim.reserveN(ctx, req.n, false)
	if r.ok {
		lim.PutTokenN(req.n)
	}
	if r.ok != req.ok {
		t.Errorf("lim.reserveN(%v, %v) = (%v) want (%v)",
			req.n, maxReserve, r.ok, req.ok)
	}
	return r
}

func TestSimpleReserve(t *testing.T) {
	lim := NewFullBurstLimiter(2)

	runReserve(t, lim, request{2, true})
	runReserve(t, lim, request{2, true})
	runReserve(t, lim, request{2, true})
}

func TestMix(t *testing.T) {
	lim := NewFullBurstLimiter(2)

	runReserve(t, lim, request{3, false}) // should return false because n > Burst
	runReserve(t, lim, request{2, true})
	run(t, lim, []allow{{3, false}}) // not enough tokens - don't allow
	runReserve(t, lim, request{2, true})
	run(t, lim, []allow{{1, true}})
}

func TestCancelInvalid(t *testing.T) {
	lim := NewFullBurstLimiter(2)

	runReserve(t, lim, request{2, true})
	r := runReserve(t, lim, request{3, false})
	r.Cancel()                           // should have no effect
	runReserve(t, lim, request{2, true}) // did not get extra tokens
}

func TestCancelLast(t *testing.T) {
	lim := NewFullBurstLimiter(2)

	runReserve(t, lim, request{2, true})
	r := runReserve(t, lim, request{2, true})
	r.Cancel() // got 2 tokens back
	runReserve(t, lim, request{2, true})
}

func TestCancelTooLate(t *testing.T) {
	lim := NewFullBurstLimiter(2)

	runReserve(t, lim, request{2, true})
	r := runReserve(t, lim, request{2, true})
	r.Cancel() // too late to cancel - should have no effect
	runReserve(t, lim, request{2, true})
}

func TestCancel0Tokens(t *testing.T) {
	lim := NewFullBurstLimiter(2)

	runReserve(t, lim, request{2, true})
	r := runReserve(t, lim, request{1, true})
	runReserve(t, lim, request{1, true})
	r.Cancel() // got 0 tokens back
	runReserve(t, lim, request{1, true})
}

func TestCancel1Token(t *testing.T) {
	lim := NewFullBurstLimiter(2)

	runReserve(t, lim, request{2, true})
	r := runReserve(t, lim, request{2, true})
	runReserve(t, lim, request{1, true})
	r.Cancel() // got 1 token back
	runReserve(t, lim, request{2, true})
}

func TestCancelMulti(t *testing.T) {
	lim := NewFullBurstLimiter(4)

	runReserve(t, lim, request{4, true})
	rA := runReserve(t, lim, request{3, true})
	runReserve(t, lim, request{1, true})
	rC := runReserve(t, lim, request{1, true})
	rC.Cancel() // get 1 token back
	rA.Cancel() // get 2 tokens back, as if C was never reserved
	runReserve(t, lim, request{3, true})
}

func TestReserveJumpBack(t *testing.T) {
	lim := NewFullBurstLimiter(2)

	runReserve(t, lim, request{2, true}) // start at t1
	runReserve(t, lim, request{1, true}) // should violate Limit,Burst
	runReserve(t, lim, request{2, true})
}

func TestReserveJumpBackCancel(t *testing.T) {
	lim := NewFullBurstLimiter(2)

	runReserve(t, lim, request{2, true}) // start at t1
	r := runReserve(t, lim, request{2, true})
	runReserve(t, lim, request{1, true})
	r.Cancel()                           // cancel at get 1 token back
	runReserve(t, lim, request{2, true}) // should violate Limit,Burst
}

func TestReserveMax(t *testing.T) {
	lim := NewFullBurstLimiter(2)
	maxT := d

	runReserveMax(t, lim, request{2, true}, maxT)
	runReserveMax(t, lim, request{1, true}, maxT) // reserve for close future
}

type wait struct {
	name   string
	ctx    context.Context
	n      int
	nilErr bool
}

func runWait(t *testing.T, lim *BurstLimiter, w wait) {
	err := lim.WaitN(w.ctx, w.n)
	if (w.nilErr && err != nil) || (!w.nilErr && err == nil) {
		errString := "<nil>"
		if !w.nilErr {
			errString = "<non-nil error>"
		}
		t.Errorf("lim.WaitN(%v, lim, %v) = %v; want %v",
			w.name, w.n, err, errString)
	}
}

func TestWaitSimple(t *testing.T) {
	lim := NewFullBurstLimiter(3)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	runWait(t, lim, wait{"already-cancelled", ctx, 1, false})

	runWait(t, lim, wait{"exceed-burst-error", context.Background(), 4, false})

	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	runWait(t, lim, wait{"act-now", ctx, 2, true})
	lim.PutTokenN(2)
	runWait(t, lim, wait{"act-later", ctx, 3, true})
}

func TestWaitTimeout(t *testing.T) {
	lim := NewFullBurstLimiter(3)

	ctx, cancel := context.WithTimeout(context.Background(), d)
	defer cancel()
	runWait(t, lim, wait{"act-now", ctx, 2, true})
	runWait(t, lim, wait{"w-timeout-err", ctx, 3, false})
}

func BenchmarkAllowN(b *testing.B) {
	lim := NewFullBurstLimiter(1)
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			lim.AllowN(1)
		}
	})
}

func BenchmarkWaitNNoDelay(b *testing.B) {
	lim := NewFullBurstLimiter(b.N)
	ctx := context.Background()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lim.WaitN(ctx, 1)
	}
}

func TestSimultaneousLongRequests(t *testing.T) {
	const (
		burst       = 5
		numRequests = 15
	)
	var (
		timeout = 1 * time.Millisecond
	)
	var (
		wg    sync.WaitGroup
		numOK = uint32(0)
	)

	// Very slow replenishing bucket.
	lim := NewFullBurstLimiter(burst)

	// Tries to take a token, atomically updates the counter and decreases the wait
	// group counter.
	f := func(i int) {
		if i < numRequests {
			defer wg.Done()
		}
		var limiterCtx = context.Background()
		var cancel context.CancelFunc
		if timeout > 0 {
			limiterCtx, cancel = context.WithTimeout(context.Background(), timeout)
			defer cancel()
		}
		if err := lim.Wait(limiterCtx); err != nil {
			t.Logf("#%d wait expect ok, got err %s", i, err)
		} else {
			atomic.AddUint32(&numOK, 1)
			defer lim.PutToken()
			t.Logf("#%d got token", i)
			time.Sleep(time.Second)
			t.Logf("#%d put token", i)
		}
	}

	wg.Add(numRequests)
	for i := 0; i < numRequests; i++ {
		go f(i)
	}
	wg.Wait()
	f(numRequests)
	if numOK != burst+1 {
		t.Errorf("numOK = %d, want %d", numOK, numRequests)
	}
}
