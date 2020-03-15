// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cause_test

import (
	"fmt"
	"io"
	"testing"

	"github.com/searKing/golang/go/error/cause"
)

func TestWithError(t *testing.T) {
	tests := []struct {
		cause error
		err   error
		want  string
	}{
		{nil, nil, ""},
		{io.EOF, nil, ""},
		{nil, io.EOF, io.EOF.Error()},
		{io.EOF, fmt.Errorf("read error"), "read error: EOF"},
		{cause.WithError(io.EOF, fmt.Errorf("read error")), fmt.Errorf("client error"), "client error: read error: EOF"},
	}

	for _, tt := range tests {
		var got string
		err := cause.WithError(tt.cause, tt.err)
		if err != nil {
			got = err.Error()
		}
		if got != tt.want {
			t.Errorf("WithError(%v, %q): got: %q, want %q", tt.cause, tt.err, got, tt.want)
		}
	}
}
