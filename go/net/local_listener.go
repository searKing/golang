package net

import (
	"fmt"
	"net"
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

// LoopbackListener returns a loopback listener on port or 0 if port is empty
func LoopbackListener(port string) net.Listener {
	var err error
	var l net.Listener
	if len(port) == 0 {
		port = "0"
	}

	for _, host := range Ipv4LoopbackHosts {
		l, err = net.Listen("tcp", net.JoinHostPort(host, port))
		if err == nil {
			return l
		}
	}
	for _, host := range Ipv6LoopbackHosts {
		l, err = net.Listen("tcp6", net.JoinHostPort(host, port))
		if err == nil {
			return l
		}
	}

	panic(fmt.Sprintf("net: failed to listen on a port: %v", err))
}
