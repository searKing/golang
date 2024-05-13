// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package multi_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/searKing/golang/go/error/multi"
)

func TestNew(t *testing.T) {
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
		got := multi.New(tt.errs...)
		if errors.Is(got, tt.want) {
			continue
		}
		if got == nil || tt.want == nil || got.Error() != tt.want.Error() {
			t.Errorf("New.Error(): got: %q, want %q", got, tt.want)
		}
	}
}
