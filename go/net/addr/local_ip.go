// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package addr

import (
	"net"
	"time"

	net_ "github.com/searKing/golang/go/net"
)

// Deprecated: Use net_.ScoreAddr instead.
func ScoreAddr(iface net.Interface, addr net.Addr) (int, net.IP) {
	return net_.ScoreAddr(iface, addr)
}

// Deprecated: Use net_.ExpectInterfaceNameFilter instead.
func ExpectInterfaceNameFilter(names ...string) func(iface net.Interface) bool {
	return net_.ExpectInterfaceNameFilter(names...)
}

// Deprecated: Use net_.ExceptInterfaceNameFilter instead.
func ExceptInterfaceNameFilter(names ...string) func(iface net.Interface) bool {
	return net_.ExceptInterfaceNameFilter(names...)
}

// Deprecated: Use net_.RoutedInterfaceNameFilter instead.
func RoutedInterfaceNameFilter() func(iface net.Interface) bool {
	return net_.RoutedInterfaceNameFilter()
}

// Deprecated: Use net_.ListenIP instead.
func ListenIP(filters ...func(iface net.Interface) bool) (net.IP, error) {
	return net_.ListenIP(filters...)
}

// Deprecated: Use net_.ListenMac instead.
func ListenMac(filters ...func(iface net.Interface) bool) (net.HardwareAddr, error) {
	return net_.ListenMac(filters...)
}

// Deprecated: Use net_.ListenAddr instead.
func ListenAddr(filters ...func(iface net.Interface) bool) (net.HardwareAddr, net.IP, error) {
	return net_.ListenAddr(filters...)
}

// Deprecated: Use net_.DialIP instead.
func DialIP(network, address string, timeout time.Duration) (net.IP, error) {
	return net_.DialIP(network, address, timeout)
}

// Deprecated: Use net_.ServeIP instead.
func ServeIP(networks, addresses []string, timeout time.Duration) (net.IP, error) {
	return net_.ServeIP(networks, addresses, timeout)
}
