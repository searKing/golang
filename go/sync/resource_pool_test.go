// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sync_test

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	sync_ "github.com/searKing/golang/go/sync"
)

// res is a resource made by a recorder.
type res struct {
	key    string
	id     int
	closed atomic.Bool
	users  atomic.Int32 // leases using the resource right now
}

// recorder counts the resources a ResourcePool makes and closes.
type recorder struct {
	mu     sync.Mutex
	made   []*res
	closed []*res
	newErr error // returned by New, if set
}

func (r *recorder) New(ctx context.Context, key string) (*res, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.newErr != nil {
		return nil, r.newErr
	}
	v := &res{key: key, id: len(r.made)}
	r.made = append(r.made, v)
	return v, nil
}

func (r *recorder) Close(key string, v *res) error {
	if v.closed.Swap(true) {
		panic(fmt.Sprintf("resource %s#%d closed twice", v.key, v.id))
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.closed = append(r.closed, v)
	return nil
}

func (r *recorder) counts() (made, closed int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.made), len(r.closed)
}

func newPool(r *recorder) *sync_.ResourcePool[string, *res] {
	return &sync_.ResourcePool[string, *res]{New: r.New, Close: r.Close}
}

func mustGet(t *testing.T, p *sync_.ResourcePool[string, *res], keys ...string) *sync_.Lease[string, *res] {
	t.Helper()
	l, err := p.Get(context.Background(), keys...)
	if err != nil {
		t.Fatalf("Get(%q): %v", keys, err)
	}
	return l
}

// getAsync calls Get on a new goroutine.
func getAsync(ctx context.Context, p *sync_.ResourcePool[string, *res], keys ...string) <-chan *sync_.Lease[string, *res] {
	c := make(chan *sync_.Lease[string, *res], 1)
	go func() {
		l, err := p.Get(ctx, keys...)
		if err != nil {
			close(c)
			return
		}
		c <- l
	}()
	return c
}

// expectBlocked fails the test if c delivers within a short while.
func expectBlocked[E any](t *testing.T, c <-chan E) {
	t.Helper()
	select {
	case v, ok := <-c:
		t.Fatalf("expected to block, got %v, %v", v, ok)
	case <-time.After(50 * time.Millisecond):
	}
}

func receive[E any](t *testing.T, c <-chan E) E {
	t.Helper()
	select {
	case v, ok := <-c:
		if !ok {
			t.Fatal("channel closed")
		}
		return v
	case <-time.After(5 * time.Second):
		t.Fatal("timed out")
	}
	panic("unreachable")
}

func TestResourcePool_Reuse(t *testing.T) {
	var r recorder
	p := newPool(&r)

	l := mustGet(t, p, "a")
	v := l.Value()
	if l.Key() != "a" || v.key != "a" {
		t.Fatalf("got key %q, value key %q; want a", l.Key(), v.key)
	}
	l.Release()

	l = mustGet(t, p, "a")
	if l.Value() != v {
		t.Errorf("Get after Release made a new resource, want the idle one reused")
	}
	l.Release()
	if made, closed := r.counts(); made != 1 || closed != 0 {
		t.Errorf("made %d, closed %d; want 1, 0", made, closed)
	}
}

func TestResourcePool_NoKeysUsesZeroKey(t *testing.T) {
	var r recorder
	p := newPool(&r)
	l := mustGet(t, p)
	if l.Key() != "" {
		t.Errorf("Key() = %q, want zero key", l.Key())
	}
	l.Release()
}

func TestResourcePool_NilNew(t *testing.T) {
	var p sync_.ResourcePool[string, int]
	if _, err := p.Get(context.Background(), "a"); err == nil {
		t.Fatal("Get with nil New succeeded, want error")
	}
}

func TestResourcePool_ReleaseIdempotent(t *testing.T) {
	var r recorder
	p := newPool(&r)
	l := mustGet(t, p, "a")
	l.Release()
	l.Release()
	l.Discard()
	if s := p.Stats()["a"]; s.Idle != 1 || s.Total != 1 || s.Leases != 0 {
		t.Errorf("Stats = %+v, want 1 idle of 1 and no leases", s)
	}
	if _, closed := r.counts(); closed != 0 {
		t.Errorf("Discard after Release closed the resource")
	}
}

func TestResourcePool_MaxResourcesPerKeyBlocks(t *testing.T) {
	var r recorder
	p := newPool(&r)
	p.MaxResourcesPerKey = 1

	l := mustGet(t, p, "a")

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	if _, err := p.Get(ctx, "a"); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("Get over limit: err = %v, want DeadlineExceeded", err)
	}

	c := getAsync(context.Background(), p, "a")
	expectBlocked(t, c)
	v := l.Value()
	l.Release()
	l2 := receive(t, c)
	if l2.Value() != v {
		t.Errorf("waiter got a new resource, want the released one")
	}
	l2.Release()
	if made, _ := r.counts(); made != 1 {
		t.Errorf("made %d, want 1", made)
	}
}

func TestResourcePool_WaitersFIFO(t *testing.T) {
	var r recorder
	p := newPool(&r)
	p.MaxResourcesPerKey = 1
	l := mustGet(t, p, "a")

	var order []int
	var mu sync.Mutex
	var wg sync.WaitGroup
	for i := range 3 {
		wg.Go(func() {
			l, err := p.Get(context.Background(), "a")
			if err != nil {
				t.Error(err)
				return
			}
			mu.Lock()
			order = append(order, i)
			mu.Unlock()
			l.Release()
		})
		// Let goroutine i queue before i+1.
		waitFor(t, func() bool { return p.Stats()["a"].Waiting == i+1 })
	}
	l.Release()
	wg.Wait()
	if fmt.Sprint(order) != "[0 1 2]" {
		t.Errorf("waiters served in order %v, want [0 1 2]", order)
	}
}

func waitFor(t *testing.T, cond func() bool) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for !cond() {
		if time.Now().After(deadline) {
			t.Fatal("condition not met in time")
		}
		time.Sleep(time.Millisecond)
	}
}

func TestResourcePool_CancelWhileWaiting(t *testing.T) {
	var r recorder
	p := newPool(&r)
	p.MaxResourcesPerKey = 1
	l := mustGet(t, p, "a")

	ctx, cancel := context.WithCancel(context.Background())
	errc := make(chan error, 1)
	go func() {
		_, err := p.Get(ctx, "a")
		errc <- err
	}()
	waitFor(t, func() bool { return p.Stats()["a"].Waiting == 1 })
	cancel()
	if err := receive(t, errc); !errors.Is(err, context.Canceled) {
		t.Fatalf("err = %v, want Canceled", err)
	}
	if s := p.Stats()["a"]; s.Waiting != 0 {
		t.Errorf("Waiting = %d after cancel, want 0", s.Waiting)
	}
	l.Release()
	if s := p.Stats()["a"]; s.Idle != 1 {
		t.Errorf("Stats = %+v, want the resource idle", s)
	}
}

func TestResourcePool_PicksMostIdleKey(t *testing.T) {
	var r recorder
	p := newPool(&r)
	if err := p.Prefill(context.Background(), "a", 1); err != nil {
		t.Fatal(err)
	}
	if err := p.Prefill(context.Background(), "b", 2); err != nil {
		t.Fatal(err)
	}
	l := mustGet(t, p, "a", "b")
	if l.Key() != "b" {
		t.Errorf("Key() = %q, want b, the key with most idle resources", l.Key())
	}
	l.Release()
}

func TestResourcePool_TiesBrokenByKeyOrder(t *testing.T) {
	var r recorder
	p := newPool(&r)
	l := mustGet(t, p, "b", "a")
	if l.Key() != "b" {
		t.Errorf("Key() = %q, want b, the first key", l.Key())
	}
	l.Release()
}

func TestResourcePool_CreatesOnLeastLoadedKey(t *testing.T) {
	var r recorder
	p := newPool(&r)
	la := mustGet(t, p, "a")
	l := mustGet(t, p, "a", "b")
	if l.Key() != "b" {
		t.Errorf("Key() = %q, want b, the key with fewest resources", l.Key())
	}
	la.Release()
	l.Release()
}

func TestResourcePool_UnknownKeyIsCreated(t *testing.T) {
	var r recorder
	p := newPool(&r)
	p.MaxResourcesPerKey = 1
	l := mustGet(t, p, "never-seen")
	l.Release()
}

func TestResourcePool_WaiterServedByAnyKey(t *testing.T) {
	var r recorder
	p := newPool(&r)
	p.MaxResourcesPerKey = 1
	la := mustGet(t, p, "a")
	lb := mustGet(t, p, "b")

	c := getAsync(context.Background(), p, "a", "b")
	expectBlocked(t, c)
	lb.Release()
	if l := receive(t, c); l.Key() != "b" {
		t.Errorf("Key() = %q, want b", l.Key())
	} else {
		l.Release()
	}
	la.Release()
}

func TestResourcePool_Discard(t *testing.T) {
	var r recorder
	p := newPool(&r)
	p.MaxResourcesPerKey = 1
	l := mustGet(t, p, "a")
	v := l.Value()

	c := getAsync(context.Background(), p, "a")
	expectBlocked(t, c)
	l.Discard()
	if !v.closed.Load() {
		t.Error("Discard did not close the resource")
	}
	l2 := receive(t, c)
	if l2.Value() == v {
		t.Error("waiter got the discarded resource")
	}
	l2.Release()
	l.Release() // no-op
	if made, closed := r.counts(); made != 2 || closed != 1 {
		t.Errorf("made %d, closed %d; want 2, 1", made, closed)
	}
}

func TestResourcePool_DisableKeepAlives(t *testing.T) {
	var r recorder
	p := newPool(&r)
	p.DisableKeepAlives = true
	p.MaxConcurrentLeases = 2
	l := mustGet(t, p, "a")
	l2 := mustGet(t, p, "a")
	if l.Value() == l2.Value() {
		t.Error("resource shared with DisableKeepAlives")
	}
	l.Release()
	l2.Release()
	if made, closed := r.counts(); made != 2 || closed != 2 {
		t.Errorf("made %d, closed %d; want 2, 2", made, closed)
	}
}

func TestResourcePool_MaxIdleResourcesPerKey(t *testing.T) {
	var r recorder
	p := newPool(&r)
	p.MaxIdleResourcesPerKey = 1
	l := mustGet(t, p, "a")
	l2 := mustGet(t, p, "a")
	l.Release()
	l2.Release()
	if s := p.Stats()["a"]; s.Idle != 1 || s.Total != 1 {
		t.Errorf("Stats = %+v, want 1 idle of 1", s)
	}
	if _, closed := r.counts(); closed != 1 {
		t.Errorf("closed %d, want 1", closed)
	}
}

func TestResourcePool_IdleResourceTimeout(t *testing.T) {
	var r recorder
	p := newPool(&r)
	p.IdleResourceTimeout = 10 * time.Millisecond
	l := mustGet(t, p, "a")
	v := l.Value()
	l.Release()
	waitFor(t, v.closed.Load)
	waitFor(t, func() bool { return p.Stats()["a"].Total == 0 })

	// A resource in use is not closed.
	l = mustGet(t, p, "a")
	time.Sleep(30 * time.Millisecond)
	if l.Value().closed.Load() {
		t.Error("resource in use closed by IdleResourceTimeout")
	}
	l.Release()
}

func TestResourcePool_MaxConcurrentLeases(t *testing.T) {
	var r recorder
	p := newPool(&r)
	p.MaxResourcesPerKey = 1
	p.MaxConcurrentLeases = 2

	l := mustGet(t, p, "a")
	l2 := mustGet(t, p, "a")
	if l.Value() != l2.Value() {
		t.Error("leases not shared")
	}
	c := getAsync(context.Background(), p, "a")
	expectBlocked(t, c)
	l.Release()
	l3 := receive(t, c)
	if s := p.Stats()["a"]; s.Leases != 2 || s.Total != 1 || s.Idle != 0 {
		t.Errorf("Stats = %+v, want 2 leases on 1 resource", s)
	}
	l2.Release()
	l3.Release()
	if s := p.Stats()["a"]; s.Leases != 0 || s.Idle != 1 {
		t.Errorf("Stats = %+v, want 1 idle", s)
	}
}

func TestResourcePool_PrefersIdleOverSharing(t *testing.T) {
	var r recorder
	p := newPool(&r)
	p.MaxConcurrentLeases = 2
	if err := p.Prefill(context.Background(), "a", 2); err != nil {
		t.Fatal(err)
	}
	l := mustGet(t, p, "a")
	l2 := mustGet(t, p, "a")
	if l.Value() == l2.Value() {
		t.Error("shared a resource while another was idle")
	}
	// Both resources are in use: share rather than create.
	l3 := mustGet(t, p, "a")
	if made, _ := r.counts(); made != 2 {
		t.Errorf("made %d, want 2", made)
	}
	l.Release()
	l2.Release()
	l3.Release()
}

func TestResourcePool_DiscardShared(t *testing.T) {
	var r recorder
	p := newPool(&r)
	p.MaxResourcesPerKey = 1
	p.MaxConcurrentLeases = 2
	l := mustGet(t, p, "a")
	l2 := mustGet(t, p, "a")
	v := l.Value()

	l.Discard()
	if v.closed.Load() {
		t.Fatal("Discard closed a resource still leased")
	}
	c := getAsync(context.Background(), p, "a")
	expectBlocked(t, c) // must not get the discarded resource
	l2.Release()
	if !v.closed.Load() {
		t.Error("discarded resource not closed after its last lease")
	}
	l3 := receive(t, c)
	if l3.Value() == v {
		t.Error("got the discarded resource")
	}
	l3.Release()
}

func TestResourcePool_NewError(t *testing.T) {
	r := recorder{newErr: errors.New("boom")}
	p := newPool(&r)
	p.MaxResourcesPerKey = 1
	if _, err := p.Get(context.Background(), "a"); !errors.Is(err, r.newErr) {
		t.Fatalf("err = %v, want %v", err, r.newErr)
	}
	r.mu.Lock()
	r.newErr = nil
	r.mu.Unlock()
	mustGet(t, p, "a").Release() // the failed New must free its slot
}

func TestResourcePool_NewErrorWakesWaiter(t *testing.T) {
	var r recorder
	block := make(chan struct{})
	p := &sync_.ResourcePool[string, *res]{
		MaxResourcesPerKey: 1,
		Close:              r.Close,
		New: func(ctx context.Context, key string) (*res, error) {
			<-block
			if len(r.made) == 0 {
				r.made = append(r.made, nil)
				return nil, errors.New("first fails")
			}
			return &res{key: key}, nil
		},
	}
	errc := make(chan error, 1)
	go func() {
		_, err := p.Get(context.Background(), "a")
		errc <- err
	}()
	waitFor(t, func() bool { return p.Stats()["a"].Total == 1 })
	c := getAsync(context.Background(), p, "a")
	waitFor(t, func() bool { return p.Stats()["a"].Waiting == 1 })
	close(block)
	if err := receive(t, errc); err == nil {
		t.Fatal("first Get succeeded, want error")
	}
	receive(t, c).Release()
}

func TestResourcePool_Prefill(t *testing.T) {
	var r recorder
	p := newPool(&r)
	p.MaxResourcesPerKey = 2
	if err := p.Prefill(context.Background(), "a", 3); err == nil {
		t.Error("Prefill over MaxResourcesPerKey succeeded, want error")
	}
	if s := p.Stats()["a"]; s.Idle != 2 {
		t.Errorf("Stats = %+v, want 2 idle", s)
	}
}

func TestResourcePool_Do(t *testing.T) {
	var r recorder
	p := newPool(&r)

	var got *res
	if err := p.Do(context.Background(), func(ctx context.Context, v *res) error {
		got = v
		return nil
	}, "a"); err != nil {
		t.Fatal(err)
	}
	if s := p.Stats()["a"]; s.Idle != 1 {
		t.Errorf("Stats = %+v, want resource returned idle", s)
	}

	errApp := errors.New("app")
	if err := p.Do(context.Background(), func(ctx context.Context, v *res) error {
		return errApp
	}, "a"); err != errApp {
		t.Errorf("err = %v, want %v", err, errApp)
	}
	if got.closed.Load() {
		t.Error("an ordinary error discarded the resource")
	}

	errBad := fmt.Errorf("device lost: %w", sync_.ErrBadResource)
	if err := p.Do(context.Background(), func(ctx context.Context, v *res) error {
		return errBad
	}, "a"); err != errBad {
		t.Errorf("err = %v, want %v", err, errBad)
	}
	if !got.closed.Load() {
		t.Error("ErrBadResource did not discard the resource")
	}
}

func TestResourcePool_DoPanic(t *testing.T) {
	var r recorder
	p := newPool(&r)
	var got *res
	func() {
		defer func() {
			if e := recover(); e != "boom" {
				t.Errorf("recover() = %v, want boom", e)
			}
		}()
		_ = p.Do(context.Background(), func(ctx context.Context, v *res) error {
			got = v
			panic("boom")
		}, "a")
	}()
	if !got.closed.Load() {
		t.Error("panic in f did not discard the resource")
	}
	if s := p.Stats()["a"]; s.Total != 0 {
		t.Errorf("Stats = %+v, want empty", s)
	}
}

func TestResourcePool_Shutdown(t *testing.T) {
	var r recorder
	p := newPool(&r)
	p.MaxResourcesPerKey = 2
	if err := p.Prefill(context.Background(), "a", 1); err != nil {
		t.Fatal(err)
	}
	l := mustGet(t, p, "a")
	l2 := mustGet(t, p, "a")

	waiter := make(chan error, 1)
	go func() {
		_, err := p.Get(context.Background(), "a")
		waiter <- err
	}()
	waitFor(t, func() bool { return p.Stats()["a"].Waiting == 1 })

	done := make(chan error, 1)
	go func() { done <- p.Shutdown(context.Background()) }()

	if err := receive(t, waiter); !errors.Is(err, sync_.ErrResourcePoolClosed) {
		t.Errorf("waiter err = %v, want ErrResourcePoolClosed", err)
	}
	if _, err := p.Get(context.Background(), "a"); !errors.Is(err, sync_.ErrResourcePoolClosed) {
		t.Errorf("Get after Shutdown: err = %v, want ErrResourcePoolClosed", err)
	}
	if err := p.Prefill(context.Background(), "a", 1); !errors.Is(err, sync_.ErrResourcePoolClosed) {
		t.Errorf("Prefill after Shutdown: err = %v, want ErrResourcePoolClosed", err)
	}
	expectBlocked(t, done) // leases outstanding
	l.Release()
	l2.Release()
	if err := receive(t, done); err != nil {
		t.Errorf("Shutdown: %v", err)
	}
	if made, closed := r.counts(); made != 2 || closed != 2 {
		t.Errorf("made %d, closed %d; want 2, 2", made, closed)
	}
	if err := p.Shutdown(context.Background()); err != nil {
		t.Errorf("second Shutdown: %v", err)
	}
}

func TestResourcePool_ShutdownContext(t *testing.T) {
	var r recorder
	p := newPool(&r)
	l := mustGet(t, p, "a")
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	if err := p.Shutdown(ctx); !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("Shutdown: err = %v, want DeadlineExceeded", err)
	}
	l.Release()
	if !l.Value().closed.Load() {
		t.Error("resource released after Shutdown not closed")
	}
}

func TestResourcePool_ShutdownCloseError(t *testing.T) {
	errClose := errors.New("close")
	p := &sync_.ResourcePool[string, int]{
		New:   func(ctx context.Context, key string) (int, error) { return 0, nil },
		Close: func(key string, v int) error { return errClose },
	}
	if err := p.Prefill(context.Background(), "a", 1); err != nil {
		t.Fatal(err)
	}
	if err := p.Shutdown(context.Background()); !errors.Is(err, errClose) {
		t.Errorf("Shutdown: err = %v, want %v", err, errClose)
	}
}

func TestResourcePool_Concurrent(t *testing.T) {
	var r recorder
	p := newPool(&r)
	const maxPerKey, maxLeases = 2, 2
	p.MaxResourcesPerKey = maxPerKey
	p.MaxConcurrentLeases = maxLeases
	p.MaxIdleResourcesPerKey = 1
	p.IdleResourceTimeout = time.Millisecond
	keys := []string{"a", "b", "c"}

	var wg sync.WaitGroup
	for i := range 64 {
		wg.Go(func() {
			for j := range 50 {
				ks := keys[(i+j)%len(keys):]
				err := p.Do(context.Background(), func(ctx context.Context, v *res) error {
					if v.closed.Load() {
						t.Error("got a closed resource")
					}
					if n := v.users.Add(1); n > maxLeases {
						t.Errorf("%d concurrent leases, want <= %d", n, maxLeases)
					}
					defer v.users.Add(-1)
					if j%7 == 0 {
						return sync_.ErrBadResource
					}
					return nil
				}, ks...)
				if err != nil && !errors.Is(err, sync_.ErrBadResource) {
					t.Error(err)
				}
				for k, s := range p.Stats() {
					if s.Total > maxPerKey {
						t.Errorf("key %s has %d resources, want <= %d", k, s.Total, maxPerKey)
					}
				}
			}
		})
	}
	wg.Wait()
	if err := p.Shutdown(context.Background()); err != nil {
		t.Fatal(err)
	}
	if made, closed := r.counts(); made != closed {
		t.Errorf("made %d, closed %d; want all closed", made, closed)
	}
}

func ExampleResourcePool() {
	// A pool of at most 2 resources per GPU.
	p := &sync_.ResourcePool[int, string]{
		New: func(ctx context.Context, gpu int) (string, error) {
			return fmt.Sprintf("model@gpu%d", gpu), nil
		},
		MaxResourcesPerKey: 2,
	}
	defer p.Shutdown(context.Background())

	// Stateless: borrow a resource from GPU 0 or 1 for one call.
	_ = p.Do(context.Background(), func(ctx context.Context, m string) error {
		fmt.Println("infer on", m)
		return nil
	}, 0, 1)

	// Stateful: hold a resource across calls.
	l, err := p.Get(context.Background(), 1)
	if err != nil {
		panic(err)
	}
	defer l.Release()
	fmt.Println("session on", l.Value())

	// Output:
	// infer on model@gpu0
	// session on model@gpu1
}

func TestResourcePool_MaxIdleResources(t *testing.T) {
	var r recorder
	p := newPool(&r)
	p.MaxIdleResources = 2
	la, lb, lc := mustGet(t, p, "a"), mustGet(t, p, "b"), mustGet(t, p, "c")
	va := la.Value()
	la.Release()
	lb.Release()
	lc.Release() // evicts a, the least recently used idle resource of any key
	if !va.closed.Load() {
		t.Error("least recently used idle resource not closed")
	}
	stats := p.Stats()
	if _, ok := stats["a"]; ok || stats["b"].Idle != 1 || stats["c"].Idle != 1 {
		t.Errorf("Stats = %+v, want b and c idle", stats)
	}
}

func TestResourcePool_CloseIdleResources(t *testing.T) {
	var r recorder
	p := newPool(&r)
	if err := p.Prefill(context.Background(), "a", 2); err != nil {
		t.Fatal(err)
	}
	l := mustGet(t, p, "a")
	p.CloseIdleResources()
	if s := p.Stats()["a"]; s.Idle != 0 || s.Total != 1 {
		t.Errorf("Stats = %+v, want only the leased resource", s)
	}
	l.Release() // no Get since CloseIdleResources: closed, not kept idle
	if s := p.Stats()["a"]; s.Total != 0 {
		t.Errorf("Stats = %+v, want the released resource closed", s)
	}
	mustGet(t, p, "a").Release() // Get makes the pool keep idle resources again
	if s := p.Stats()["a"]; s.Idle != 1 {
		t.Errorf("Stats = %+v, want the released resource kept idle", s)
	}
}

// BenchmarkResourcePool_MaxIdleResources releases resources of many keys over
// MaxIdleResources, so that each Release evicts the least recently used idle one.
func BenchmarkResourcePool_MaxIdleResources(b *testing.B) {
	for _, keys := range []int{10, 1000, 100000} {
		b.Run(fmt.Sprintf("keys=%d", keys), func(b *testing.B) {
			p := &sync_.ResourcePool[int, int]{
				New:              func(ctx context.Context, key int) (int, error) { return key, nil },
				MaxIdleResources: keys / 2,
			}
			ctx := context.Background()
			for k := range keys {
				if err := p.Prefill(ctx, k, 1); err != nil {
					b.Fatal(err)
				}
			}
			b.ResetTimer()
			for i := range b.N {
				l, err := p.Get(ctx, i%keys)
				if err != nil {
					b.Fatal(err)
				}
				l.Release()
			}
		})
	}
}

func TestResourcePool_NewPanic(t *testing.T) {
	panicking := true
	p := &sync_.ResourcePool[string, int]{
		New: func(ctx context.Context, key string) (int, error) {
			if panicking {
				panic("boom")
			}
			return 1, nil
		},
		MaxResourcesPerKey: 1,
	}
	func() {
		defer func() {
			if e := recover(); e != "boom" {
				t.Errorf("recover() = %v, want boom", e)
			}
		}()
		_, _ = p.Get(context.Background(), "a")
	}()
	panicking = false
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	l, err := p.Get(ctx, "a") // the panicking New must free its slot
	if err != nil {
		t.Fatal(err)
	}
	l.Release()
	if err := p.Shutdown(ctx); err != nil {
		t.Errorf("Shutdown: %v", err)
	}
}

// BenchmarkResourcePool_Waiters releases resources of a key while many callers
// wait for resources of other keys.
func BenchmarkResourcePool_Waiters(b *testing.B) {
	for _, waiters := range []int{0, 100, 10000} {
		b.Run(fmt.Sprintf("waiters=%d", waiters), func(b *testing.B) {
			p := &sync_.ResourcePool[int, int]{
				New:                func(ctx context.Context, key int) (int, error) { return key, nil },
				MaxResourcesPerKey: 1,
			}
			ctx, cancel := context.WithCancel(context.Background())
			var wg sync.WaitGroup
			for k := 1; k <= waiters; k++ {
				l, err := p.Get(ctx, k)
				if err != nil {
					b.Fatal(err)
				}
				defer l.Release()
				wg.Go(func() { _, _ = p.Get(ctx, k) })
			}
			for i := 1; i <= waiters; i++ {
				for p.Stats()[i].Waiting < 1 {
					time.Sleep(time.Microsecond)
				}
			}
			b.ResetTimer()
			for range b.N {
				l, err := p.Get(ctx, 0)
				if err != nil {
					b.Fatal(err)
				}
				l.Release()
			}
			b.StopTimer()
			cancel()
			wg.Wait()
		})
	}
}

func TestResourcePool_WaitersFIFOAcrossKeys(t *testing.T) {
	var r recorder
	p := newPool(&r)
	p.MaxResourcesPerKey = 1
	la, lb := mustGet(t, p, "a"), mustGet(t, p, "b")

	first := getAsync(context.Background(), p, "a", "b")
	waitFor(t, func() bool { return p.Stats()["b"].Waiting == 1 })
	second := getAsync(context.Background(), p, "b")
	waitFor(t, func() bool { return p.Stats()["b"].Waiting == 2 })

	lb.Release()
	l := receive(t, first) // queued first for b
	expectBlocked(t, second)
	if s := p.Stats(); s["a"].Waiting != 0 || s["b"].Waiting != 1 {
		t.Errorf("Stats = %+v, want the first waiter dequeued from a and b", s)
	}
	l.Release()
	receive(t, second).Release()
	la.Release()
}

func TestResourcePool_CloseIdleResourcesServesWaiter(t *testing.T) {
	var r recorder
	p := newPool(&r)
	p.MaxResourcesPerKey = 1
	l := mustGet(t, p, "a")
	v := l.Value()
	c := getAsync(context.Background(), p, "a")
	waitFor(t, func() bool { return p.Stats()["a"].Waiting == 1 })
	p.CloseIdleResources()
	l.Release()
	l2 := receive(t, c)
	if l2.Value() != v {
		t.Error("waiter got a new resource, want the released one")
	}
	l2.Release()
}

// A caller waiting while the only resource is being made can share it once made.
func TestResourcePool_WaiterSharesResourceJustMade(t *testing.T) {
	made := make(chan struct{})
	p := &sync_.ResourcePool[string, *res]{
		New: func(ctx context.Context, key string) (*res, error) {
			<-made
			return &res{key: key}, nil
		},
		MaxResourcesPerKey:  1,
		MaxConcurrentLeases: 2,
	}
	first := getAsync(context.Background(), p, "a")
	waitFor(t, func() bool { return p.Stats()["a"].Total == 1 })
	second := getAsync(context.Background(), p, "a")
	waitFor(t, func() bool { return p.Stats()["a"].Waiting == 1 })
	close(made)
	l := receive(t, first)
	defer l.Release()
	l2 := receive(t, second) // l is still held
	defer l2.Release()
	if l.Value() != l2.Value() {
		t.Error("leases not shared")
	}
}

// A key of a dynamic type that is not comparable panics, without breaking the pool.
func TestResourcePool_UnhashableKey(t *testing.T) {
	p := &sync_.ResourcePool[any, int]{
		New: func(ctx context.Context, key any) (int, error) { return 1, nil },
	}
	func() {
		defer func() {
			if recover() == nil {
				t.Error("Get with an unhashable key did not panic")
			}
		}()
		_, _ = p.Get(context.Background(), []int{1})
	}()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	done := make(chan error, 1)
	go func() {
		l, err := p.Get(ctx, "ok")
		if err == nil {
			l.Release()
		}
		done <- err
	}()
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-ctx.Done():
		t.Fatal("pool unusable after a panicking Get")
	}
}

// The keys passed to Get belong to the caller, who may reuse them while Get waits.
func TestResourcePool_KeysNotRetained(t *testing.T) {
	var r recorder
	p := newPool(&r)
	p.MaxResourcesPerKey = 1
	l := mustGet(t, p, "a")
	keys := []string{"a"}
	ctx, cancel := context.WithCancel(context.Background())
	errc := make(chan error, 1)
	go func() {
		_, err := p.Get(ctx, keys...)
		errc <- err
	}()
	waitFor(t, func() bool { return p.Stats()["a"].Waiting == 1 })
	keys[0] = "b"
	cancel()
	if err := receive(t, errc); !errors.Is(err, context.Canceled) {
		t.Errorf("err = %v, want Canceled", err)
	}
	l.Release()
}

func TestResourcePool_ShutdownReportsCloseErrorToEachCall(t *testing.T) {
	errClose := errors.New("close")
	p := &sync_.ResourcePool[string, int]{
		New:   func(ctx context.Context, key string) (int, error) { return 0, nil },
		Close: func(key string, v int) error { return errClose },
	}
	l, err := p.Get(context.Background(), "a")
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := p.Shutdown(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("Shutdown: err = %v, want Canceled", err)
	}
	l.Release()
	if err := p.Shutdown(context.Background()); !errors.Is(err, errClose) {
		t.Errorf("second Shutdown: err = %v, want %v", err, errClose)
	}
}

func TestResourcePool_ClosePanicStillFreesResource(t *testing.T) {
	p := &sync_.ResourcePool[string, int]{
		New:   func(ctx context.Context, key string) (int, error) { return 0, nil },
		Close: func(key string, v int) error { panic("boom") },
	}
	l, err := p.Get(context.Background(), "a")
	if err != nil {
		t.Fatal(err)
	}
	func() {
		defer func() { _ = recover() }()
		l.Discard()
	}()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := p.Shutdown(ctx); err != nil {
		t.Errorf("Shutdown: %v", err)
	}
}

func TestResourcePool_DoNilFunc(t *testing.T) {
	var r recorder
	p := newPool(&r)
	defer func() {
		if recover() == nil {
			t.Error("Do with a nil function did not panic")
		}
		if made, _ := r.counts(); made != 0 {
			t.Errorf("made %d, want 0", made)
		}
	}()
	_ = p.Do(context.Background(), nil, "a")
}

// A key not equal to itself, such as NaN, could never be found again in the pool.
func TestResourcePool_KeyNotEqualToItself(t *testing.T) {
	type pair struct{ a, b float64 }
	check := func(t *testing.T, get func(ctx context.Context) error, prefill func(ctx context.Context) error, usable func(ctx context.Context) error) {
		t.Helper()
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		call := func(name string, f func(ctx context.Context) error) error {
			errc := make(chan error, 1)
			go func() {
				defer func() {
					if e := recover(); e != nil {
						errc <- fmt.Errorf("%s panicked: %v", name, e)
					}
				}()
				errc <- f(ctx)
			}()
			select {
			case err := <-errc:
				return err
			case <-ctx.Done():
				t.Fatalf("%s blocked: the pool is stuck", name)
				return nil
			}
		}
		if err := call("Get", get); !errors.Is(err, sync_.ErrResourcePoolInvalidKey) {
			t.Errorf("Get: err = %v, want ErrResourcePoolInvalidKey", err)
		}
		if err := call("Prefill", prefill); !errors.Is(err, sync_.ErrResourcePoolInvalidKey) {
			t.Errorf("Prefill: err = %v, want ErrResourcePoolInvalidKey", err)
		}
		if err := call("Get of a valid key", usable); err != nil {
			t.Errorf("Get of a valid key: %v", err)
		}
	}

	t.Run("float64", func(t *testing.T) {
		p := &sync_.ResourcePool[float64, int]{
			New: func(ctx context.Context, key float64) (int, error) { return 1, nil },
		}
		check(t,
			func(ctx context.Context) error { _, err := p.Get(ctx, 1, math.NaN()); return err },
			func(ctx context.Context) error { return p.Prefill(ctx, math.NaN(), 1) },
			func(ctx context.Context) error {
				l, err := p.Get(ctx, 1)
				if err == nil {
					l.Release()
				}
				return err
			})
	})
	t.Run("any", func(t *testing.T) {
		p := &sync_.ResourcePool[any, int]{
			New: func(ctx context.Context, key any) (int, error) { return 1, nil },
		}
		check(t,
			func(ctx context.Context) error { _, err := p.Get(ctx, math.NaN()); return err },
			func(ctx context.Context) error { return p.Prefill(ctx, math.NaN(), 1) },
			func(ctx context.Context) error {
				l, err := p.Get(ctx, "a")
				if err == nil {
					l.Release()
				}
				return err
			})
	})
	t.Run("struct", func(t *testing.T) {
		p := &sync_.ResourcePool[pair, int]{
			New: func(ctx context.Context, key pair) (int, error) { return 1, nil },
		}
		check(t,
			func(ctx context.Context) error { _, err := p.Get(ctx, pair{1, math.NaN()}); return err },
			func(ctx context.Context) error { return p.Prefill(ctx, pair{1, math.NaN()}, 1) },
			func(ctx context.Context) error {
				l, err := p.Get(ctx, pair{1, 2})
				if err == nil {
					l.Release()
				}
				return err
			})
	})
}
