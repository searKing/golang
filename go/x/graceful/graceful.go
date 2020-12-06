package graceful

import (
	"context"
	"sync"

	context_ "github.com/searKing/golang/go/context"
	"github.com/searKing/golang/go/runtime"
)

// StarFunc is the type of the function invoked by Graceful to start the server
type StartFunc func(context.Context) error

// ShutdownFunc is the type of the function invoked by Graceful to shutdown the server
type ShutdownFunc func(context.Context) error

// Graceful sets up graceful handling of SIGINT and SIGTERM, typically for an HTTP server.
// When signal is trapped, the shutdown handler will be invoked with a context.
func Graceful(ctx context.Context, start StartFunc, shutdown ShutdownFunc) (err error) {
	defer runtime.LogPanic.Recover()
	if start == nil {
		start = func(ctx context.Context) error { return nil }
	}
	if shutdown == nil {
		shutdown = func(ctx context.Context) error { return nil }
	}

	ctx = context_.WithShutdownSignal(ctx)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case <-ctx.Done():
			err = shutdown(ctx)
		}
	}()

	// Start the server
	if err := start(ctx); err != nil {
		return err
	}

	wg.Wait()
	return err
}
