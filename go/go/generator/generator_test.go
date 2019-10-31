package generator_test

import (
	"testing"
	"time"

	"github.com/searKing/golang/go/go/generator"
)

// Test the basic function calling behavior. Correct queueing
// behavior is tested elsewhere, since After and AfterFunc share
// the same code.
func TestGeneratorFunc(t *testing.T) {
	const unit = 25 * time.Millisecond
	var i, j int
	n := 10
	c := make(chan bool)
	supplierC := make(chan interface{})
	f := func(msg interface{}) {
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

	generator.GeneratorFunc(supplierC, f)
	<-c
}

func TestGenerator_Stop(t *testing.T) {
	const unit = 25 * time.Millisecond
	var i, j int
	n := 10
	accept := 5
	c := make(chan bool)
	supplierC := make(chan interface{})
	var g *generator.Generator
	f := func(msg interface{}) {
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

	g = generator.GeneratorFunc(supplierC, f)
	<-c
}

func TestGenerator_Next(t *testing.T) {
	const unit = 25 * time.Millisecond
	var i, j int
	n := 10
	accept := 5
	c := make(chan bool)
	supplierC := make(chan interface{})
	var g *generator.Generator
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

	g = generator.NewGenerator(supplierC)

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
