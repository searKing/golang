package sync_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/pkg/errors"
	sync_ "github.com/searKing/golang/go/sync"
)

func TestSubject_PublishSignal(t *testing.T) {
	var s sync_.Subject
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
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			err := s.PublishSignal(ctx, nil)
			if err != nil {
				t.Fatalf("PublishSignal: %s", err)
				return
			}
		}()
		<-awake // Will deadlock if no goroutine wakes up
		select {
		case <-awake:
			t.Fatal("too many goroutines awake")
		default:
		}
		n--
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := s.PublishSignal(ctx, nil)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("PublishSignal: %s", err)
		return
	}
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
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			err := s.PublishBroadcast(ctx, nil)
			if err != nil {
				t.Fatalf("PublishBroadcast: %s", err)
				return
			}
		}()
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := s.PublishBroadcast(ctx, nil)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("PublishBroadcast: %s", err)
		return
	}
}
