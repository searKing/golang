// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package exception

type IndexOutOfBoundsException struct {
	*RuntimeException
}

func NewIndexOutOfBoundsException() *IndexOutOfBoundsException {
	return &IndexOutOfBoundsException{
		RuntimeException: NewRuntimeException(),
	}
}

func NewIndexOutOfBoundsException1(message string) *IndexOutOfBoundsException {
	return &IndexOutOfBoundsException{
		RuntimeException: NewRuntimeException1(message),
	}
}

func NewIndexOutOfBoundsException2(message string, cause Throwable) *IndexOutOfBoundsException {
	return &IndexOutOfBoundsException{
		RuntimeException: NewRuntimeException2(message, cause),
	}
}
