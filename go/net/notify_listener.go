// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package net

import (
	"net"
	"sync"
)

var _ net.Listener = (*NotifyListener)(nil)

type notifyAddr string

func (notifyAddr) Network() string  { return "notify+net" }
func (f notifyAddr) String() string { return string(f) }

type NotifyListener struct {
	C chan net.Conn

	mu       sync.Mutex
	doneChan chan struct{}
}

// NewNotifyListener creates a new Listener that will recv
// conns on its channel.
func NewNotifyListener() *NotifyListener {
	c := make(chan net.Conn)

	l := &NotifyListener{
		C: c,
	}
	return l
}

// Addr returns the listener's network address.
func (l *NotifyListener) Addr() net.Addr {
	return notifyAddr("notify_listener")
}

func (l *NotifyListener) Accept() (net.Conn, error) {
	select {
	case <-l.getDoneChan():
		return nil, ErrListenerClosed
	case c, ok := <-l.C:
		if !ok {
			// Already closed. Don't Accept again.
			return nil, ErrListenerClosed
		}
		return c, nil
	}
}

// DoneC returns whether this listener has been closed
// for multi producers of C
func (l *NotifyListener) DoneC() <-chan struct{} {
	return l.getDoneChan()
}

// Stop prevents the NotifyListener from firing.
// To ensure the channel is empty after a call to Close, check the
// return value and drain the channel.
func (l *NotifyListener) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.closeDoneChanLocked()

	return l.drainChanLocked()
}

func (l *NotifyListener) drainChanLocked() error {
	var err error
	// Drain the connections enqueued for the listener.
L:
	for {
		select {
		case c, ok := <-l.C:
			if !ok {
				// Already closed. Don't close again.
				return nil
			}
			if c == nil {
				continue
			}
			if cerr := c.Close(); cerr != nil && err == nil {
				err = cerr
			}
		default:
			break L
		}
	}
	return err
}

func (l *NotifyListener) getDoneChan() <-chan struct{} {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.getDoneChanLocked()
}

func (l *NotifyListener) getDoneChanLocked() chan struct{} {
	if l.doneChan == nil {
		l.doneChan = make(chan struct{})
	}
	return l.doneChan
}

func (l *NotifyListener) closeDoneChanLocked() {
	ch := l.getDoneChanLocked()
	select {
	case <-ch:
		// Already closed. Don't close again.
	default:
		// Safe to close here. We're the only closer, guarded
		// by s.mu.
		close(ch)
	}
}
