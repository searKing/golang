package goroutine

import (
	"fmt"
	"strings"
	"testing"
)

func TestLock(t *testing.T) {
	oldDebug := DebugGoroutines
	DebugGoroutines = true
	defer func() { DebugGoroutines = oldDebug }()

	g := NewLock()
	g.MustCheck()

	sawPanic := make(chan interface{})
	go func() {
		defer func() { sawPanic <- recover() }()
		g.MustCheck() // should panic
	}()
	e := <-sawPanic
	if e == nil {
		t.Fatal("did not see panic from check in other goroutine")
	}
	if !strings.Contains(fmt.Sprint(e), "wrong goroutine") {
		t.Errorf("expected on see panic about running on the wrong goroutine; got %v", e)
	}
}
