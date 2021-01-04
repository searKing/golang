package sync_test

import (
	"context"
	sync_ "github.com/searKing/golang/go/sync"
	"sync"
	"testing"
)

func TestSubject_PublishSignal(t *testing.T) {
	s := sync_.Subject{}
	n := 2
	awake := make(chan bool, n)
	var wg sync.WaitGroup
	eventC, _ := s.Subscribe()
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			wg.Done()
			<-eventC
			awake <- true
		}()
	}
	// Wait for everyone to run.
	wg.Wait()
	for n > 0 {
		select {
		case <-awake:
			t.Fatal("goroutine not asleep")
		default:
		}
		s.PublishSignal(context.Background(), nil)
		<-awake // Will deadlock if no goroutine wakes up
		select {
		case <-awake:
			t.Fatal("too many goroutines awake")
		default:
		}
		n--
	}
	s.PublishSignal(context.Background(), nil)
}

func TestSubject_PublishBroadcast(t *testing.T) {
	s := sync_.Subject{}
	n := 2
	awake := make(chan bool, n)
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		eventC, _ := s.Subscribe()
		wg.Add(1)
		go func() {
			wg.Done()
			<-eventC
			awake <- true
		}()
	}
	// Wait for everyone to run.
	wg.Wait()
	{
		select {
		case <-awake:
			t.Fatal("goroutine not asleep")
		default:
		}
		s.PublishBroadcast(context.Background(), nil)
		for n > 0 {
			<-awake // Will deadlock if no goroutine wakes up
			n--
		}
		select {
		case <-awake:
			t.Fatal("too many goroutines awake")
		default:
		}
	}
}
