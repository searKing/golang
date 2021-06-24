package context

import (
	"context"
	"os"
	"os/signal"
	"sync"

	os_ "github.com/searKing/golang/go/os"
)

var onlyOneSignalHandler = make(chan struct{})

// WithShutdownSignal registered for signals. A context.Context is returned.
// If no signals are provided, incoming os_.ShutdownSignals signals will be relayed.
// Otherwise, just the provided signals will.
// which is done on one of these incoming signals. If a second signal is caught, the program
// is terminated with exit code 1.
// Only one of Signal should be called, and only can be called once.
func WithShutdownSignal(parent context.Context, sig ...os.Signal) context.Context {
	if len(sig) == 0 {
		sig = os_.ShutdownSignals
	}
	close(onlyOneSignalHandler) // panics when called twice

	var shutdownSignalC chan os.Signal
	shutdownSignalC = make(chan os.Signal, 2)

	ctx, cancel := context.WithCancel(parent)
	signal.Notify(shutdownSignalC, sig...)
	go func() {
		defer signal.Stop(shutdownSignalC)
		select {
		case <-shutdownSignalC:
			cancel()
		case <-ctx.Done():
			return
		}

		select {
		case <-shutdownSignalC:
			os.Exit(1) // second signal. Exit directly.
		case <-ctx.Done():
		}
	}()

	return ctx
}

// WithSignal registered for signals. A context.Context is returned.
// which is done on one of these incoming signals.
// signals can be stopped by stopSignal, context will never be Done() if stopped.
// WithSignal returns a copy of the parent context registered for signals.
// If the parent's context is already done than d, WithSignal(parent, sig...) is semantically
// equivalent to parent. The returned context's Done channel is closed when one of these
// incoming signals, when the returned cancel function is called, or when the parent context's
// Done channel is closed, whichever happens first.
// Canceling this context releases resources associated with it, so code should
// call cancel as soon as the operations running in this Context complete.
// Deprecated: Use signal.NotifyContext instead since go1.16.
func WithSignal(parent context.Context, sig ...os.Signal) (ctx context.Context, cancel context.CancelFunc) {
	var c chan os.Signal
	c = make(chan os.Signal, 1)

	ctx, cancel_ := context.WithCancel(parent)

	var once sync.Once
	cancel = func() {
		once.Do(func() {
			// https://groups.google.com/g/golang-nuts/c/pZwdYRGxCIk/m/qpbHxRRPJdUJ
			// https://groups.google.com/g/golang-nuts/c/bfuwmhbZHYw/m/PNnKk4hGytgJ
			// The close function should only be used on a channel used to send data,
			// after all data has been sent. It does not make sense to send more data
			// after the channel has been closed; if all data has not been sent, the
			// channel should not be closed. Similarly, it does not make sense to
			// close the channel twice.
			//
			// Note that it is only necessary to close a channel if the receiver is
			// looking for a close. Closing the channel is a control signal on the
			// channel indicating that no more data follows.
			// BUT
			// Channels persist as long as one reference to them persists,
			// since without difficult (and in general impossible) analysis
			// it's hard to prove a second reference won't arise.
			defer close(c)

			signal.Stop(c)
			cancel_()
		})
	}

	signal.Notify(c, sig...)
	go func() {
		select {
		case <-c:
		case <-ctx.Done():
		}
		cancel()
	}()

	return ctx, cancel
}
