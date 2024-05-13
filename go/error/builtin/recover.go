// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package builtin

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"
)

// The Recover function allows a program to manage behavior of a
// panicking goroutine.
// affect as recoverHandler(recover())
// brokenPipeMsg is called when panic for broken pipe
func Recover(writer io.Writer, recoverHandler func(err any) any, brokenPipeMsg func() string) any {
	if err := recover(); err != nil {
		// Check for a broken connection, as it is not really a
		// condition that warrants a panic stack trace.
		var brokenPipe = ErrorIsBrokenPipe(err)
		if writer != nil {
			goErr := fmt.Errorf("panic %v", err)

			var msg string
			if brokenPipeMsg != nil {
				msg = brokenPipeMsg()
			}

			if brokenPipe {
				_, _ = writer.Write([]byte(fmt.Sprintf("[Recovery] brokenPipe %+v\n%s", goErr, msg)))
			} else {
				_, _ = writer.Write([]byte(fmt.Sprintf("[Recovery] %s panic recovered:\n%+v",
					timeFormat(time.Now()), goErr)))
			}
		}
		if recoverHandler != nil {
			return recoverHandler(err)
		}
		return err
	}
	return nil
}
func timeFormat(t time.Time) string {
	var timeString = t.Format("2006/01/02 - 15:04:05")
	return timeString
}

func ErrorIsBrokenPipe(err any) bool {
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
