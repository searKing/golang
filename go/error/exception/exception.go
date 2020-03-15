// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package exception

type Exception struct {
	*ThrowableObject
}

func NewException() *Exception {
	return &Exception{
		ThrowableObject: NewThrowable(),
	}
}

func NewException1(message string) *Exception {
	return &Exception{
		ThrowableObject: NewThrowable1(message),
	}
}

func NewException2(message string, cause Throwable) *Exception {
	return &Exception{
		ThrowableObject: NewThrowable2(message, cause),
	}
}

func NewException4(message string, cause Throwable, enableSuppression, writableStackTrace bool) *Exception {
	return &Exception{
		ThrowableObject: NewThrowable4(message, cause, enableSuppression, writableStackTrace),
	}
}
