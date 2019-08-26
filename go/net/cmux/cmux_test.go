// Copyright 2019 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmux_test

import (
	"context"
	"github.com/searKing/golang/go/net/cmux"
	"github.com/searKing/golang/go/testing/leakcheck"
	"io"
	"log"
	"net"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"golang.org/x/net/http2"
)

const (
	handleHTTP1Close   = 1
	handleHTTP1Request = 2
	handleAnyClose     = 3
	handleAnyRequest   = 4
)

func TestTimeout(t *testing.T) {
	defer leakcheck.Check(t)
	lis := testListener(t)
	defer lis.Close()
	result := make(chan int, 5)
	testDuration := time.Millisecond * 500
	m := cmux.New(context.Background())
	m.SetReadTimeout(testDuration)
	http1Listener := m.Match(cmux.HTTP1Fast())
	defer http1Listener.Close()
	anyListener := m.Match(cmux.Any())
	defer anyListener.Close()

	ctx, cancelFn := context.WithCancel(context.TODO())
	defer cancelFn()
	go func() {
		_ = m.Serve(lis)
	}()
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			con, err := http1Listener.Accept()
			if err != nil {
				result <- handleHTTP1Close
			} else {
				_, _ = con.Write([]byte("http1Listener"))
				result <- handleHTTP1Request
				select {
				case <-ctx.Done():
					break
				}
				_ = con.Close()
			}
			select {
			case <-ctx.Done():
				return
			}
		}
	}()
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			con, err := anyListener.Accept()
			if err != nil {
				result <- handleAnyClose
			} else {
				_, err = con.Write([]byte("any"))
				result <- handleAnyRequest
				select {
				case <-ctx.Done():
					break
				}
				_ = con.Close()
			}
		}
	}()
	time.Sleep(testDuration) // wait to prevent timeouts on slow test-runners
	client, err := net.Dial("tcp", lis.Addr().String())
	if err != nil {
		log.Fatal("testTimeout client failed: ", err)
	}
	defer func() {
		_ = client.Close()
	}()
	time.Sleep(testDuration / 2)
	if len(result) != 0 {
		log.Print("tcp ")
		t.Fatal("testTimeout failed: accepted to fast: ", len(result))
	}
	//_ = client.SetReadDeadline(time.Now().Add(testDuration * 3))
	buffer := make([]byte, 10)
	rl, err := client.Read(buffer)
	if err != nil {
		t.Fatal("testTimeout failed: client error: ", err, rl)
	}
	_ = m.Close()
	if rl != 3 {
		log.Print("testTimeout failed: response from wrong sevice ", rl)
	}
	if string(buffer[0:3]) != "any" {
		log.Print("testTimeout failed: response from wrong sevice ")
	}
	time.Sleep(testDuration * 2)
	if len(result) != 2 {
		t.Fatal("testTimeout failed: accepted to less: ", len(result))
	}
	if a := <-result; a != handleAnyRequest {
		t.Fatal("testTimeout failed: anyListener rule did not match")
	}
	if a := <-result; a != handleHTTP1Close {
		t.Fatal("testTimeout failed: no close an http rule")
	}
}

func TestRead(t *testing.T) {
	defer leakcheck.Check(t)
	errCh := make(chan error)
	defer func() {
		select {
		case err := <-errCh:
			t.Fatal(err)
		default:
		}
	}()
	const payload = "hello world\r\n"
	const mult = 2

	writer, reader := net.Pipe()
	go func() {
		if _, err := io.WriteString(writer, strings.Repeat(payload, mult)); err != nil {
			t.Fatal(err)
		}
		if err := writer.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	l := newChanListener()
	l.Notify(reader)
	defer l.Close()
	muxl := cmux.New(context.Background())
	defer muxl.Close()

	// Register a bogus matcher to force buffering exactly the right amount.
	// Before this fix, this would trigger a bug where `Read` would incorrectly
	// report `io.EOF` when only the buffer had been consumed.
	muxl.Match(cmux.MatcherFunc(func(w io.Writer, r io.Reader) bool {
		var b [len(payload)]byte
		_, _ = r.Read(b[:])
		return false
	}))
	anyl := muxl.Match(cmux.Any())
	go safeServe(errCh, muxl, l)
	muxedConn, err := anyl.Accept()
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < mult; i++ {
		var b [len(payload)]byte
		n, err := muxedConn.Read(b[:])
		if err != nil {
			t.Error(err)
			continue
		}
		if e := len(b); n != e {
			t.Errorf("expected to read %d bytes, but read %d bytes", e, n)
		}
	}
	var b [1]byte
	if _, err := muxedConn.Read(b[:]); err != io.EOF {
		t.Errorf("unexpected error %v, expected %v", err, io.EOF)
	}

}

func TestAny(t *testing.T) {
	defer leakcheck.Check(t)
	errCh := make(chan error)
	defer func() {
		for {
			select {
			case err, ok := <-errCh:
				if !ok {
					return
				}
				t.Fatal(err)
			default:
				close(errCh)
				return
			}
		}
	}()
	l := testListener(t)
	defer l.Close()

	muxl := cmux.New(context.Background())
	defer muxl.Close()

	httpl := muxl.Match(cmux.Any())

	go runTestHTTPServer(errCh, httpl)
	go safeServe(errCh, muxl, l)

	runTestHTTP1Client(t, l.Addr())
}

func TestTLS(t *testing.T) {
	generateTLSCert(t)
	defer cleanupTLSCert(t)
	defer leakcheck.Check(t)
	errCh := make(chan error)
	defer func() {
		for {
			select {
			case err, ok := <-errCh:
				if !ok {
					return
				}
				t.Fatal(err)
			default:
				close(errCh)
				return
			}
		}
	}()
	l := testListener(t)
	defer l.Close()

	muxl := cmux.New(context.Background())
	defer muxl.Close()
	tlsl := muxl.Match(cmux.TLS())
	httpl := muxl.Match(cmux.Any())

	go runTestTLSServer(errCh, tlsl)
	go runTestHTTPServer(errCh, httpl)
	go safeServe(errCh, muxl, l)

	runTestHTTP1Client(t, l.Addr())
	runTestTLSClient(t, l.Addr())
}

func TestHTTP2(t *testing.T) {
	defer leakcheck.Check(t)
	errCh := make(chan error)
	defer func() {
		for {
			select {
			case err, ok := <-errCh:
				if !ok {
					return
				}
				t.Fatal(err)
			default:
				close(errCh)
				return
			}
		}
	}()
	writer, reader := net.Pipe()
	go func() {
		if _, err := io.WriteString(writer, http2.ClientPreface); err != nil {
			t.Fatal(err)
		}
		if err := writer.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	l := newChanListener()
	l.Notify(reader)
	muxl := cmux.New(context.Background())
	defer muxl.Close()
	// Register a bogus matcher that only reads one byte.
	muxl.Match(cmux.MatcherFunc(func(w io.Writer, r io.Reader) bool {
		var b [1]byte
		_, _ = r.Read(b[:])
		return false
	}))
	h2l := muxl.Match(cmux.HTTP2())
	go safeServe(errCh, muxl, l)
	muxedConn, err := h2l.Accept()
	_ = l.Close()
	if err != nil {
		t.Fatal(err)
	}
	var b [len(http2.ClientPreface)]byte
	var n int
	// We have the sniffed buffer first...
	if n, err = muxedConn.Read(b[:]); err == io.EOF {
		t.Fatal(err)
	}
	// and then we read from the source.
	if _, err = muxedConn.Read(b[n:]); err != nil && err != io.EOF {
		t.Fatal(err)
	}
	if string(b[:]) != http2.ClientPreface {
		t.Errorf("got unexpected read %s, expected %s", b, http2.ClientPreface)
	}
}

func TestHTTP2MatchHeaderField(t *testing.T) {
	testHTTP2HeaderField(t, cmux.HTTP2HeaderFieldEqual, "value", "value", "anothervalue")
}

func TestHTTP2MatchHeaderFieldPrefix(t *testing.T) {
	testHTTP2HeaderField(t, cmux.HTTP2HeaderFieldPrefix, "application/grpc+proto", "application/grpc", "application/json")
}

func TestHTTPGoRPC(t *testing.T) {
	defer leakcheck.Check(t)
	errCh := make(chan error)
	defer func() {
		for {
			select {
			case err, ok := <-errCh:
				if !ok {
					return
				}
				t.Fatal(err)
			default:
				close(errCh)
				return
			}
		}
	}()
	l := testListener(t)
	defer l.Close()

	muxl := cmux.New(context.Background())
	defer muxl.Close()
	httpl := muxl.Match(cmux.HTTP2(), cmux.HTTP1Fast())
	rpcl := muxl.Match(cmux.Any())

	go runTestHTTPServer(errCh, httpl)
	go runTestRPCServer(errCh, rpcl)
	go safeServe(errCh, muxl, l)

	runTestHTTP1Client(t, l.Addr())
	runTestRPCClient(t, l.Addr())
}

func TestErrorHandler(t *testing.T) {
	defer leakcheck.Check(t)
	errCh := make(chan error)
	defer func() {
		for {
			select {
			case err, ok := <-errCh:
				if !ok {
					return
				}
				t.Fatal(err)
			default:
				close(errCh)
				return
			}
		}
	}()
	l := testListener(t)
	defer l.Close()

	muxl := cmux.New(context.Background())
	defer muxl.Close()
	httpl := muxl.Match(cmux.HTTP2(), cmux.HTTP1Fast())

	go runTestHTTPServer(errCh, httpl)
	go safeServe(errCh, muxl, l)

	var errCount uint32
	muxl.HandleError(cmux.ErrorHandlerFunc(func(err error) bool {
		if atomic.AddUint32(&errCount, 1) == 1 {
			if _, ok := err.(cmux.ErrNotMatched); !ok {
				t.Errorf("unexpected error: %v", err)
			}
		}
		return true
	}))

	c, clean := safeDial(t, l.Addr())
	defer clean()

	var num int
	for atomic.LoadUint32(&errCount) == 0 {
		if err := c.Call("TestRPCRcvr.Test", rpcVal, &num); err == nil {
			// The connection is simply closed.
			t.Errorf("unexpected rpc success after %d errors", atomic.LoadUint32(&errCount))
		}
	}
}

func TestMultipleMatchers(t *testing.T) {
	defer leakcheck.Check(t)
	errCh := make(chan error)
	defer func() {
		for {
			select {
			case err, ok := <-errCh:
				if !ok {
					return
				}
				t.Fatal(err)
			default:
				close(errCh)
				return
			}
		}
	}()
	l := testListener(t)
	defer l.Close()

	matcher := func(w io.Writer, r io.Reader) bool {
		return true
	}
	unmatcher := func(w io.Writer, r io.Reader) bool {
		return false
	}

	muxl := cmux.New(context.Background())
	defer muxl.Close()
	lis := muxl.Match(cmux.MatcherFunc(unmatcher), cmux.MatcherFunc(matcher), cmux.MatcherFunc(unmatcher))

	go runTestHTTPServer(errCh, lis)
	go safeServe(errCh, muxl, l)

	runTestHTTP1Client(t, l.Addr())
}

func TestClose(t *testing.T) {
	defer leakcheck.Check(t)
	errCh := make(chan error)
	defer func() {
		for {
			select {
			case err, ok := <-errCh:
				if !ok {
					return
				}
				t.Fatal(err)
			default:
				close(errCh)
				return
			}
		}
	}()
	l := newChanListener()

	c1, c2 := net.Pipe()

	muxl := cmux.New(context.Background())
	defer muxl.Close()
	anyl := muxl.Match(cmux.Any())

	go safeServe(errCh, muxl, l)

	l.Notify(c1)

	// First connection goes through.
	if _, err := anyl.Accept(); err != nil {
		t.Fatal(err)
	}

	// Second connection is sent
	l.Notify(c2)

	// Listener is closed.
	l.Close()

	// Second connection either goes through or it is closed.
	if _, err := anyl.Accept(); err != nil {
		if err != cmux.ErrListenerClosed {
			t.Fatal(err)
		}
		// The error is either io.ErrClosedPipe or net.OpError wrapping
		// a net.pipeError depending on the go version.
		if _, err := c2.Read([]byte{}); !strings.Contains(err.Error(), "closed") {
			t.Fatalf("connection is not closed and is leaked: %v", err)
		}
	}
}
