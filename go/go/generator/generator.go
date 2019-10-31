package generator

import "context"

// runtimeGenerator is an implement of Generator's behavior actually.
type runtimeGenerator struct {
	// fired func, as callback when supplierC is consumed successfully
	// arg for msg receiver
	// msg for msg to be delivered
	f   func(arg interface{}, msg interface{})
	arg interface{}

	ctx    context.Context
	cancel context.CancelFunc

	// data src
	supplierC <-chan interface{}
	// data dst
	consumerC chan interface{}
}

func (g *runtimeGenerator) start() {
	go func() {
		for {
			select {
			case <-g.ctx.Done():
				return
			case s, ok := <-g.supplierC:
				if !ok {
					return
				}
				g.f(g.arg, s)
			}
		}
	}()
}

func (g *runtimeGenerator) stop() bool {
	select {
	case <-g.ctx.Done():
		return false
	default:
		g.cancel()
		return true
	}
}

// Generator is as in python or ES6
// Generator functions allow you to declare a function that behaves like an iterator, i.e. it can be used in a for loop.
// see https://wiki.python.org/moin/Generators
type Generator struct {
	// Used by Next, to notify or deliver what is generated, as next in python or ES6
	C <-chan interface{}

	r runtimeGenerator
}

// Stop prevents the Generator from firing, with the channel drained.
// Stop ensures the channel is empty after a call to Stop.
// It returns true if the call stops the generator, false if the generator has already
// expired or been stopped.
// Stop does not close the channel, to prevent a read from the channel succeeding
// incorrectly.
//
// This cannot be done concurrent to other receives from the Generator's
// channel.
//
// For a generator created with GeneratorFunc(supplierC, f), if t.Stop returns false, then the generator
// has already expired and the function f has been started in its own goroutine;
// Stop does not wait for f to complete before returning.
// If the caller needs to know whether f is completed, it must coordinate
// with f explicitly.
func (g *Generator) Stop() bool {
	if g.r.f == nil || g.r.ctx == nil || g.r.cancel == nil {
		panic("generator: Stop called on uninitialized Generator")
	}
	return g.r.stop()
}

// Next behaves like an iterator, i.e. it can be used in a for loop.
// It's a grammar sugar for chan
func (g *Generator) Next() (msg interface{}, ok bool) {
	msg, ok = <-g.C
	return
}

func NewGenerator(supplierC <-chan interface{}) *Generator {
	c := make(chan interface{})

	ctx, cancel := context.WithCancel(context.Background())
	g := &Generator{
		C: c,
		r: runtimeGenerator{
			f:         sendChan,
			arg:       c,
			ctx:       ctx,
			cancel:    cancel,
			supplierC: supplierC,
			consumerC: c,
		},
	}
	g.r.start()
	return g
}

// GeneratorFunc waits for the supplierC to supply and then calls f
// in its own goroutine every time. It returns a Generator that can
// be used to cancel the call using its Stop method.
// Consume will be stopped when supplierC is closed.
func GeneratorFunc(supplierC <-chan interface{}, f func(msg interface{})) *Generator {
	ctx, cancel := context.WithCancel(context.Background())
	g := &Generator{
		r: runtimeGenerator{
			f:         goFunc,
			arg:       f,
			ctx:       ctx,
			cancel:    cancel,
			supplierC: supplierC,
		},
	}
	g.r.start()
	return g

}

func sendChan(c interface{}, msg interface{}) {
	// Non-blocking send of msg on c.
	// Used in NewGenerator, it cannot block anyway (buffer).
	// Used in NewGeneratorTicker, dropping sends on the floor is
	// the desired behavior when the reader gets behind,
	// because the sends are periodic.
	select {
	case c.(chan interface{}) <- msg:
	default:
	}
	return
}

func goFunc(arg interface{}, msg interface{}) {
	go arg.(func(interface{}))(msg)
}
