// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package recovery

import (
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/searKing/golang/go/error/builtin"
)

func ServerInterceptor(f func(c *gin.Context, err interface{})) gin.HandlerFunc {
	return ServerInterceptorWithWriter(gin.DefaultErrorWriter, f)
}

func ServerInterceptorWithWriter(out io.Writer, f func(c *gin.Context, err interface{})) gin.HandlerFunc {
	var logger *log.Logger
	if out != nil {
		logger = log.New(out, "\n\n\x1b[31m", log.LstdFlags)
	}

	return func(c *gin.Context) {
		defer func() {
			builtin.Recover(logger, func(err interface{}) {
				var brokenPipe = builtin.ErrorIsBrokenPipe(err)
				// If the connection is dead, we can't write a status to it.
				if brokenPipe {
					_ = c.Error(err.(error)) // nolint: errcheck
					c.Abort()
				} else {
					c.AbortWithStatus(http.StatusInternalServerError)
				}
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
