package sync

import (
	"context"
	"fmt"
	"sync"

	"github.com/searKing/golang/go/errors"
	"github.com/searKing/golang/go/pragma"
	"github.com/searKing/golang/go/sync/atomic"
)

// Subject implements a condition variable like with channel, a rendezvous point
// for goroutines waiting for or announcing the occurrence
// of an event.
//
// The caller typically cannot assume that the condition is true when
// Subscribe chan returns. Instead, the caller should Wait in a loop:
//
//    time.After(timeout, c.PublishBroadcast()) // for timeout or periodic event
//    c.PublishBroadcast() // for async notify event directly
//    eventC, cancel := c.Subscribe()
//    for !condition() {
//        select{
//        case event, closed := <- eventC:
//            ... make use of event ...
//        }
//    }
//    ... make use of condition ...
//
type Subject struct {
	noCopy pragma.DoNotCopy

	mu          sync.Mutex
	subscribers map[*subscriber]struct{}

	inShutdown atomic.Bool // true when when server is in shutdown
}

type subscriber struct {
	mu   sync.Mutex // guard close of channel msgC
	msgC chan interface{}

	once  sync.Once
	doneC chan struct{} // closed when when subscriber is in shutdown, like removed.
}

func (s *subscriber) Shutdown() {
	if s == nil {
		return
	}
	s.once.Do(func() {
		close(s.doneC)
		s.mu.Lock()
		defer s.mu.Unlock()
		close(s.msgC)
	})
}

// publish wakes a listener waiting on c to consume the event.
// event will be dropped if ctx is Done before event is received.
func (s *subscriber) publish(ctx context.Context, event interface{}) error {
	// guard of msgC's close
	s.mu.Lock()
	defer s.mu.Unlock()
	select {
	case <-ctx.Done():
		// event dropped because of publisher
		return ctx.Err()
	case <-s.doneC:
		// event dropped because of subscriber
		return fmt.Errorf("event dropped because of subscriber unsubscribed")
	default:
	}

	select {
	case <-ctx.Done():
		// event dropped because of publisher
		return ctx.Err()
	case <-s.doneC:
		// event dropped because of subscriber
		return fmt.Errorf("event dropped because of subscriber unsubscribed")
	case s.msgC <- event:
		// event consumed
		return nil
	}
}

// Subscribe returns a channel that's closed when awoken by PublishSignal or PublishBroadcast.
// never be canceled. Successive calls to Subscribe return different values.
// The close of the Subscribe channel may happen asynchronously,
// after the cancel function returns.
func (s *Subject) Subscribe() (<-chan interface{}, context.CancelFunc) {
	listener := &subscriber{
		msgC:  make(chan interface{}),
		doneC: make(chan struct{}),
	}
	s.trackChannel(listener, true)
	return listener.msgC, func() {
		s.trackChannel(listener, false)
	}
}

// PublishSignal wakes one listener waiting on c, if there is any.
// PublishSignal blocks until event is received or dropped.
func (s *Subject) PublishSignal(ctx context.Context, event interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	var wg sync.WaitGroup
	var errs []error
	for listener := range s.subscribers {
		wg.Add(1)
		go func(listener *subscriber) {
			defer wg.Done()
			err := listener.publish(ctx, event)
			if err != nil {
				errs = append(errs, err)
			}
		}(listener)
		break
	}
	wg.Wait()
	return errors.Multi(errs...)
}

// PublishBroadcast wakes all listeners waiting on c.
// PublishBroadcast blocks until event is received or dropped.
// event will be dropped if ctx is Done before event is received.
func (s *Subject) PublishBroadcast(ctx context.Context, event interface{}) error {
	var wg sync.WaitGroup
	var errs []error
	func() {
		s.mu.Lock()
		defer s.mu.Unlock()
		for listener := range s.subscribers {
			wg.Add(1)
			go func(listener *subscriber) {
				defer wg.Done()
				err := listener.publish(ctx, event)
				if err != nil {
					errs = append(errs, err)
				}
			}(listener)
		}
	}()
	wg.Wait()
	return errors.Multi(errs...)
}

func (s *Subject) trackChannel(c *subscriber, add bool) {
	func() {
		s.mu.Lock()
		defer s.mu.Unlock()
		if s.subscribers == nil {
			s.subscribers = make(map[*subscriber]struct{})
		}
		_, has := s.subscribers[c]
		if has {
			if add {
				return
			}
			delete(s.subscribers, c)
			return
		}
		if add {
			s.subscribers[c] = struct{}{}
		}
	}()
	if !add {
		c.Shutdown()
	}
}
