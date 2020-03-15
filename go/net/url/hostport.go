// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package url

import (
	"errors"
	"fmt"
	"net"
)

// ParseHostPort takes the user input target string, returns formatted host and port info.
// If target doesn't specify a port, set the port to be the defaultPort.
// If target is in IPv6 format and host-name is enclosed in sqarue brackets, brackets
// are strippd when setting the host.
// examples:
// target: "www.google.com" returns host: "www.google.com", port: "443"
// target: "ipv4-host:80" returns host: "ipv4-host", port: "80"
// target: "[ipv6-host]" returns host: "ipv6-host", port: "443"
// target: ":80" returns host: "localhost", port: "80"
// target: ":" returns host: "localhost", port: "443"
func ParseHostPort(scheme string, hostport string, getDefaultPort func(schema string) (string, error)) (host, port string, err error) {
	if hostport == "" {
		return "", "", errors.New("missing hostport")
	}
	if getDefaultPort == nil {
		return "", "", errors.New("missing getDefaultPort")
	}

	// justIP, no port follow
	if ip := net.ParseIP(hostport); ip != nil {
		// hostport is an IPv4 or IPv6(without brackets) address
		port, err := getDefaultPort(scheme)
		if err != nil {
			return "", "", err
		}
		return hostport, port, nil
	}

	if host, port, err = net.SplitHostPort(hostport); err == nil {
		// hostport has port, i.e ipv4-host:port, [ipv6-host]:port, host-name:port
		if host == "" {
			// Keep consistent with net.Dial(): If the host is empty, as in ":80", the local system is assumed.
			host = "localhost"
		}
		if port == "" {
			// If the port field is empty(hostport ends with colon), e.g. "[::1]:", defaultPort is used.
			port, err = getDefaultPort(scheme)
			if err != nil {
				return "", "", err
			}
		}
		return host, port, nil
	}
	// missing port
	defaultPort, err := getDefaultPort(scheme)
	if err != nil {
		return "", "", err
	}
	if host, port, err = net.SplitHostPort(hostport + ":" + defaultPort); err == nil {
		// hostport doesn't have port
		return host, port, nil
	}
	return "", "", fmt.Errorf("invalid hostport address %v", hostport)
}
