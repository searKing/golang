// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mux

import "log"

// WithMaxIdleConns controls the maximum number of idle (keep-alive)
// connections across all hosts. Zero means no limit.
func WithMaxIdleConns(maxIdleConns int) ServerOption {
	return ServerOptionFunc(func(c *Server) {
		c.maxIdleConns = maxIdleConns
	})
}

// WithErrorLog specifies an optional logger for errors
func WithErrorLog(errorLog *log.Logger) ServerOption {
	return ServerOptionFunc(func(c *Server) {
		c.errorLog = errorLog
	})
}
