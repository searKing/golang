// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package net

import (
	"net"
	"strings"
)

var (
	_ net.Addr = strAddr("")
	_ net.Addr = multiAddrs{}
)

// strAddr is a net.Addr backed by either a TCP "ip:port" string, or
// the empty string if unknown.
type strAddr string

func (a strAddr) Network() string {
	if a != "" {
		// Per the documentation on net/http.Request.RemoteAddr, if this is
		// set, it's set to the IP:port of the peer (hence, TCP):
		// https://golang.org/pkg/net/http/#Request
		//
		// If we want to support Unix sockets later, we can
		// add our own grpc-specific convention within the
		// grpc codebase to set RemoteAddr to a different
		// format, or probably better: we can attach it to the
		// context and use that from serverHandlerTransport.RemoteAddr.
		return "tcp"
	}
	return ""
}

func (a strAddr) String() string { return string(a) }

type multiAddrs []net.Listener

// Network returns the address's network name, "tcp,udp".
func (ls multiAddrs) Network() string {
	var networkStrs []string
	for _, ln := range ls {
		networkStrs = append(networkStrs, ln.Addr().Network())
	}
	return strings.Join(networkStrs, ",")
}

func (ls multiAddrs) String() string {
	var addrStrs []string
	for _, l := range ls {
		addrStrs = append(addrStrs, l.Addr().String())
	}
	return strings.Join(addrStrs, ",")
}
