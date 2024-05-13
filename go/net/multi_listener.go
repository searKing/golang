// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package net

import (
	"log"
	"net"
	"sync"
	"time"

	"github.com/searKing/golang/go/errors"
)

// MultiListener is a net.Listener that accepts all the connections from all listeners.
type MultiListener struct {
	*NotifyListener
	listeners multiAddrs
	// ErrorLog specifies an optional logger for errors accepting
	// connections, unexpected behavior from handlers, and
	// underlying FileSystem errors.
	// If nil, logging is done via the log package's standard logger.
	ErrorLog *log.Logger

	once sync.Once
}

func NewMultiListener(listeners ...net.Listener) *MultiListener {
	return &MultiListener{
		NotifyListener: NewNotifyListener(),
		listeners:      listeners,
	}
}

// Addr returns the listener's network address.
func (l *MultiListener) Addr() net.Addr {
	return l.listeners
}
func (l *MultiListener) Accept() (net.Conn, error) {
	if len(l.listeners) == 0 {
		return nil, ErrListenerClosed
	}
	l.once.Do(func() {
		for _, listener := range l.listeners {
			go l.serve(listener)
		}
	})
	return l.NotifyListener.Accept()
}

func (l *MultiListener) Close() error {
	var errs []error
	errs = append(errs, l.NotifyListener.Close())
	for _, listener := range l.listeners {
		errs = append(errs, listener.Close())
	}
	return errors.Multi(errs...)
}

func (l *MultiListener) logf(format string, args ...any) {
	if l.ErrorLog != nil {
		l.ErrorLog.Printf(format, args...)
	} else {
		log.Printf(format, args...)
	}
}

// serve accepts incoming connections on the Listener lis, send the
// conn accepted to the NotifyListener for each.
//
// Serve always returns a non-nil error and closes l.
// After Close, the returned error is ErrListenerClosed.
func (l *MultiListener) serve(lis net.Listener) error {
	defer l.Close()
	var tempDelay time.Duration // how long to sleep on accept failure
	for {
		// Accept waits for and returns the next connection to the listener.
		conn, err := lis.Accept()
		if err != nil {
			select {
			case <-l.DoneC():
				return ErrListenerClosed
			default:
			}

			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				l.logf("multi listener: Accept error: %v; retrying in %v", err, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			return err
		}
		tempDelay = 0

		select {
		case <-l.DoneC():
			return ErrListenerClosed
		case l.C <- conn:
		}
	}
}
