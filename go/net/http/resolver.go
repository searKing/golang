// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http

import (
	"net/http"

	url_ "github.com/searKing/golang/go/net/url"
)

// RequestWithTarget reset Host in url.Url by resolver.Target
func RequestWithTarget(req *http.Request, target string) {
	url_.ResolveWithTarget(req.Context(), req.URL, target)
}
