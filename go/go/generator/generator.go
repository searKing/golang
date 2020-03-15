// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package generator

import (
	"context"
)

type Yield func(msg interface{}) (ok bool)

// Generator behaves like Generator in python or ES6
// Generator function contains one or more yield statement.
// Generator functions allow you to declare a function that behaves like an iterator, i.e. it can be used in a for loop.
// see https://wiki.python.org/moin/Generators
// see https://www.programiz.com/python-programming/generator
// see https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Statements/function*
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
// For a generator created with GeneratorFuncWithSupplier(supplierC, f), if t.Stop returns false, then the generator
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

func (g *Generator) StoppedC() context.Context {
	if g.r.f == nil || g.r.ctx == nil || g.r.cancel == nil {
		panic("generator: StoppedC called on uninitialized Generator")
	}
	return g.r.ctx
}

func (g *Generator) Stopped() bool {
	if g.r.f == nil || g.r.ctx == nil || g.r.cancel == nil {
		panic("generator: Stopped called on uninitialized Generator")
	}
	return g.r.stopped()
}

// Next behaves like an iterator, i.e. it can be used in a for loop.
// It's a grammar sugar for chan
func (g *Generator) Next() (msg interface{}, ok bool) {
	msg, ok = <-g.C
	return
}

// Yield is a grammar sugar for data src of generator
// ok returns true if msg sent; false if consume canceled
// If a function contains at least one yield statement (it may contain other yield or return statements),
// it becomes a generator function. Both yield and return will return some value from a function.
// The difference is that, while a return statement terminates a function entirely,
// yield statement pauses the function saving all its states and later continues from there on successive calls.
func (g *Generator) Yield(supplierC chan<- interface{}) Yield {
	return func(msg interface{}) (ok bool) {
		select {
		case <-g.StoppedC().Done():
			return false
		case supplierC <- msg:
			return true
		}
	}
}

// Simply speaking, a generator is a function that returns an object (iterator) which we can iterate over (one value at a time).

// GeneratorFunc returns an object (iterator) which we can iterate over (one value at a time).
// It returns a Generator that can be used to cancel the call using its Stop method.
// Iterate will be stopped when f is return or Stop is called.
func GeneratorFunc(f func(yield Yield)) *Generator {
	supplierC := make(chan interface{})
	g := GeneratorWithSupplier(supplierC)
	supplierF := func(args ...interface{}) {
		yield := g.Yield(supplierC)
		f(yield)
		close(supplierC)
	}
	go supplierF()
	return g
}

// GeneratorWithSupplier is like GeneratorFunc.
// But it's data src is from supplierC.
// Iterate will be stopped when supplierC is closed or Stop is called.
func GeneratorWithSupplier(supplierC <-chan interface{}) *Generator {
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

// Deprecated: Use GeneratorFunc wrapped by closure instead.
func GeneratorVariadicFunc(f func(yield Yield, args ...interface{})) func(args ...interface{}) *Generator {
	supplierC := make(chan interface{})
	g := GeneratorWithSupplier(supplierC)
	supplierF := func(args ...interface{}) {
		yield := g.Yield(supplierC)
		f(yield, args...)
		close(supplierC)
	}
	return func(args ...interface{}) *Generator {
		go supplierF(args...)
		return g
	}
}

// GeneratorFuncWithSupplier waits for the supplierC to supply and then calls f
// in its own goroutine every time. It returns a Generator that can
// be used to cancel the call using its Stop method.
// Consume will be stopped when supplierC is closed.
func GeneratorFuncWithSupplier(supplierC <-chan interface{}, f func(msg interface{})) *Generator {
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

func sendChan(ctx context.Context, c interface{}, msg interface{}) {
	// Blocking send of msg on c.
	select {
	case <-ctx.Done():
		return
	case c.(chan interface{}) <- msg:
	}
	return
}

// arg -> func
func goFunc(ctx context.Context, f interface{}, msg interface{}) {
	go f.(func(interface{}))(msg)
}
