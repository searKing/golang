package sync

import (
	"context"
	"github.com/searKing/golang/go/pragma"
	"sync"
)

// The caller typically cannot assume that the condition is true when
// Subscribe chan returns. Instead, the caller should Wait in a loop:
//
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

	mu       sync.Mutex
	channels map[chan interface{}]struct{}

	doneC chan struct{}

	checker pragma.CopyChecker
}

// Subscribe returns a channel that's closed when awoken by PublishSignal or PublishBroadcast.
// never be canceled. Successive calls to Subscribe return different values.
// The close of the Subscribe channel may happen asynchronously,
// after the cancel function returns.
func (s *Subject) Subscribe() (<-chan interface{}, context.CancelFunc) {
	eventC := make(chan interface{})
	s.trackChannel(eventC, true)
	return eventC, func() {
		s.trackChannel(eventC, false)
	}
}

// PublishSignal wakes one listener waiting on c, if there is any.
func (s *Subject) PublishSignal(ctx context.Context, event interface{}) {
	s.checker.Check()
	s.mu.Lock()
	defer s.mu.Unlock()
	for eventC := range s.channels {
		go s.publish(ctx, event, eventC)
		break
	}
}

// PublishBroadcast wakes all listeners waiting on c.
func (s *Subject) PublishBroadcast(ctx context.Context, event interface{}) {
	s.checker.Check()
	s.mu.Lock()
	defer s.mu.Unlock()
	for eventC := range s.channels {
		go s.publish(ctx, event, eventC)
	}
}

func (s *Subject) publish(ctx context.Context, event interface{}, ch chan interface{}) {
	select {
	case <-ctx.Done():
		s.trackChannel(ch, false)
	case ch <- event:
	}
}

func (s *Subject) trackChannel(c chan interface{}, add bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.channels == nil {
		s.channels = make(map[chan interface{}]struct{})
	}
	if add {
		s.channels[c] = struct{}{}
	} else {
		delete(s.channels, c)
		close(c)
	}
}
