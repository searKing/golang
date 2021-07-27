package net

import (
	"fmt"
	"net"
)

// HostportOrDefault takes the user input target string and default port, returns formatted host and port info.
// If target doesn't specify a port, set the port to be the defaultPort.
// If target is in IPv6 format and host-name is enclosed in square brackets, brackets
// are stripped when setting the host.
// examples:
// target: "www.google.com" defaultPort: "443" returns host: "www.google.com", port: "443"
// target: "ipv4-host:80" defaultPort: "443" returns host: "ipv4-host", port: "80"
// target: "[ipv6-host]" defaultPort: "443" returns host: "ipv6-host", port: "443"
// target: ":80" defaultPort: "443" returns host: "localhost", port: "80"
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

// RandomPort returns a random port by a temporary listen on :0
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
