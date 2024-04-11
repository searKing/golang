// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sync_test

import (
	"fmt"
	"path"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"sort"
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	sync_ "github.com/searKing/golang/go/exp/sync"
	runtime_ "github.com/searKing/golang/go/runtime"
)

func fixGC() {
	// TODO: Fix #45315 and remove this extra call.
	//
	// Unfortunately, it's possible for the sweep termination condition
	// to flap, so with just one runtime.GC call, a freed object could be
	// missed, leading this test to fail. A second call reduces the chance
	// of this happening to zero, because sweeping actually has to finish
	// to move on to the next GC, during which nothing will happen.
	//
	// See https://github.com/golang/go/issues/46500 and
	// https://github.com/golang/go/issues/45315 for more details.
	runtime.GOMAXPROCS(1)
}
func caller() string {
	function, file, line := runtime_.GetCallerFuncFileLine(3)
	return fmt.Sprintf("%s() %s:%d", path.Base(function), filepath.Base(file), line)
}

func testFixedPoolLenAndCap[E any](t *testing.T, p *sync_.FixedPool[E], l, c int) {
	gotLen := p.Len()
	gotCap := p.Cap()
	if (gotLen != l && c >= 0) || (gotCap != c && c >= 0) {
		t.Fatalf("%s, got %d|%d; want %d|%d", caller(), gotLen, gotCap, l, c)
	}
}

func TestNewFixedPool(t *testing.T) {
	fixGC()

	// disable GC so we can control when it happens.
	defer debug.SetGCPercent(debug.SetGCPercent(-1))
	var i int
	f := func() string {
		defer func() { i++ }()
		return strconv.Itoa(i)
	}
	var p = sync_.NewFixedPool[string](f, 2)
	testFixedPoolLenAndCap(t, p, 2, 2)
	if g := p.TryGet(); g == nil || g.Value != "0" {
		t.Fatalf("got %#v; want 0", g)
	}
	testFixedPoolLenAndCap(t, p, 1, 2)
	p.Emplace("a")
	testFixedPoolLenAndCap(t, p, 2, 3)
	p.Emplace("b")
	testFixedPoolLenAndCap(t, p, 3, 4)
	if g := p.TryGet(); g == nil || g.Value != "1" {
		t.Fatalf("got %#v; want 1", g)
	}
	if g := p.Get(); g == nil || g.Value != "a" {
		t.Fatalf("got %#v; want a", g)
	}
	testFixedPoolLenAndCap(t, p, 1, 4)
	if g := p.Get(); g.Value != "b" {
		t.Fatalf("got %#v; want b", g)
	}
	testFixedPoolLenAndCap(t, p, 0, 4)
	if g := p.TryGet(); g != nil {
		t.Fatalf("got %#v; want nil", g)
	}
	// After one GC, the victim cache should keep them alive.
	runtime.GC()
	// drop all the items taken by Get and not be referenced by any
	testFixedPoolLenAndCap(t, p, 2, 2)
	// A second GC should drop the victim cache.
	runtime.GC()
	testFixedPoolLenAndCap(t, p, 2, 2)

	// Put in a large number of items, so they spill into
	// stealable space.
	n := 100
	for i := 0; i < n; i++ {
		p.Emplace("c")
		testFixedPoolLenAndCap(t, p, i+1+2, i+1+2)
	}
	testFixedPoolLenAndCap(t, p, 102, 102)
	for i := 0; i < n; i++ {
		if g := p.Get(); g == nil {
			t.Fatalf("got empty")
		}
	}
	testFixedPoolLenAndCap(t, p, 2, 102)
	if g := p.TryGet(); g == nil {
		t.Fatalf("got empty")
	}
	testFixedPoolLenAndCap(t, p, 1, 102)
	// After one GC, the victim cache should keep them alive.
	runtime.GC()
	// drop all the items taken by Get and not be referenced by any
	testFixedPoolLenAndCap(t, p, 3, 3)
	// A second GC should drop the victim cache.
	runtime.GC()
	testFixedPoolLenAndCap(t, p, 2, 2)
}
func TestNewCachePool(t *testing.T) {
	fixGC()
	// disable GC so we can control when it happens.
	defer debug.SetGCPercent(debug.SetGCPercent(-1))
	var p = sync_.NewCachedPool[string](nil)
	if p.TryGet() != nil {
		t.Fatal("expected empty")
	}
	p.Emplace("a")
	testFixedPoolLenAndCap(t, p, 1, 1)
	p.Emplace("b")
	testFixedPoolLenAndCap(t, p, 2, 2)
	if g := p.Get(); g.Value != "a" {
		t.Fatalf("got %#v; want a", g)
	}
	testFixedPoolLenAndCap(t, p, 1, 2)
	if g := p.Get(); g.Value != "b" {
		t.Fatalf("got %#v; want b", g)
	}
	testFixedPoolLenAndCap(t, p, 0, 2)
	if g := p.TryGet(); g != nil {
		t.Fatalf("got %#v; want nil", g)
	}
	testFixedPoolLenAndCap(t, p, 0, 2)

	// Put in a large number of items, so they spill into
	// stealable space.
	n := 100
	for i := 0; i < n; i++ {
		p.Emplace("c")
		testFixedPoolLenAndCap(t, p, i+1, i+1+2)
	}
	testFixedPoolLenAndCap(t, p, 100, 102)
	for i := 0; i < n; i++ {
		if g := p.Get(); g.Value != "c" {
			t.Fatalf("got %#v; want a", g)
		}
	}
	testFixedPoolLenAndCap(t, p, 0, 102)
	if g := p.TryGet(); g != nil {
		t.Fatalf("got %#v; want nil", g)
	}
	testFixedPoolLenAndCap(t, p, 0, 102)
}

func TestNewTempPool(t *testing.T) {
	fixGC()
	// disable GC so we can control when it happens.
	defer debug.SetGCPercent(debug.SetGCPercent(-1))
	var p = sync_.NewTempPool[string](nil)
	if p.TryGet() != nil {
		t.Fatal("expected empty")
	}
	testFixedPoolLenAndCap(t, p, 0, 0)
	p.Emplace("a")

	testFixedPoolLenAndCap(t, p, 1, 1)
	p.Emplace("b")
	testFixedPoolLenAndCap(t, p, 2, 2)
	if g := p.Get(); g.Value != "a" {
		t.Fatalf("got %#v; want a", g)
	}
	testFixedPoolLenAndCap(t, p, 1, 2)
	if g := p.Get(); g.Value != "b" {
		t.Fatalf("got %#v; want b", g)
	}
	testFixedPoolLenAndCap(t, p, 0, 2)

	// Put in a large number of items, so they spill into
	// stealable space.
	for i := 0; i < 100; i++ {
		p.Emplace("c")
		testFixedPoolLenAndCap(t, p, i+1, i+1+2)
	}
	testFixedPoolLenAndCap(t, p, 100, 102)
	// After one GC, the victim cache should keep them alive.
	runtime.GC()
	// drop all the items taken by Get and not be referenced by any
	testFixedPoolLenAndCap(t, p, 100, 100)
	if g := p.Get(); g.Value != "c" {
		t.Fatalf("got %#v; want c after GC", g)
	}
	testFixedPoolLenAndCap(t, p, 99, 100)
	// A second GC should drop the victim cache.
	runtime.GC()
	testFixedPoolLenAndCap(t, p, 0, 0)
	if g := p.TryGet(); g != nil {
		t.Fatalf("got %#v; want nil after second GC", g)
	}
	testFixedPoolLenAndCap(t, p, 0, 0)
}

func TestFixedPoolNilNew(t *testing.T) {
	fixGC()
	// disable GC so we can control when it happens.
	defer debug.SetGCPercent(debug.SetGCPercent(-1))
	const LEN = 10
	const CAP = 20
	var p = (&sync_.FixedPool[string]{
		New:             nil,
		MinResidentSize: 0,
		MaxResidentSize: LEN,
		MaxCapacity:     CAP,
	}).Init()

	testFixedPoolLenAndCap(t, p, 0, 0)

	if p.TryGet() != nil {
		t.Fatal("expected empty")
	}
	p.Emplace("a")
	testFixedPoolLenAndCap(t, p, 1, 1)
	p.Emplace("b")
	testFixedPoolLenAndCap(t, p, 2, 2)
	if g := p.Get(); g.Value != "a" {
		t.Fatalf("got %#v; want a", g)
	}
	testFixedPoolLenAndCap(t, p, 1, 2)
	if g := p.Get(); g.Value != "b" {
		t.Fatalf("got %#v; want b", g)
	}
	testFixedPoolLenAndCap(t, p, 0, 2)
	if g := p.TryGet(); g != nil {
		t.Fatalf("got %#v; want nil", g)
	}
	testFixedPoolLenAndCap(t, p, 0, 2)

	// Put in a large number of items, so they spill into
	// stealable space.
	for i := 0; i < 100; i++ {
		p.Emplace("c")
		testFixedPoolLenAndCap(t, p, i+1, i+1+2)
	}
	testFixedPoolLenAndCap(t, p, 100, 102)
	// After one GC, the victim cache should keep them alive.
	runtime.GC()
	testFixedPoolLenAndCap(t, p, 100, 100)
	if g := p.Get(); g.Value != "c" {
		t.Fatalf("got %#v; want c after GC", g)
	}
	testFixedPoolLenAndCap(t, p, 99, 100)
	// A second GC should drop the victim cache, try put into local first.
	runtime.GC()
	testFixedPoolLenAndCap(t, p, LEN, LEN)

	// drain keep-alive cache
	for i := 0; i < LEN; i++ {
		if g := p.Get(); g == nil || g.Value != "c" {
			t.Fatalf("#%d: got %#v; want c after GC", i, g)
		}
	}
	if g := p.TryGet(); g != nil {
		t.Fatalf("got %#v; want nil after second GC", g)
	}
	testFixedPoolLenAndCap(t, p, 0, LEN)
	// After one GC, the victim cache should keep them alive.
	// After one GC, the got object will be GC, as no reference
	runtime.GC()
	testFixedPoolLenAndCap(t, p, LEN, LEN)
	// A second GC should drop the victim cache, try put into local first.
	runtime.GC()
	testFixedPoolLenAndCap(t, p, LEN, LEN)
}

func TestFixedPoolNew(t *testing.T) {
	fixGC()
	// disable GC so we can control when it happens.
	defer debug.SetGCPercent(debug.SetGCPercent(-1))

	const MinLen = 2
	const LEN = 10
	const CAP = 20
	i := 0
	var p = (&sync_.FixedPool[int]{
		New: func() int {
			i++
			return i
		},
		MinResidentSize: MinLen,
		MaxResidentSize: LEN,
		MaxCapacity:     CAP,
	}).Init()
	testFixedPoolLenAndCap(t, p, MinLen, MinLen)

	if v := p.Get(); v.Value != 1 {
		t.Fatalf("got %v; want 1", v.Value)
	}
	if v := p.Get(); v.Value != 2 {
		t.Fatalf("got %v; want 2", v.Value)
	}

	p.Emplace(42)
	if v := p.Get(); v.Value != 42 {
		t.Fatalf("got %v; want 42", v)
	}

	if v := p.Get(); v.Value != 3 {
		t.Fatalf("got %v; want 3", v)
	}
}

func TestFixedPoolGCRetryPut(t *testing.T) {
	// disable GC so we can control when it happens.
	defer debug.SetGCPercent(debug.SetGCPercent(-1))
	const LEN = 1
	const CAP = 2
	var p = (&sync_.FixedPool[string]{
		New:             nil,
		MinResidentSize: 0,
		MaxResidentSize: LEN,
		MaxCapacity:     CAP,
	}).Init()

	testFixedPoolLenAndCap(t, p, 0, 0)

	if p.TryGet() != nil {
		t.Fatal("expected empty")
	}

	// Put in a large number of items, so they spill into
	// stealable space.
	var N = 4
	for i := 0; i < N; i++ {
		p.Emplace(strconv.Itoa(i))
		testFixedPoolLenAndCap(t, p, i+1, i+1)
	}
	testFixedPoolLenAndCap(t, p, N, N)
	// After one GC, the victim cache should keep them alive.
	runtime.GC()
	testFixedPoolLenAndCap(t, p, N, N)
	if g := p.Get(); g.Value != "0" {
		t.Fatalf("got %#v; want c after GC", g)
	}
	testFixedPoolLenAndCap(t, p, N-1, N)
	// A second GC should drop the victim cache, try put into local first.
	runtime.GC()
	testFixedPoolLenAndCap(t, p, LEN, LEN)

	// drain keep-alive cache
	for i := 1; i < LEN+1; i++ {
		if g := p.Get(); g == nil {
			t.Fatalf("#%d: got nil; want %q after GC", i, strconv.Itoa(i))
		}
	}
	testFixedPoolLenAndCap(t, p, 0, LEN)
	{
		if g := p.TryGet(); g != nil {
			t.Fatalf("got %#v; want nil after second GC", g.Value)
		}
	}
	testFixedPoolLenAndCap(t, p, 0, LEN)
	// After one GC, the victim cache should keep them alive.
	runtime.GC()
	testFixedPoolLenAndCap(t, p, LEN, LEN)
	// A second GC should drop the victim cache, try put into local first.
	runtime.GC()
	testFixedPoolLenAndCap(t, p, LEN, LEN)
	{
		if g := p.TryGet(); g == nil {
			t.Fatalf("got nil; want %q after GC", g.Value)
		}
	}
	testFixedPoolLenAndCap(t, p, 0, LEN)
}

func TestFixedPoolGCReFillLocal(t *testing.T) {
	// disable GC so we can control when it happens.
	defer debug.SetGCPercent(debug.SetGCPercent(-1))
	const LEN = 1
	const CAP = 2
	var p = (&sync_.FixedPool[string]{
		New:             nil,
		MinResidentSize: 0,
		MaxResidentSize: LEN,
		MaxCapacity:     CAP,
	}).Init()
	// Put in a large number of items, so they spill into
	// stealable space.
	for i := 0; i < CAP*2; i++ {
		p.Emplace(strconv.Itoa(i))
	}
	testFixedPoolLenAndCap(t, p, CAP*2, CAP*2)

	// drain all cache
	for i := 0; i < 2*CAP; i++ {
		g := p.Get()
		if i < LEN {
			if g == nil || g.Value != strconv.Itoa(i) {
				t.Fatalf("#%d: got %#v; want %q after GC", i, g, strconv.Itoa(i))
			}
		} else {
			if g == nil {
				t.Fatalf("#%d: got %#v; want %q after GC", i, g, strconv.Itoa(i))
			}
		}
		testFixedPoolLenAndCap(t, p, CAP*2-i-1, CAP*2)
	}
	testFixedPoolLenAndCap(t, p, 0, CAP*2)
	if g := p.TryGet(); g != nil {
		t.Fatalf("got %#v; want nil after second GC", g)
	}
	testFixedPoolLenAndCap(t, p, 0, CAP*2)

	// After one GC, the victim cache should keep them alive.
	runtime.GC()
	testFixedPoolLenAndCap(t, p, LEN, LEN)
	// A second GC should drop the victim cache, try put into local first.
	runtime.GC()
	testFixedPoolLenAndCap(t, p, LEN, LEN)

	// drain all cache
	for i := 0; i < LEN; i++ {
		g := p.Get()
		if g == nil {
			t.Fatalf("#%d: got nil; want not nil after GC", i)
		}
	}
	testFixedPoolLenAndCap(t, p, 0, LEN)
	if g := p.TryGet(); g != nil {
		t.Fatalf("got %#v; want nil after second GC", g)
	}
	testFixedPoolLenAndCap(t, p, 0, LEN)
}

// Test that Pool does not hold pointers to previously cached resources.
func TestFixedPoolGC(t *testing.T) {
	testFixedPool(t, true)
}

// Test that Pool releases resources on GC.
func TestFixedPoolRelease(t *testing.T) {
	testFixedPool(t, false)
}

func testFixedPool(t *testing.T, drain bool) {
	var p sync_.FixedPool[*string]
	const N = 100
loop:
	for try := 0; try < 3; try++ {
		if try == 1 && testing.Short() {
			testFixedPoolLenAndCap(t, &p, 0, 0)
			break
		}
		var fin, fin1 uint32
		for i := 0; i < N; i++ {
			v := new(string)
			runtime.SetFinalizer(v, func(vv *string) {
				atomic.AddUint32(&fin, 1)
			})
			p.Emplace(v)
		}
		if drain {
			for i := 0; i < N; i++ {
				p.Get()
			}
		}
		for i := 0; i < 5; i++ {
			runtime.GC()
			time.Sleep(time.Duration(i*100+10) * time.Millisecond)
			// 1 pointer can remain on stack or elsewhere
			if fin1 = atomic.LoadUint32(&fin); fin1 >= N-1 {
				continue loop
			}
		}
		t.Fatalf("only %v out of %v resources are finalized on try %v", fin1, N, try)
	}
}

func TestFixedPoolStress(t *testing.T) {
	const P = 10
	N := int(1e6)
	if testing.Short() {
		N /= 100
	}
	var p sync_.FixedPool[any]
	done := make(chan bool)
	for i := 0; i < P; i++ {
		go func() {
			var v any = 0
			for j := 0; j < N; j++ {
				if v == nil {
					v = 0
				}
				p.Emplace(v)
				e := p.Get()
				if e != nil && e.Value != 0 {
					t.Errorf("expect 0, got %v", v)
					break
				}
			}
			done <- true
		}()
	}
	for i := 0; i < P; i++ {
		<-done
	}
}

func BenchmarkFixedPool(b *testing.B) {
	var p sync_.FixedPool[int]
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			p.Emplace(1)
			p.Get()
		}
	})
}

func BenchmarkPoolOverflow(b *testing.B) {
	var p sync_.FixedPool[int]
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for b := 0; b < 100; b++ {
				p.Emplace(1)
			}
			for b := 0; b < 100; b++ {
				p.Get()
			}
		}
	})
}

// Simulate object starvation in order to force Ps to steal items
// from other Ps.
func BenchmarkPoolStarvation(b *testing.B) {
	var p sync_.FixedPool[int]
	count := 100
	// Reduce number of putted items by 33 %. It creates items starvation
	// that force P-local storage to steal items from other Ps.
	countStarved := count - int(float32(count)*0.33)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for b := 0; b < countStarved; b++ {
				p.Emplace(1)
			}
			for b := 0; b < count; b++ {
				p.Get()
			}
		}
	})
}

var globalSink any

func BenchmarkPoolSTW(b *testing.B) {
	// Take control of GC.
	defer debug.SetGCPercent(debug.SetGCPercent(-1))

	var mstats runtime.MemStats
	var pauses []uint64

	var p sync_.FixedPool[any]
	for i := 0; i < b.N; i++ {
		// Put a large number of items into a pool.
		const N = 100000
		var item any = 42
		for i := 0; i < N; i++ {
			p.Emplace(item)
		}
		// Do a GC.
		runtime.GC()
		// Record pause time.
		runtime.ReadMemStats(&mstats)
		pauses = append(pauses, mstats.PauseNs[(mstats.NumGC+255)%256])
	}

	// Get pause time stats.
	sort.Slice(pauses, func(i, j int) bool { return pauses[i] < pauses[j] })
	var total uint64
	for _, ns := range pauses {
		total += ns
	}
	// ns/op for this benchmark is average STW time.
	b.ReportMetric(float64(total)/float64(b.N), "ns/op")
	b.ReportMetric(float64(pauses[len(pauses)*95/100]), "p95-ns/STW")
	b.ReportMetric(float64(pauses[len(pauses)*50/100]), "p50-ns/STW")
}

func BenchmarkPoolExpensiveNew(b *testing.B) {
	// Populate a pool with items that are expensive to construct
	// to stress pool cleanup and subsequent reconstruction.

	// Create a ballast so the GC has a non-zero heap size and
	// runs at reasonable times.
	globalSink = make([]byte, 8<<20)
	defer func() { globalSink = nil }()

	// Create a pool that's "expensive" to fill.
	var p sync_.FixedPool[any]
	var nNew uint64
	p.New = func() any {
		atomic.AddUint64(&nNew, 1)
		time.Sleep(time.Millisecond)
		return 42
	}
	var mstats1, mstats2 runtime.MemStats
	runtime.ReadMemStats(&mstats1)
	b.RunParallel(func(pb *testing.PB) {
		// Simulate 100X the number of goroutines having items
		// checked out from the Pool simultaneously.
		items := make([]*sync_.FixedPoolElement[any], 100)
		var sink []byte
		for pb.Next() {
			// Stress the pool.
			for i := range items {
				items[i] = p.Get()
				// Simulate doing some work with this
				// item checked out.
				sink = make([]byte, 32<<10)
			}
			for i, v := range items {
				p.Put(v)
				items[i] = nil
			}
		}
		_ = sink
	})
	runtime.ReadMemStats(&mstats2)

	b.ReportMetric(float64(mstats2.NumGC-mstats1.NumGC)/float64(b.N), "GCs/op")
	b.ReportMetric(float64(nNew)/float64(b.N), "New/op")
}
