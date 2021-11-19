// Copyright 2021 The searKing Author. All rights reserved.
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

	ioutil_ "github.com/searKing/golang/go/io/ioutil"
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
func GlogPreRotate(name string) {
	// Write header.
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "Log file created at: %s by %s\n", time.Now().Format("2006/01/02 15:04:05"), program)
	fmt.Fprintf(&buf, "Running on machine: %s\n", host)
	fmt.Fprintf(&buf, "Binary: Built with %s %s for %s/%s\n", runtime.Compiler, runtime.Version(), runtime.GOOS, runtime.GOARCH)
	fmt.Fprintf(&buf, "Log line format: [IWEF]yyyymmdd hh:mm:ss.uuuuuu threadid file:line(func)] msg\n")
	ioutil_.AppendAll(name, buf.Bytes())
}
