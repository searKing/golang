// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testing

import (
	"strings"
	"testing"
)

// Assert asserts that the predicate is true, and formats its arguments
// using default formatting, analogous to [fmt.Println],
func Assert[T testing.TB](t T, p bool, args ...any) {
	if !p {
		var sb strings.Builder
		sb.WriteString("assertion failed")
		if len(args) > 0 {
			sb.WriteString(":")
			args = append([]any{sb.String()}, args...)
		}
		t.Helper()
		t.Fatal(args...)
	}
}

// Assertf asserts that the predicate is true, and formats its arguments
// according to the format, analogous to [fmt.Printf].
func Assertf[T testing.TB](t T, p bool, format string, args ...any) {
	if !p {
		var sb strings.Builder
		sb.WriteString("assertion failed")
		if format != "" {
			sb.WriteString(": ")
			sb.WriteString(format)
		}
		t.Helper()
		t.Fatalf(sb.String(), args...)
	}
}
