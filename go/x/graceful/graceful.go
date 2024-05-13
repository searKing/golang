package graceful

import (
	"context"
	"fmt"
	"sync"

	"github.com/searKing/golang/go/errors"
	"github.com/searKing/golang/go/runtime"
)

type Handler struct {
	Name         string       // server name
	StartFunc    StartFunc    // func to start the server
	ShutdownFunc ShutdownFunc // func to shutdown the server
}

// StartFunc is the type of the function invoked by Graceful to start the server
type StartFunc func(context.Context) error

// ShutdownFunc is the type of the function invoked by Graceful to shutdown the server
type ShutdownFunc func(context.Context) error

// Graceful sets up graceful handling of context done, typically for an HTTP server.
// When context is done, the shutdown handler will be invoked with a context.
// Example:
//
//	ctx is wrapped WithShutdownSignal(ctx)
//	When signal is trapped, the shutdown handler will be invoked with a context.
func Graceful(ctx context.Context, handlers ...Handler) (err error) {
	if len(handlers) == 0 {
		return nil
	}
	defer runtime.LogPanic.Recover()
	var wg sync.WaitGroup

	var mu sync.Mutex
	var errs []error

	for _, h := range handlers {
		start := h.StartFunc
		shutdown := h.ShutdownFunc
		if start == nil {
			start = func(ctx context.Context) error { return nil }
		}
		if shutdown == nil {
			shutdown = func(ctx context.Context) error { return nil }
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			select {
			case <-ctx.Done():
				err = shutdown(ctx)
				if err != nil {
					err = fmt.Errorf("graceful shutdown %s: %w", h.Name, err)
					mu.Lock()
					defer mu.Unlock()
					errs = append(errs, err)
				}
			}
		}()

		// Start the server
		if err := start(ctx); err != nil {
			return fmt.Errorf("graceful start %s: %w", h.Name, err)
		}
	}

	wg.Wait()
	return errors.Multi(errs...)
}
