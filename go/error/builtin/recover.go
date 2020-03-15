// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package builtin

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

func Recover(logger *log.Logger, recoverHandler func(err interface{}), brokenPipeMsg func() string) {
	if err := recover(); err != nil {
		// Check for a broken connection, as it is not really a
		// condition that warrants a panic stack trace.
		var brokenPipe = ErrorIsBrokenPipe(err)
		if logger != nil {
			reset := string([]byte{27, 91, 48, 109})

			goErr := fmt.Errorf("panic %v", err)

			var msg string
			if brokenPipeMsg != nil {
				msg = brokenPipeMsg()
			}

			if brokenPipe {
				logger.Printf("[Recovery] brokenPipe %+v\n%s%s", goErr, msg, reset)
			} else {
				logger.Printf("[Recovery] %s panic recovered:\n%+v%s",
					timeFormat(time.Now()), goErr, reset)
			}
		}
		if recoverHandler != nil {
			recoverHandler(err)
		}
	}
}
func timeFormat(t time.Time) string {
	var timeString = t.Format("2006/01/02 - 15:04:05")
	return timeString
}

func ErrorIsBrokenPipe(err interface{}) bool {
	var brokenPipe bool
	if ne, ok := err.(*net.OpError); ok {
		if se, ok := ne.Err.(*os.SyscallError); ok {
			if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
				brokenPipe = true
			}
		}
	}
	return brokenPipe
}
