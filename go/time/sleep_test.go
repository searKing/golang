package time_test

import (
	"errors"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	time_ "github.com/searKing/golang/go/time"
)

func benchmark(b *testing.B, bench func(n int)) {
	// Create equal number of garbage timers on each P before starting
	// the benchmark.
	var wg sync.WaitGroup
	garbageAll := make([][]*time_.Timer, runtime.GOMAXPROCS(0))
	for i := range garbageAll {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			garbage := make([]*time_.Timer, 1<<15)
			for j := range garbage {
				garbage[j] = time_.AfterFunc(time.Hour, nil)
			}
			garbageAll[i] = garbage
		}(i)
	}
	wg.Wait()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			bench(1000)
		}
	})
	b.StopTimer()

	for _, garbage := range garbageAll {
		for _, t := range garbage {
			t.Stop()
		}
	}
}
func BenchmarkReset(b *testing.B) {
	benchmark(b, func(n int) {
		t := time_.NewTimer(time.Hour)
		for i := 0; i < n; i++ {
			t.Reset(time.Hour)
		}
		t.Stop()
	})
}

func testReset(d time.Duration) error {
	t0 := time.NewTimer(2 * d)
	time.Sleep(d)
	if !t0.Reset(3 * d) {
		return errors.New("resetting unfired timer returned false")
	}
	time.Sleep(2 * d)
	select {
	case <-t0.C:
		return errors.New("timer fired early")
	default:
	}
	time.Sleep(2 * d)
	select {
	case <-t0.C:
	default:
		return errors.New("reset timer did not fire")
	}

	if t0.Reset(50 * time.Millisecond) {
		return errors.New("resetting expired timer returned true")
	}
	return nil
}

func TestReset(t *testing.T) {
	// We try to run this test with increasingly larger multiples
	// until one works so slow, loaded hardware isn't as flaky,
	// but without slowing down fast machines unnecessarily.
	const unit = 25 * time.Millisecond
	tries := []time.Duration{
		1 * unit,
		3 * unit,
		7 * unit,
		15 * unit,
	}
	var err error
	for _, d := range tries {
		err = testReset(d)
		if err == nil {
			t.Logf("passed using duration %v", d)
			return
		}
	}
	t.Error(err)
}

func checkZeroPanicString(t *testing.T) {
	e := recover()
	s, _ := e.(string)
	if want := "called on uninitialized Timer"; !strings.Contains(s, want) {
		t.Errorf("panic = %v; want substring %q", e, want)
	}
}

func TestZeroTimerResetPanics(t *testing.T) {
	defer checkZeroPanicString(t)
	var tr time_.Timer
	tr.Reset(1)
}

func TestTimer_Reset(t *testing.T) {
	tr := time_.NewTimer(25 * time.Millisecond)
	defer func() {
		tr.Stop()
	}()
	<-tr.C
}
