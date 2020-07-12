// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package net

import "net"

// type check
var _ net.Error = ErrListenerClosed

type errListenerClosed string

func (e errListenerClosed) Error() string   { return string(e) }
func (e errListenerClosed) Temporary() bool { return false }
func (e errListenerClosed) Timeout() bool   { return false }

// ErrListenerClosed is returned from Listener.Accept when the underlying
// listener is closed.
var ErrListenerClosed = errListenerClosed("net: listener closed")
