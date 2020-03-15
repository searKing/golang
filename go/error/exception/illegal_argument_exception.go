// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package exception

type IllegalArgumentException struct {
	*RuntimeException
}

func NewIllegalArgumentException() *IllegalArgumentException {
	return &IllegalArgumentException{
		RuntimeException: NewRuntimeException(),
	}
}

func NewIllegalArgumentException1(message string) *IllegalArgumentException {
	return &IllegalArgumentException{
		RuntimeException: NewRuntimeException1(message),
	}
}

func NewIllegalArgumentException2(message string, cause Throwable) *IllegalArgumentException {
	return &IllegalArgumentException{
		RuntimeException: NewRuntimeException2(message, cause),
	}
}
