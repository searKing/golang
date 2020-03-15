// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package exception_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/searKing/golang/go/error/exception"
)

type ThrowableTests struct {
	input  exception.Throwable
	output []string
}

var (
	throwableTests = []ThrowableTests{
		{
			exception.NewThrowable(),
			[]string{"runtime/debug.Stack"},
		},
		{
			exception.NewThrowable1("throwable exception"),
			[]string{"throwable exception"},
		},
		{
			exception.NewThrowable2("throwable exception2", exception.NewThrowable1("throwable exception1")),
			[]string{"throwable exception1", "throwable exception2"},
		},
	}
)

func TestThrowable(t *testing.T) {
	for m, test := range throwableTests {
		runs := bytes.Buffer{}
		test.input.PrintStackTrace1(&runs)
		output := runs.String()
		for n, outputExpect := range test.output {
			if !strings.Contains(runs.String(), output) {
				t.Errorf("#%d[%d]: %v: got %s runs; expected %s", m, n, test.input, output, outputExpect)
				continue
			}
		}
	}

}
