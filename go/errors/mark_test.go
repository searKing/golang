// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors_test

import (
	"errors"
	"reflect"
	"testing"

	errors_ "github.com/searKing/golang/go/errors"
)

func TestMarkReturnsNil(t *testing.T) {
	if err := errors_.Mark(nil); err != nil {
		t.Errorf("errors_.Mark(nil) = %v, want nil", err)
	}
	if err := errors_.Mark(nil, nil); err != nil {
		t.Errorf("errors_.Mark(nil, nil) = %v, want nil", err)
	}
	mark := errors.New("mark")
	if err := errors_.Mark(nil, mark); err != nil {
		t.Errorf("errors_.Mark(nil, %v) = %v, want nil", mark, err)
	}
}

func TestMark(t *testing.T) {
	err := errors.New("err")
	mark1 := errors.New("mark1")
	mark2 := errors.New("mark2")
	for _, test := range []struct {
		err  error
		errs []error
		want error
	}{{
		err:  err,
		errs: []error{mark1},
		want: err,
	}, {
		err:  err,
		errs: []error{mark1, mark2},
		want: err,
	}, {
		err:  err,
		errs: []error{mark1, nil, mark2},
		want: err,
	}} {
		got := errors_.Mark(test.err, test.errs...).(interface{ Unwrap() error }).Unwrap()
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("Mark(%v) = %v; want %v", test.errs, got, test.want)
		}
	}
}

func TestMarkErrorMethod(t *testing.T) {
	err := errors.New("err")
	mark1 := errors.New("mark1")
	mark2 := errors.New("mark2")
	for _, test := range []struct {
		err  error
		errs []error
		want string
	}{{
		err:  err,
		errs: []error{mark1},
		want: "err",
	}, {
		err:  err,
		errs: []error{mark1, mark2},
		want: "err",
	}, {
		err:  err,
		errs: []error{mark1, nil, mark2},
		want: "err",
	}} {
		got := errors_.Mark(test.err, test.errs...).Error()
		if got != test.want {
			t.Errorf("Mark(%v).Error() = %q; want %q", test.errs, got, test.want)
		}
	}
}

func TestMarkErrorIs(t *testing.T) {
	err := errors.New("err")
	mark1 := errors.New("mark1")
	mark2 := errors.New("mark2")
	for _, test := range []struct {
		err     error
		errs    []error
		want    []error
		notWant []error
	}{{
		err:     nil,
		errs:    []error{mark1},
		want:    []error{nil},
		notWant: []error{err, mark1, mark2},
	}, {
		err:     err,
		errs:    []error{mark1},
		want:    []error{err, mark1},
		notWant: []error{nil, mark2},
	}, {
		err:     err,
		errs:    []error{mark1, mark2},
		want:    []error{err, mark1, mark2},
		notWant: []error{nil},
	}, {
		err:     err,
		errs:    []error{mark1, nil, mark2},
		want:    []error{err, mark1, mark2},
		notWant: []error{nil},
	}} {
		got := errors_.Mark(test.err, test.errs...)
		for _, want := range test.want {
			if !errors.Is(got, want) {
				t.Errorf("errors.Is(Mark(%v), %v) = %t; want %t", test.errs, want, got, true)
			}
		}
	}
}
