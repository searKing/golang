// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors_test

import (
	"errors"
	"testing"

	errors_ "github.com/searKing/golang/go/errors"
)

func TestMultiReturnsNil(t *testing.T) {
	if err := errors_.Multi(); err != nil {
		t.Errorf("errors_.Multi() = %v, want nil", err)
	}
	if err := errors_.Multi(nil); err != nil {
		t.Errorf("errors_.Multi(nil) = %v, want nil", err)
	}
	if err := errors_.Multi(nil, nil); err != nil {
		t.Errorf("errors_.Multi(nil, nil) = %v, want nil", err)
	}
}

func TestMultiErrorMethod(t *testing.T) {
	err1 := errors.New("err1")
	err2 := errors.New("err2")
	for _, test := range []struct {
		errs []error
		want string
	}{{
		errs: []error{err1},
		want: "err1",
	}, {
		errs: []error{err1, err2},
		want: "err1\nerr2",
	}, {
		errs: []error{err1, nil, err2},
		want: "err1\nerr2",
	}} {
		got := errors_.Multi(test.errs...).Error()
		if got != test.want {
			t.Errorf("Multi(%v).Error() = %q; want %q", test.errs, got, test.want)
		}
	}
}
