// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http

import (
	"net"
	"net/http"
	"strings"

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
	return ipFromHostPort(req.RemoteAddr)
}

// ServerIP implements a best effort algorithm to return the real server IP, it parses
// LocalAddrContextKey from request context to get server IP.
func ServerIP(req *http.Request) string {
	if addr, ok := req.Context().Value(http.LocalAddrContextKey).(net.Addr); ok {
		return ipFromHostPort(addr.String())
	}
	return ""
}

func ipFromHostPort(hp string) string {
	h, _, err := net.SplitHostPort(hp)
	if err != nil {
		return ""
	}
	if len(h) > 0 && h[0] == '[' {
		return h[1 : len(h)-1]
	}
	return h
}
