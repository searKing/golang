// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logrus

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	os_ "github.com/searKing/golang/go/os"
)

var (
	host    = "unknownhost"
	program = filepath.Base(os.Args[0])
)

func init() {
	h, err := os.Hostname()
	if err == nil {
		host = shortHostname(h)
	}
}

// shortHostname returns its argument, truncating at the first period.
// For instance, given "www.google.com" it returns "www".
func shortHostname(hostname string) string {
	if i := strings.Index(hostname, "."); i >= 0 {
		return hostname[:i]
	}
	return hostname
}

// GlogRotateHeader append rotate header to a file named by filename.
func GlogRotateHeader(name string) {
	// Write header.
	var buf bytes.Buffer
	_, _ = fmt.Fprintf(&buf, "Log file created at: %s by %s\n", time.Now().Format("2006/01/02 15:04:05"), program)
	_, _ = fmt.Fprintf(&buf, "Running on machine: %s\n", host)
	_, _ = fmt.Fprintf(&buf, "Binary: Built with %s %s for %s/%s\n", runtime.Compiler, runtime.Version(), runtime.GOOS, runtime.GOARCH)
	_, _ = fmt.Fprintf(&buf, "Log line format: [IWEF]yyyymmdd hh:mm:ss.uuuuuu threadid file:line(func)] msg\n")
	_ = os_.AppendAll(name, buf.Bytes())
}
