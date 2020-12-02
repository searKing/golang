// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package net

import (
	"fmt"
	"net"

	"github.com/searKing/golang/go/errors"
)

var Ipv4LoopbackHosts = []string{"localhost", "127.0.0.1"}
var Ipv6LoopbackHosts = []string{"localhost", "[::1]", "::1"}

func IsLoopbackHost(host string) bool {
	for _, v := range Ipv4LoopbackHosts {
		if host == v {
			return true
		}
	}
	for _, v := range Ipv6LoopbackHosts {
		if host == v {
			return true
		}
	}
	return false
}

// LoopbackListener returns a loopback listener on a first usable port in ports or 0 if ports is empty
func LoopbackListener(ports ...string) (net.Listener, error) {
	var errs []error
	if len(ports) == 0 {
		ports = append(ports, "0")
	}

	for _, port := range ports {
		for _, host := range Ipv4LoopbackHosts {
			l, err := net.Listen("tcp", net.JoinHostPort(host, port))
			if err == nil {
				errs = append(errs, err)
				return l, nil
			}
		}
		for _, host := range Ipv6LoopbackHosts {
			l, err := net.Listen("tcp6", net.JoinHostPort(host, port))
			if err == nil {
				errs = append(errs, err)
				return l, nil
			}
		}
	}

	return nil, fmt.Errorf("net: failed to listen on ports: %v", errors.Multi(errs...))
}
