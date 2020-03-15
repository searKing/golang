// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package url

import (
	"fmt"
	"net/url"
	"strconv"
)

func ParseStandardURL(rawUrl string) (*url.URL, error) {

	u, err := url.Parse(rawUrl) // Just url_.Parse (url_ is shadowed for godoc).
	if err != nil {
		return nil, err
	}
	// The host's colon:port should be normalized. See Issue 14836.
	u.Host = removeEmptyPort(u.Host)

	return u, nil
}

var portMap = map[string]string{
	"http":   "80",
	"https":  "443",
	"socks5": "1080",
	"stun":   "3478",
	"turn":   "3478",
	"stuns":  "5349",
	"turns":  "5349",
}

func ParseURL(rawUrl string, getDefaultPort func(schema string) (string, error)) (u *url.URL, host string, port int, err error) {
	if getDefaultPort == nil {
		getDefaultPort = func(schema string) (string, error) {
			port, ok := portMap[schema]
			if ok {
				return port, nil
			}
			return "", fmt.Errorf("malformed schema:%s", schema)
		}
	}

	standUrl, err := ParseStandardURL(rawUrl)
	if err != nil {
		return nil, "", -1, err
	}
	var hostport string
	if standUrl.Opaque == "" {
		hostport = standUrl.Host
	} else {
		hostport = standUrl.Opaque
	}
	rawHost, rawPort, err := ParseHostPort(standUrl.Scheme, hostport, getDefaultPort)
	if err != nil {
		return nil, "", -1, err
	}
	host = rawHost
	if port, err = strconv.Atoi(rawPort); err != nil {
		return nil, "", -1, fmt.Errorf("malformed port:%s", rawPort)
	}
	return standUrl, host, port, nil
}
