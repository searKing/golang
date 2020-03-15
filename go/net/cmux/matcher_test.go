// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmux_test

import (
	"io"
	"io/ioutil"
	"net"
	"strings"
	"sync"
	"testing"

	"github.com/searKing/golang/go/net/cmux"
	"github.com/searKing/golang/go/testing/leakcheck"
)

func TestHTTP1Fast(t *testing.T) {
	defer leakcheck.Check(t)
	const payload = "GET /version HTTP/1.1\r\n"
	const mult = 2

	test(t, cmux.HTTP1Fast(), payload, mult)
}

func TestHTTP1(t *testing.T) {
	defer leakcheck.Check(t)
	const payload = "GET /version HTTP/1.1\r\n"
	const mult = 2

	test(t, cmux.HTTP1(), payload, mult)
}

func test(t *testing.T, matcher cmux.Matcher, payload string, mult int) {
	errCh := make(chan error)
	defer func() {
		select {
		case err := <-errCh:
			t.Fatal(err)
		default:
		}
	}()

	writer, reader := net.Pipe()
	defer reader.Close()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if _, err := io.WriteString(writer, strings.Repeat(payload, mult)); err != nil {
			t.Fatal(err)
		}
		_ = writer.Close()

	}()
	if !matcher.Match(nil, reader) {
		t.Errorf("expect false but accept true")
	}
	_, _ = ioutil.ReadAll(reader)

	wg.Wait()
}
