// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mux_test

import (
	"bytes"
	"crypto/rand"
	"crypto/tls"
	"errors"
	"fmt"
	"go/build"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/exec"
	"strings"
	"sync"
	"testing"
	"time"

	net_ "github.com/searKing/golang/go/net"
	"github.com/searKing/golang/go/net/mux"
	"github.com/searKing/golang/go/sync/atomic"
	"github.com/searKing/golang/go/testing/leakcheck"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/hpack"
)

const (
	testHTTP1Resp = "http1"
	rpcVal        = 1234
)

func safeServe(errCh chan<- error, muxl *mux.Server, l net.Listener) {
	if err := muxl.Serve(l); err != nil {
		if err == mux.ErrServerClosed || err == mux.ErrListenerClosed {
			return
		}
		if strings.Contains(err.Error(), "use of closed") {
			return
		}
		errCh <- err
	}
}

func safeDial(t *testing.T, addr net.Addr) (*rpc.Client, func()) {
	c, err := rpc.Dial(addr.Network(), addr.String())
	if err != nil {
		t.Fatal(err)
	}
	return c, func() {
		if err := c.Close(); err != nil {
			t.Fatal(err)
		}
	}
}

type chanListener struct {
	net.Listener
	connCh     chan net.Conn
	inShutdown atomic.Bool
}

func newChanListener() *chanListener {
	return &chanListener{connCh: make(chan net.Conn, 1)}
}

func (l *chanListener) Notify(conn net.Conn) {
	if l.inShutdown.Load() {
		return
	}
	l.connCh <- conn
}

func (l *chanListener) Accept() (net.Conn, error) {
	if c, ok := <-l.connCh; ok {
		return c, nil
	}
	return nil, errors.New("use of closed network connection")
}

func (l *chanListener) Close() error {
	if l.inShutdown.Load() {
		return nil
	}

	l.inShutdown.Store(true)

	close(l.connCh)

	if l.Listener == nil {
		return nil
	}
	return l.Listener.Close()
}

func testListener(t leakcheck.Errorfer) net.Listener {
	l, err := net_.LoopbackListener()
	if err != nil {
		t.Errorf(err.Error())
	}
	return net_.OnceCloseListener(l)
}

type testHTTP1Handler struct{}

func (h *testHTTP1Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, testHTTP1Resp)
}

func runTestHTTPServer(errCh chan<- error, l net.Listener) {
	var mu sync.Mutex
	conns := make(map[net.Conn]struct{})

	defer func() {
		mu.Lock()
		for c := range conns {
			if err := c.Close(); err != nil {
				errCh <- err
			}
		}
		mu.Unlock()
	}()

	s := &http.Server{
		Handler: &testHTTP1Handler{},
		ConnState: func(c net.Conn, state http.ConnState) {
			mu.Lock()
			switch state {
			case http.StateNew:
				conns[c] = struct{}{}
			case http.StateClosed:
				delete(conns, c)
			}
			mu.Unlock()
		},
	}
	if err := s.Serve(l); err != mux.ErrListenerClosed {
		errCh <- err
	}
}

func generateTLSCert(t *testing.T) {
	err := exec.Command("go", "run", build.Default.GOROOT+"/src/crypto/tls/generate_cert.go", "--host", "*").Run()
	if err != nil {
		t.Fatal(err)
	}
}

func cleanupTLSCert(t *testing.T) {
	err := os.Remove("cert.pem")
	if err != nil {
		t.Error(err)
	}
	err = os.Remove("key.pem")
	if err != nil {
		t.Error(err)
	}
}

func runTestTLSServer(errCh chan<- error, l net.Listener) {
	certificate, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		errCh <- err
		log.Printf("1")
		return
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{certificate},
		Rand:         rand.Reader,
	}

	tlsl := tls.NewListener(l, config)
	runTestHTTPServer(errCh, tlsl)
}

func runTestHTTP1Client(t *testing.T, addr net.Addr) {
	runTestHTTPClient(t, "http", addr)
}

func runTestTLSClient(t *testing.T, addr net.Addr) {
	runTestHTTPClient(t, "https", addr)
}

func runTestHTTPClient(t *testing.T, proto string, addr net.Addr) {
	client := http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	r, err := client.Get(proto + "://" + addr.String())
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err = r.Body.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != testHTTP1Resp {
		t.Fatalf("invalid response: want=%s got=%s", testHTTP1Resp, b)
	}
}

type TestRPCRcvr struct{}

func (r TestRPCRcvr) Test(i int, j *int) error {
	*j = i
	return nil
}

func runTestRPCServer(errCh chan<- error, l net.Listener) {
	s := rpc.NewServer()
	if err := s.Register(TestRPCRcvr{}); err != nil {
		errCh <- err
	}
	for {
		c, err := l.Accept()
		if err != nil {
			if err != mux.ErrListenerClosed {
				errCh <- err
			}
			return
		}
		go s.ServeConn(c)
	}
}

func runTestRPCClient(t *testing.T, addr net.Addr) {
	c, clean := safeDial(t, addr)
	defer clean()

	var num int
	if err := c.Call("TestRPCRcvr.Test", rpcVal, &num); err != nil {
		t.Fatal(err)
	}

	if num != rpcVal {
		t.Errorf("wrong rpc response: want=%d got=%v", rpcVal, num)
	}
}

func testHTTP2HeaderField(
	t *testing.T,
	matcherConstructor func(sendSetting bool,
		expects ...hpack.HeaderField) mux.MatcherFunc,
	headerValue string,
	matchValue string,
	notMatchValue string,
) {
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
	name := "name"
	writer, reader := net.Pipe()
	go func() {
		if _, err := io.WriteString(writer, http2.ClientPreface); err != nil {
			t.Fatal(err)
		}
		var buf bytes.Buffer
		enc := hpack.NewEncoder(&buf)
		if err := enc.WriteField(hpack.HeaderField{Name: name, Value: headerValue}); err != nil {
			t.Fatal(err)
		}
		framer := http2.NewFramer(writer, nil)
		if err := framer.WriteSettingsAck(); err != nil {
			t.Fatal(err)
		}

		if err := framer.WriteHeaders(http2.HeadersFrameParam{
			StreamID:      1,
			BlockFragment: buf.Bytes(),
			EndStream:     true,
			EndHeaders:    true,
		}); err != nil {
			t.Fatal(err)
		}
		if err := writer.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	l := newChanListener()
	l.Notify(reader)
	// Register a bogus matcher that only reads one byte.
	muxl := mux.HandleListener(mux.MatcherFunc(func(w io.Writer, r io.Reader) bool {
		var b [1]byte
		_, _ = r.Read(b[:])
		return false
	}))
	defer muxl.Close()

	// Create a matcher that cannot match the response.
	//muxl.Match(matcherConstructor(false, hpack.HeaderField{Name: name, Value: notMatchValue}))
	// Then match with the expected field.
	h2l := mux.HandleListener(matcherConstructor(false, hpack.HeaderField{Name: name, Value: matchValue}))
	defer h2l.Close()

	srv := mux.NewServer()
	go func() {
		safeServe(errCh, srv, l)
	}()
	muxedConn, err := h2l.Accept()
	_ = l.Close()
	if err != nil {
		t.Fatal(err)
	}
	var b [len(http2.ClientPreface)]byte
	// We have the sniffed buffer first...
	if _, err := muxedConn.Read(b[:]); err == io.EOF {
		t.Fatal(err)
	}
	if string(b[:]) != http2.ClientPreface {
		t.Errorf("got unexpected read %s, expected %s", b, http2.ClientPreface)
	}
}
