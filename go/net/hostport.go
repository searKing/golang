package net

import (
	"fmt"
	"net"
)

func HostportOrDefault(hostport string, defHostport string) string {
	host, port, _ := SplitHostPort(hostport)
	defHost, defPort, _ := SplitHostPort(defHostport)
	if host == "" {
		host = defHost
	}
	if port == "" {
		port = defPort
	}
	return net.JoinHostPort(host, port)
}

// SplitHostPort splits a network address of the form "host:port",
// "host%zone:port", "[host]:port" or "[host%zone]:port" into host or
// host%zone and port.
// Different to net.SplitHostPort, host or port can be not set
//
// A literal IPv6 address in hostport must be enclosed in square
// brackets, as in "[::1]:80", "[::1%lo0]:80".
//
// See func Dial for a description of the hostport parameter, and host
// and port results.
func SplitHostPort(hostport string) (host, port string, err error) {
	host, portStr, err := net.SplitHostPort(hostport)
	if err != nil {
		// If adding a port makes it valid, the previous error
		// was not due to an invalid address and we can append a port.
		host, _, err := net.SplitHostPort(hostport + ":1234")
		return host, "", err
	}
	return host, portStr, err
}

func RandomPort(host string) (int, error) {
	if host == "" {
		host = "localhost"
	}
	ln, err := net.Listen("tcp", net.JoinHostPort(host, "0"))
	if err != nil {
		return 0, fmt.Errorf("could not generate random port: %w", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	err = ln.Close()
	if err != nil {
		return 0, fmt.Errorf("could not generate random port: %w", err)
	}
	return port, nil

}
