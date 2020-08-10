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

func Recovery(f func(c *gin.Context, err interface{})) gin.HandlerFunc {
	return RecoveryWithWriter(gin.DefaultErrorWriter, f)
}

func RecoveryWithWriter(out io.Writer, f func(c *gin.Context, err interface{})) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			builtin.Recover(out, func(err interface{}) interface{} {
				var brokenPipe = builtin.ErrorIsBrokenPipe(err)
				// If the connection is dead, we can't write a status to it.
				if brokenPipe {
					_ = c.Error(err.(error)) // nolint: errcheck
					c.Abort()
				} else {
					c.AbortWithStatus(http.StatusInternalServerError)
				}
				return nil
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
