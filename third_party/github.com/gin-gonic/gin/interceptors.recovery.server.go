// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gin

import (
	"io"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/searKing/golang/go/error/builtin"
)

//go:generate go-option -type "recovery"
type recovery struct {
	recoveryHandler func(c *gin.Context, err any) error
}

// Recovery returns a middleware that recovers from any panics and writes a 500 if there was one.
func Recovery(opts ...RecoveryOption) gin.HandlerFunc {
	return RecoveryWithWriter(gin.DefaultErrorWriter, opts...)
}

// RecoveryWithWriter returns a middleware for a given writer
// that recovers from any panics and writes a 500 if there was one.
func RecoveryWithWriter(out io.Writer, opts ...RecoveryOption) gin.HandlerFunc {
	var opt recovery
	opt.ApplyOptions(WithRecoveryHandler(RecoverHandler))
	opt.ApplyOptions(opts...)
	return func(c *gin.Context) {
		defer func() {
			builtin.Recover(out, func(err any) any {
				if opt.recoveryHandler != nil {
					return opt.recoveryHandler(c, err)
				}
				return err
			}, func() string {
				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				headers := strings.Split(string(httpRequest), "\r\n")
				for idx, header := range headers {
					current := strings.Split(header, ":")
					if current[0] == "Authorization" {
						headers[idx] = current[0] + ": *"
					}
				}
				return string(httpRequest)
			})
		}()
		c.Next() // execute all the handlers
	}
}

func RecoverHandler(c *gin.Context, err any) error {
	var brokenPipe = builtin.ErrorIsBrokenPipe(err)
	// If the connection is dead, we can't write a status to it.
	if brokenPipe {
		_ = c.Error(err.(error)) // nolint: errcheck
		c.Abort()
	} else {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	return nil
}
