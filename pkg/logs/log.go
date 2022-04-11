// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logs

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/searKing/golang/go/errors"
	"github.com/searKing/golang/go/runtime"
	"github.com/sirupsen/logrus"
)

func init() {
	runtime.LogPanic.AppendHandler(logPanic)
	runtime.NeverPanicButLog.AppendHandler(logPanic)

	errors.ErrorHandlers = append(errors.ErrorHandlers, func(err error) {
		if err != nil {
			caller, file, line := runtime.GetShortCallerFuncFileLine(2)
			logrus.Errorf("Observed an error: %s at %s() %s:%d", err, caller, file, line)
		}
	})
}

// InitLog initializes logs the way we want for kubernetes.
func InitLog() {
	log.SetPrefix(fmt.Sprintf("[%s] ", os.Args[0]))
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	//log.SetFlags(0)
}

// logPanic logs the caller tree when a panic occurs (except in the special case of http.ErrAbortHandler).
func logPanic(r interface{}) {
	if r == nil || r == http.ErrAbortHandler {
		// honor the http.ErrAbortHandler sentinel panic value:
		//   ErrAbortHandler is a sentinel panic value to abort a handler.
		//   While any panic from ServeHTTP aborts the response to the client,
		//   panicking with ErrAbortHandler also suppresses logging of a stack trace to the server's error log.
		return
	}

	const size = 64 << 10
	stacktrace := runtime.GetCallStack(size)
	if _, ok := r.(string); ok {
		logrus.Errorf("Observed a panic: %s\n%s", r, stacktrace)
	} else {
		logrus.Errorf("Observed a panic: %#v (%v)\n%s", r, r, stacktrace)
	}
}
