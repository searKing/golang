// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package exception

type IllegalStateException struct {
	*RuntimeException
}

func NewIllegalStateException() *IllegalStateException {
	return &IllegalStateException{
		RuntimeException: NewRuntimeException(),
	}
}

func NewIllegalStateException1(message string) *IllegalStateException {
	return &IllegalStateException{
		RuntimeException: NewRuntimeException1(message),
	}
}

func NewIllegalStateException2(message string, cause Throwable) *IllegalStateException {
	return &IllegalStateException{
		RuntimeException: NewRuntimeException2(message, cause),
	}
}
