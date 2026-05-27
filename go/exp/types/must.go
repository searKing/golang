// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

// Must is a helper that wraps a call to a function returning (T, error)
// and panics if the error is non-nil. It is intended for use in variable
// initializations such as
//
//	var t = types.Must(template.New("name").Parse("text"))
//	var t = types.Must(os.Open("file.txt"))
func Must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}
