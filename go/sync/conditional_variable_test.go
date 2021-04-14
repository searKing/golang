package sync_test

import (
	"context"
	"sync"
	"testing"

	sync_ "github.com/searKing/golang/go/sync"
)

func TestConditionVariable_Wait(t *testing.T) {
	var m sync.Mutex
	var c sync_.ConditionVariable
	n := 2
	running := make(chan bool, n)
	awake := make(chan bool, n)
	for i := 0; i < n; i++ {
		go func() {
			m.Lock()
			running <- true
			c.Wait(&m)
			awake <- true
			m.Unlock()
		}()
	}
	for i := 0; i < n; i++ {
		<-running // Wait for everyone to run.
	}
	for n > 0 {
		select {
		case <-awake:
			t.Fatal("goroutine not asleep")
		default:
		}
		m.Lock()
		c.Signal()
		m.Unlock()
		<-awake // Will deadlock if no goroutine wakes up
		select {
		case <-awake:
			t.Fatal("too many goroutines awake")
		default:
		}
		n--
	}
	c.Signal()
}

func TestConditionVariable_WaitPred(t *testing.T) {
	var m sync.Mutex
	var c sync_.ConditionVariable
	n := 2
	running := make(chan bool, n)
	awake := make(chan bool, n)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		m.Lock()
		running <- true
		c.WaitPred(&m, func() bool {
			return ctx.Err() != nil
		})
		awake <- true
		m.Unlock()
	}()
	<-running // Wait for everyone to run.

	select {
	case <-awake:
		t.Fatal("goroutine not asleep")
	default:
	}
	m.Lock()
	c.Signal()
	m.Unlock()
	cancel()
	<-awake // Will deadlock if no goroutine wakes up
	select {
	case <-awake:
		t.Fatal("too many goroutines awake")
	default:
	}
	n--
}
