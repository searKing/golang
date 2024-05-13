// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package generator_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/searKing/golang/go/go/generator"
)

// Test the basic function calling behavior. Correct queueing
// behavior is tested elsewhere, since After and AfterFunc share
// the same code.
func TestGeneratorFuncWithSupplier(t *testing.T) {
	const unit = 25 * time.Millisecond
	var i, j int
	n := 10
	c := make(chan bool)
	supplierC := make(chan any)
	f := func(msg any) {
		j++
		if j == n {
			c <- true
		}
	}

	go func() {
		for {
			i++
			if i <= n {
				supplierC <- i
				if i == 0 {
					return
				}
				time.Sleep(unit)
				continue
			}
			return
		}
	}()

	generator.GeneratorFuncWithSupplier(supplierC, f)
	<-c
}

func TestGenerator_Stop(t *testing.T) {
	const unit = 25 * time.Millisecond
	var i, j int
	n := 10
	accept := 5
	c := make(chan bool)
	supplierC := make(chan any)
	var g *generator.Generator
	f := func(msg any) {
		j++
		if j == accept {
			g.Stop()
			c <- true
		}
	}

	go func() {
		for {
			i++
			if i <= n {
				supplierC <- i
				if i == 0 {
					return
				}
				time.Sleep(unit)
				continue
			}
			return
		}
	}()

	g = generator.GeneratorFuncWithSupplier(supplierC, f)
	<-c
}

func TestGenerator_Next(t *testing.T) {
	const unit = 25 * time.Millisecond
	var i, j int
	n := 10
	accept := 5
	c := make(chan bool)
	supplierC := make(chan any)
	go func() {
		for {
			i++
			if i <= n {
				supplierC <- i
				if i == 0 {
					return
				}
				time.Sleep(unit)
				continue
			}
			return
		}
	}()

	g := generator.GeneratorWithSupplier(supplierC)

	go func() {
		for {
			_, ok := g.Next()
			if !ok {
				break
			}
			j++
			if j == accept {
				g.Stop()
				break
			}
		}
		c <- true
	}()
	<-c
}

func TestGenerator_Yield(t *testing.T) {
	var g *generator.Generator

	supplierC := make(chan any, 100)
	supplierF := func(i int) {
		yield := g.Yield(supplierC)
		if !yield(i) {
			return
		}
		if !yield(i + 10) {
			return
		}
		close(supplierC)
	}
	g = generator.GeneratorWithSupplier(supplierC)
	supplierF(10)

	for msg := range g.C {
		fmt.Println(msg)
	}
}

func TestGeneratorVariadicFunc(t *testing.T) {
	g := generator.GeneratorVariadicFunc(func(yield generator.Yield, args ...any) {
		i := (args[0]).(int)
		if !yield(i) {
			return
		}
		if !yield(i + 10) {
			return
		}
	})

	gen := g(10)

	for msg := range gen.C {
		fmt.Println(msg)
	}
}

func TestGeneratorFuncClosure(t *testing.T) {
	g := func(i int) *generator.Generator {
		return generator.GeneratorFunc(func(yield generator.Yield) {
			if !yield(i) {
				return
			}
			if !yield(i + 10) {
				return
			}
		})
	}

	gen := g(10)

	for msg := range gen.C {
		fmt.Println(msg)
	}
}

func TestNewGeneratorFunc(t *testing.T) {
	i := 10
	gen := generator.GeneratorFunc(func(yield generator.Yield) {
		if !yield(i) {
			return
		}
		if !yield(i + 10) {
			return
		}
	})

	for msg := range gen.C {
		fmt.Println(msg)
	}
}
