// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mux_test

import (
	"context"
	"io"
	"log"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/searKing/golang/go/net/mux"
	"github.com/searKing/golang/go/testing/leakcheck"

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
	loopbackLis := testListener(t)
	defer loopbackLis.Close()
	http1Result := make(chan int, 5)
	anyResult := make(chan int, 5)
	testDuration := time.Millisecond * 500
	muxer := mux.NewServeMux()
	muxer.SetReadTimeout(testDuration)
	http1Listener := muxer.HandleListener(mux.HTTP1Fast())
	defer http1Listener.Close()
	anyListener := muxer.HandleListener(mux.Any())
	defer anyListener.Close()

	srv := mux.NewServer()
	defer srv.Close()
	srv.Handler = muxer

	ctx, cancelFn := context.WithCancel(context.TODO())
	defer cancelFn()
	go func() {
		_ = srv.Serve(loopbackLis)
	}()
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			conn, err := http1Listener.Accept()
			if err != nil {
				http1Result <- handleHTTP1Close
			} else {
				_, _ = conn.Write([]byte("http1Listener"))
				http1Result <- handleHTTP1Request
				select {
				case <-ctx.Done():
					break
				}
				_ = conn.Close()
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
			conn, err := anyListener.Accept()
			if err != nil {
				anyResult <- handleAnyClose
			} else {
				_, err = conn.Write([]byte("any"))
				anyResult <- handleAnyRequest
				select {
				case <-ctx.Done():
					break
				}
				_ = conn.Close()
			}
		}
	}()
	time.Sleep(testDuration) // wait to prevent timeouts on slow test-runners
	client, err := net.Dial("tcp", loopbackLis.Addr().String())
	if err != nil {
		log.Fatal("testTimeout client failed: ", err)
	}
	defer client.Close()
	time.Sleep(testDuration / 2)
	if len(anyResult) != 0 {
		log.Print("tcp ")
		t.Fatal("testTimeout failed: accepted to fast: ", len(anyResult))
	}
	//_ = client.SetReadDeadline(time.Now().Add(testDuration * 3))
	buffer := make([]byte, 10)
	rl, err := client.Read(buffer)
	if err != nil {
		t.Fatal("testTimeout failed: client error: ", err, rl)
	}
	_ = srv.Close()
	if rl != len("any") {
		log.Print("testTimeout failed: response from wrong service ", rl)
	}
	if string(buffer[0:3]) != "any" {
		log.Print("testTimeout failed: response from wrong service ")
	}
	time.Sleep(testDuration * 2)
	if len(http1Result) != 1 {
		t.Fatal("testTimeout failed: http1 accepted to less: ", len(http1Result))
	}
	if len(anyResult) != 1 {
		t.Fatal("testTimeout failed: any accepted to less: ", len(anyResult))
	}
	if a := <-http1Result; a != handleHTTP1Close {
		t.Fatal("testTimeout failed: no close an http rule")
	}
	if a := <-anyResult; a != handleAnyRequest {
		t.Fatal("testTimeout failed: anyListener rule did not match")
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

	muxer := mux.NewServeMux()
	// Register a bogus matcher to force buffering exactly the right amount.
	// Before this fix, this would trigger a bug where `Read` would incorrectly
	// report `io.EOF` when only the buffer had been consumed.
	_ = muxer.HandleListener(mux.MatcherFunc(func(w io.Writer, r io.Reader) bool {
		var b [len(payload)]byte
		_, _ = r.Read(b[:])
		return false
	}))
	anyl := muxer.HandleListener(mux.Any())

	l := newChanListener()
	l.Notify(reader)
	defer l.Close()
	srv := mux.NewServer()
	defer srv.Close()
	srv.Handler = muxer
	go safeServe(errCh, srv, l)
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
	errCh := make(chan error, 5)
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

	var wg sync.WaitGroup
	func() {
		muxer := mux.NewServeMux()
		httpl := muxer.HandleListener(mux.Any())

		srv := mux.NewServer()
		defer srv.Close()
		srv.Handler = muxer

		wg.Add(1)
		go func() {
			defer wg.Done()
			runTestHTTPServer(errCh, httpl)
		}()
		wg.Add(1)
		go func() {
			defer wg.Done()
			safeServe(errCh, srv, l)
		}()
		runTestHTTP1Client(t, l.Addr())
	}()
	wg.Wait()
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
	muxer := mux.NewServeMux()

	tlsl := muxer.HandleListener(mux.TLS())
	httpl := muxer.HandleListener(mux.Any())

	srv := mux.NewServer()
	defer srv.Close()
	srv.Handler = muxer

	go runTestTLSServer(errCh, tlsl)
	go runTestHTTPServer(errCh, httpl)
	go safeServe(errCh, srv, l)

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
	muxer := mux.NewServeMux()

	// Register a bogus matcher that only reads one byte.
	muxer.HandleListener(mux.MatcherFunc(func(w io.Writer, r io.Reader) bool {
		var b [1]byte
		_, _ = r.Read(b[:])
		return false
	}))
	h2l := muxer.HandleListener(mux.HTTP2())

	l := newChanListener()
	l.Notify(reader)
	srv := mux.NewServer()
	defer srv.Close()
	srv.Handler = muxer
	go safeServe(errCh, srv, l)
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
	testHTTP2HeaderField(t, mux.HTTP2HeaderFieldEqual, "value", "value", "anothervalue")
}

func TestHTTP2MatchHeaderFieldPrefix(t *testing.T) {
	testHTTP2HeaderField(t, mux.HTTP2HeaderFieldPrefix, "application/grpc+proto", "application/grpc", "application/json")
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

	muxer := mux.NewServeMux()

	httpl := muxer.HandleListener(mux.MatcherAny(mux.HTTP2(), mux.HTTP1Fast()))

	rpcl := muxer.HandleListener(mux.Any())

	srv := mux.NewServer()
	defer srv.Close()
	srv.Handler = muxer

	go runTestHTTPServer(errCh, httpl)
	go runTestRPCServer(errCh, rpcl)
	go safeServe(errCh, srv, l)

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
	muxer := mux.NewServeMux()

	srv := mux.NewServer()
	defer srv.Close()
	srv.Handler = muxer
	srv.Handler = muxer

	httpl := muxer.HandleListener(mux.MatcherAny(mux.HTTP2(), mux.HTTP1Fast()))

	go runTestHTTPServer(errCh, httpl)
	go safeServe(errCh, srv, l)

	var errCount uint32
	srv.HandleError(mux.ErrorHandlerFunc(func(err error) bool {
		if atomic.AddUint32(&errCount, 1) == 1 {
			t.Logf("got an expected error: %v", err)
		}
		return true
	}))

	//runTestRPCClient(t, l.Addr())
	c, clean := safeDial(t, l.Addr())
	defer clean()

	l.Close()

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

	muxer := mux.NewServeMux()
	srv := mux.NewServer()
	defer srv.Close()
	srv.Handler = muxer

	lis := muxer.HandleListener(mux.MatcherAny(mux.MatcherFunc(unmatcher), mux.MatcherFunc(matcher), mux.MatcherFunc(unmatcher)))

	go runTestHTTPServer(errCh, lis)
	go safeServe(errCh, srv, l)

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
	muxer := mux.NewServeMux()

	srv := mux.NewServer()
	defer srv.Close()
	srv.Handler = muxer

	anyl := muxer.HandleListener(mux.Any())

	go safeServe(errCh, srv, l)

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
		if err != mux.ErrListenerClosed {
			t.Fatal(err)
		}
		// The error is either io.ErrClosedPipe or net.OpError wrapping
		// a net.pipeError depending on the go version.
		if _, err := c2.Read([]byte{}); !strings.Contains(err.Error(), "closed") {
			t.Fatalf("connection is not closed and is leaked: %v", err)
		}
	}
}
