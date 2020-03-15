// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tcp

import (
	"io"
	"net"
	"sync"
	"time"
)

// connReader is the io.onMsgRead wrapper used by *conn. It combines a
// selectively-activated io.LimitedReader (to bound request header
// read sizes) with support for selectively keeping an io.onMsgRead.onMsgRead
// call blocked in a background goroutine to wait for activity and
// trigger a CloseNotifier channel.
type connReader struct {
	conn *conn

	mu      sync.Mutex // guards following
	hasByte bool
	byteBuf [1]byte
	cond    *sync.Cond
	inRead  bool
	aborted bool  // set true before conn.rwc deadline is set to past
	remain  int64 // bytes remaining
}

func (cr *connReader) lock() {
	cr.mu.Lock()
	if cr.cond == nil {
		cr.cond = sync.NewCond(&cr.mu)
	}
}

func (cr *connReader) unlock() { cr.mu.Unlock() }

func (cr *connReader) startBackgroundRead() {
	cr.lock()
	defer cr.unlock()
	if cr.inRead {
		panic("invalid concurrent Body.onMsgRead call")
	}
	if cr.hasByte {
		return
	}
	cr.inRead = true
	cr.conn.rwc.SetReadDeadline(time.Time{})
	go cr.backgroundRead()
}

// conn -> byteBuf -> p
func (cr *connReader) backgroundRead() {
	n, err := cr.conn.rwc.Read(cr.byteBuf[:])
	cr.lock()
	if n == 1 {
		cr.hasByte = true
		// We were past the end of the previous request's body already
		// (since we wouldn't be in a background read otherwise), so
		// this is a pipelined HTTP request. Prior to Go 1.11 we used to
		// send on the CloseNotify channel and cancel the context here,
		// but the behavior was documented as only "may", and we only
		// did that because that's how CloseNotify accidentally behaved
		// in very early Go releases prior to context support. Once we
		// added context support, people used a onMsgHandle's
		// Request.Context() and passed it along. Having that context
		// cancel on pipelined HTTP requests caused problems.
		// Fortunately, almost nothing uses HTTP/1.x pipelining.
		// Unfortunately, apt-get does, or sometimes does.
		// New Go 1.11 behavior: don't fire CloseNotify or cancel
		// contexts on pipelined requests. Shouldn't affect people, but
		// fixes cases like Issue 23921. This does mean that a client
		// closing their TCP connection after sending a pipelined
		// request won't cancel the context, but we'll catch that on any
		// write failure (in checkConnErrorWriter.Write).
		// If the server never writes, yes, there are still contrived
		// server & client behaviors where this fails to ever cancel the
		// context, but that's kinda why HTTP/1.x pipelining died
		// anyway.
	}
	if ne, ok := err.(net.Error); ok && cr.aborted && ne.Timeout() {
		// Ignore this error. It's the expected error from
		// another goroutine calling abortPendingRead.
	} else if err != nil {
		cr.handleReadError(err)
	}
	cr.aborted = false
	cr.inRead = false
	cr.unlock()
	cr.cond.Broadcast()
}

func (cr *connReader) abortPendingRead() {
	cr.lock()
	defer cr.unlock()
	if !cr.inRead {
		return
	}
	cr.aborted = true
	cr.conn.rwc.SetReadDeadline(aLongTimeAgo)
	for cr.inRead {
		cr.cond.Wait()
	}
	cr.conn.rwc.SetReadDeadline(time.Time{})
}

func (cr *connReader) setReadLimit(remain int64) { cr.remain = remain }
func (cr *connReader) setInfiniteReadLimit()     { cr.remain = maxInt64 }
func (cr *connReader) hitReadLimit() bool        { return cr.remain <= 0 }

// handleReadError is called whenever a onMsgRead from the client returns a
// non-nil error.
//
// The provided non-nil err is almost always io.EOF or a "use of
// closed network connection". In any case, the error is not
// particularly interesting, except perhaps for debugging during
// development. Any error means the connection is dead and we should
// down its context.
//
// It may be called from multiple goroutines.
func (cr *connReader) handleReadError(_ error) {
	cr.conn.cancelCtx()
}

func (cr *connReader) Read(p []byte) (n int, err error) {
	cr.lock()
	if cr.inRead {
		cr.unlock()
		panic("invalid concurrent Body.onMsgRead call")
	}
	if cr.hitReadLimit() {
		cr.unlock()
		return 0, io.EOF
	}
	if len(p) == 0 {
		cr.unlock()
		return 0, nil
	}
	if int64(len(p)) > cr.remain {
		p = p[:cr.remain]
	}
	if cr.hasByte {
		p[0] = cr.byteBuf[0]
		cr.hasByte = false
		cr.unlock()
		return 1, nil
	}
	cr.inRead = true
	cr.unlock()
	n, err = cr.conn.rwc.Read(p)

	cr.lock()
	cr.inRead = false
	if err != nil {
		cr.handleReadError(err)
	}
	cr.remain -= int64(n)
	cr.unlock()

	cr.cond.Broadcast()
	return n, err
}
