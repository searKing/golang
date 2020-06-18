// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/urfave/negroni"

	negroni_ "github.com/searKing/golang/third_party/github.com/urfave/negroni"
)

// WrapGinF is a helper function for wrapping gin.HandlerFunc
// Returns a negroni middleware
func UseNegroni(n *negroni.Negroni) gin.HandlerFunc {
	return func(c *gin.Context) {
		negroni_.Clone(n).With(negroni.WrapFunc(func(w http.ResponseWriter, r *http.Request) {
			c.Next()
		})).ServeHTTP(c.Writer, c.Request)
	}
}

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
