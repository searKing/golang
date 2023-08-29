// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !go1.20

package errors_test

import (
	"errors"
	"reflect"
	"testing"

	errors_ "github.com/searKing/golang/go/errors"
)

func TestMulti(t *testing.T) {
	err1 := errors.New("err1")
	err2 := errors.New("err2")
	for _, test := range []struct {
		errs []error
		want []error
	}{{
		errs: []error{err1},
		want: []error{err1},
	}, {
		errs: []error{err1, err2},
		want: []error{err1},
	}, {
		errs: []error{err1, nil, err2},
		want: []error{err1},
	}, {
		errs: []error{nil, err1, nil, err2},
		want: []error{err1},
	}} {
		me := errors_.Multi(test.errs...)
		var got []error
		if i, ok := me.(interface{ Unwrap() []error }); ok {
			got = i.Unwrap()
		} else if i, ok := me.(interface{ Unwrap() error }); ok {
			got = append(got, i.Unwrap())
		}
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("Multi(%v) = %v; want %v", test.errs, got, test.want)
		}
		if len(got) != cap(got) {
			t.Errorf("Multi(%v) returns errors with len=%v, cap=%v; want len==cap", test.errs, len(got), cap(got))
		}
	}
}
