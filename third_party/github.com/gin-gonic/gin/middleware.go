// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gin

import (
	"github.com/gin-gonic/gin"
)

func UseHTTPPreflight() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Allow overriding the HTTP method. The reason for this is
		// that some libraries/environments to not support PATCH and
		// DELETE requests, e.g. Flash in a browser and parts of Java
		if newMethod := c.GetHeader("X-HTTP-Method-Override"); newMethod != "" {
			c.Request.Method = newMethod
		}

		// Add nosniff to all responses https://golang.org/src/net/http/server.go#L1429
		c.Header("X-Content-Type-Options", "nosniff")
	}
}
