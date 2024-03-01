// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import (
	"cmp"
	"fmt"
	"io"
	"reflect"
)

// Optional represents a Value that may be null.
type Optional[E any] struct {
	Value E
	Valid bool // Valid is true if Value is not NULL
}

func (o Optional[E]) Format(s fmt.State, verb rune) {
	if o.Valid {
		_, _ = fmt.Fprintf(s, "%"+string(verb), o.Value)
		return
	}
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = fmt.Fprintf(s, "null %s: %+v", reflect.TypeOf(o.Value).String(), o.Value)
			return
		}
		fallthrough
	case 's', 'q':
		_, _ = io.WriteString(s, "null")
	}
}

func (o Optional[E]) String() string {
	if o.Valid {
		return fmt.Sprintf("%v", o.Value)
	}
	return "null"
}

// ValueOr returns the contained value if available, another value otherwise
func (o Optional[E]) ValueOr(e E) E {
	if o.Valid {
		return o.Value
	}
	return e
}

func CompareOptional[E cmp.Ordered](a, b Optional[E]) int {
	if a.Valid && !b.Valid {
		return 1
	}
	if !a.Valid && b.Valid {
		return -1
	}
	if !a.Valid && !b.Valid {
		return 0
	}

	if a.Value == b.Value {
		return 0
	}
	if a.Value < b.Value {
		return -1
	}
	return 1
}

func Opt[E any](e E) Optional[E] {
	return Optional[E]{
		Value: e,
		Valid: true,
	}
}

// NullOpt is the null Optional.
//
// Deprecated: Use a literal types.Optional[E]{} instead.
func NullOpt[E any]() Optional[E] {
	var zeroE E
	return Optional[E]{
		Value: zeroE,
		Valid: false,
	}
}
