// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http

import (
	"net/http"
	"strings"

	net_ "github.com/searKing/golang/go/net"
	strings_ "github.com/searKing/golang/go/strings"
)

// ClientIP implements a best effort algorithm to return the real client IP, it parses
// X-Real-IP and X-Forwarded-For in order to work properly with reverse-proxies such us: nginx or haproxy.
// Use X-Forwarded-For before X-Real-Ip as nginx uses X-Real-Ip with the proxy's IP.
func ClientIP(req *http.Request) string {
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-Forwarded-For
	// https://cloud.google.com/appengine/docs/flexible/python/reference/request-headers
	clientIP := strings_.ValueOrDefault(
		strings.TrimSpace(strings.Split(req.Header.Get("X-Forwarded-For"), ",")[0]),
		strings.TrimSpace(req.Header.Get("X-Real-Ip")),
		strings.TrimSpace(req.Header.Get("X-Appengine-Remote-Addr")))
	if clientIP != "" {
		return clientIP
	}

	if ip, _, err := net_.SplitHostPort(strings.TrimSpace(req.RemoteAddr)); err == nil {
		return ip
	}

	return ""
}
