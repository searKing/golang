// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package exception

type NullPointerException struct {
	*Exception
}

func NewNullPointerException() Throwable {
	return &NullPointerException{
		Exception: NewException(),
	}
}

func NewNullPointerException1(message string) Throwable {
	return &NullPointerException{
		Exception: NewException1(message),
	}
}
