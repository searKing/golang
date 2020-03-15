// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http

import (
	"net/http"
	"net/url"
	"regexp"
)

var (
	reForwardedHost  = regexp.MustCompile(`host=([^,]+)`)
	reForwardedProto = regexp.MustCompile(`proto=(https?)`)
)

// GetProxySchemeAndHost extracts the host and used protocol (either HTTP or HTTPS)
// from the given request. If `allowForwarded` is set, the X-Forwarded-Host,
// X-Forwarded-Proto and Forwarded headers will also be checked to
// support proxies.
func GetProxySchemeAndHost(r *http.Request, allowForwarded bool) (scheme, host string) {
	if r == nil {
		return
	}
	if r.TLS != nil {
		scheme = "https"
	} else {
		scheme = "http"
	}

	host = r.Host

	if !allowForwarded {
		return
	}

	if h := r.Header.Get("X-Forwarded-Host"); h != "" {
		host = h
	}

	if h := r.Header.Get("X-Forwarded-Proto"); h == "http" || h == "https" {
		scheme = h
	}

	if h := r.Header.Get("Forwarded"); h != "" {
		if r := reForwardedHost.FindStringSubmatch(h); len(r) == 2 {
			host = r[1]
		}

		if r := reForwardedProto.FindStringSubmatch(h); len(r) == 2 {
			scheme = r[1]
		}
	}

	return
}

func ResolveProxyUrl(u *url.URL, r *http.Request, allowForwarded bool) *url.URL {
	if u == nil {
		return nil
	}
	scheme, host := GetProxySchemeAndHost(r, allowForwarded)
	u.Scheme = scheme
	u.Host = host
	return u
}
