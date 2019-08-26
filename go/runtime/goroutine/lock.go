// Defensive debug-only utility to track that functions run on the
// goroutine that they're supposed to.

package goroutine

import (
	"errors"
	"github.com/searKing/golang/go/error/must"
	"os"
)

var DebugGoroutines = os.Getenv("DEBUG_GOROUTINES") == "1"

type Lock uint64

// NewLock returns a GoRoutine Lock
func NewLock() Lock {
	if !DebugGoroutines {
		return 0
	}
	return Lock(ID())
}

// Check if caller's goroutine is locked
func (g Lock) Check() error {
	if !DebugGoroutines {
		return nil
	}
	if ID() != uint64(g) {
		return errors.New("running on the wrong goroutine")
	}
	return nil
}

func (g Lock) MustCheck() {
	must.Must(g.Check())
}

// Check whether caller's goroutine escape lock
func (g Lock) CheckNotOn() error {
	if !DebugGoroutines {
		return nil
	}
	if ID() == uint64(g) {
		return errors.New("running on the wrong goroutine")
	}
	return nil
}

func (g Lock) MustCheckNotOn() {
	must.Must(g.CheckNotOn())
}
