// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors_test

import (
	"fmt"
	"testing"

	"github.com/searKing/golang/go/error/multi"
	"github.com/searKing/golang/go/errors"
)

func TestMulti(t *testing.T) {
	tests := []struct {
		errs []error
		want error
	}{
		{nil, nil},
		{[]error{}, nil},
		{[]error{fmt.Errorf("foo")}, fmt.Errorf("foo")},
		{[]error{fmt.Errorf("foo"), fmt.Errorf("fun")}, multi.New(fmt.Errorf("foo"), fmt.Errorf("fun"))},
	}

	for _, tt := range tests {
		got := errors.Multi(tt.errs...)
		if got == tt.want {
			continue
		}
		if got == nil || tt.want == nil || got.Error() != tt.want.Error() {
			t.Errorf("Multi.Error(): got: %q, want %q", got, tt.want)
		}
	}
}
