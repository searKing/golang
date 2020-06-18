// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gin

import "github.com/gin-gonic/gin"

var NopHandlerFunc = gin.HandlerFunc(func(c *gin.Context) {
	if c == nil {
		return
	}
	c.Next()
})
