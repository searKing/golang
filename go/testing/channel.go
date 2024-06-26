// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testing

import (
	"context"
)

// DefaultChanBufferSize is the default buffer size of the underlying channel.
const DefaultChanBufferSize = 1

// Channel wraps a generic channel and provides a timed receive operation.
type Channel struct {
	ch chan any
}

// Send sends value on the underlying channel.
func (c *Channel) Send(value any) {
	c.ch <- value
}

// SendContext sends value on the underlying channel, or returns an error if
// the context expires.
func (c *Channel) SendContext(ctx context.Context, value any) error {
	select {
	case c.ch <- value:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// SendOrFail attempts to send value on the underlying channel.  Returns true
// if successful or false if the channel was full.
func (c *Channel) SendOrFail(value any) bool {
	select {
	case c.ch <- value:
		return true
	default:
		return false
	}
}

// ReceiveOrFail returns the value on the underlying channel and true, or nil
// and false if the channel was empty.
func (c *Channel) ReceiveOrFail() (any, bool) {
	select {
	case got := <-c.ch:
		return got, true
	default:
		return nil, false
	}
}

// Receive returns the value received on the underlying channel, or the error
// returned by ctx if it is closed or cancelled.
func (c *Channel) Receive(ctx context.Context) (any, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case got := <-c.ch:
		return got, nil
	}
}

// Replace clears the value on the underlying channel, and sends the new value.
//
// It's expected to be used with a size-1 channel, to only keep the most
// up-to-date item. This method is inherently racy when invoked concurrently
// from multiple goroutines.
func (c *Channel) Replace(value any) {
	for {
		select {
		case c.ch <- value:
			return
		case <-c.ch:
		}
	}
}

// Clear drains all values on the underlying channel, or the error
// returned by ctx if it is closed or cancelled.
func (c *Channel) Clear(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-c.ch:
			continue
		default:
			return nil
		}
	}
}

// NewChannel returns a new Channel.
func NewChannel() *Channel {
	return NewChannelWithSize(DefaultChanBufferSize)
}

// NewChannelWithSize returns a new Channel with a buffer of bufSize.
func NewChannelWithSize(bufSize int) *Channel {
	return &Channel{ch: make(chan any, bufSize)}
}
