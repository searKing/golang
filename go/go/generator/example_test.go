package generator_test

import (
	"fmt"

	"github.com/searKing/golang/go/go/generator"
)

// see https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Statements/function*
// function* generator(i) {
//  yield i;
//  yield i + 10;
// }
//
// var gen = generator(10);
//
// console.log(gen.next().value);
// // expected output: 10
//
// console.log(gen.next().value);
// // expected output: 20
func ExampleGeneratorFunc() {
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

	// Output:
	// 10
	// 20
}
func ExampleGenerator_Next() {
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

	for {
		msg, ok := gen.Next()
		if !ok {
			return
		}
		fmt.Println(msg)
	}

	// Output:
	// 10
	// 20
}

func ExampleGeneratorWithSupplier() {
	var g *generator.Generator

	supplierC := make(chan interface{})
	supplierF := func(i int) {
		consumer := g.Yield(supplierC)
		if !consumer(i) {
			return
		}
		if !consumer(i + 10) {
			return
		}
		close(supplierC)
	}
	g = generator.GeneratorWithSupplier(supplierC)
	go supplierF(10)

	for msg := range g.C {
		fmt.Println(msg)
	}

	// Output:
	// 10
	// 20
}

func ExampleGeneratorAdvanceFunc() {
	g := func(i int) *generator.Generator {
		g := generator.GeneratorAdvanceFunc(func(yield generator.Yield, args ...interface{}) {
			i := (args[0]).(int)
			if !yield(i) {
				return
			}
			if !yield(i + 10) {
				return
			}
		})

		gen := g(i)
		return gen
	}

	gen := g(10)

	for msg := range gen.C {
		fmt.Println(msg)
	}

	// Output:
	// 10
	// 20
}
