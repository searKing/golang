// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package net

import (
	"errors"
	"net"
	"time"

	strings_ "github.com/searKing/golang/go/strings"
	"golang.org/x/net/nettest"
)

// This code is borrowed from https://github.com/uber/tchannel-go/blob/dev/localip.go

// ScoreAddr scores how likely the given addr is to be a remote address and returns the
// IP to use when listening. Any address which receives a negative score should not be used.
// Scores are calculated as:
// -1 for any unknown IP addreseses.
// +300 for IPv4 addresses
// +100 for non-local addresses, extra +100 for "up" interfaces.
// +100 for routable addresses
// -50 for local mac addr.
func ScoreAddr(iface net.Interface, addr net.Addr) (int, net.IP) {
	var ip net.IP
	if netAddr, ok := addr.(*net.IPNet); ok {
		ip = netAddr.IP
	} else if netIP, ok := addr.(*net.IPAddr); ok {
		ip = netIP.IP
	} else {
		return -1, nil
	}

	var score int
	if ip.To4() != nil {
		score += 300
	}
	if iface.Flags&net.FlagLoopback == 0 && !ip.IsLoopback() {
		score += 100
		if iface.Flags&net.FlagUp != 0 {
			score += 100
		}
	}
	_, routable := isRoutableIP("ip", ip)
	if routable {
		score -= 25
	}
	if isLocalMacAddr(iface.HardwareAddr) {
		score -= 50
	}
	return score, ip
}

// filter is a interface filter which returns false if the interface is _not_ to listen on
func listenAddr(interfaces []net.Interface, filter func(iface net.Interface) bool) (net.HardwareAddr, net.IP, error) {
	if filter == nil {
		filter = func(iface net.Interface) bool { return true }
	}
	bestScore := -1
	var bestIP net.IP
	var bestMac net.HardwareAddr
	// Select the highest scoring IP as the best IP.
	for _, iface := range interfaces {
		if !filter(iface) {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			// Skip this interface if there is an error.
			continue
		}

		for _, addr := range addrs {
			score, ip := ScoreAddr(iface, addr)
			if score > bestScore {
				bestScore = score
				bestIP = ip
				bestMac = iface.HardwareAddr
			}
		}
	}

	if bestScore == -1 {
		return nil, nil, errors.New("no addresses to listen on")
	}

	return bestMac, bestIP, nil
}

// ExpectInterfaceNameFilter
// If you want to listen specified interfaces (and the loopback) give the name of the interface (eg eth0) here.
func ExpectInterfaceNameFilter(names ...string) func(iface net.Interface) bool {
	return func(iface net.Interface) bool {
		if len(names) == 0 {
			return true
		}
		return strings_.SliceContainsAny(names, iface.Name)
	}
}

// ExceptInterfaceNameFilter
// you can specify which interface _not_ to listen on
func ExceptInterfaceNameFilter(names ...string) func(iface net.Interface) bool {
	return func(iface net.Interface) bool {
		if len(names) == 0 {
			return true
		}
		return !strings_.SliceContainsAny(names, iface.Name)
	}
}

// RoutedInterfaceNameFilter returns a network interface that can route IP
// traffic and satisfies flags.
//
// The provided network must be "ip", "ip4" or "ip6".
func RoutedInterfaceNameFilter() func(iface net.Interface) bool {
	return func(iface net.Interface) bool {
		rifs, err := nettest.RoutedInterface("ip", net.FlagUp|net.FlagBroadcast)
		if err != nil {
			return true
		}

		return rifs.Name == iface.Name
	}
}

// ListenIP returns the IP to bind to in Listen. It tries to find an IP that can be used
// by other machines to reach this machine.
// filters is interface filters any return false if the interface is _not_ to listen on
func ListenIP(filters ...func(iface net.Interface) bool) (net.IP, error) {
	_, ip, err := ListenAddr(filters...)
	return ip, err
}

// ListenMac returns the Mac to bind to in Listen. It tries to find an Mac that can be used
// by other machines to reach this machine.
// filters is interface filters any return false if the interface is _not_ to listen on
func ListenMac(filters ...func(iface net.Interface) bool) (net.HardwareAddr, error) {
	mac, _, err := ListenAddr(filters...)
	return mac, err
}

// ListenAddr returns the Mac and IP to bind to in Listen. It tries to find an Mac and IP that can be used
// by other machines to reach this machine.
// filters is interface filters any return false if the interface is _not_ to listen on
func ListenAddr(filters ...func(iface net.Interface) bool) (net.HardwareAddr, net.IP, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, nil, err
	}
	return listenAddr(interfaces, func(iface net.Interface) bool {
		for _, filter := range filters {
			if filter != nil && !filter(iface) {
				return false
			}
		}
		return true
	})
}

// DialIP returns the local IP to in Dial.
func DialIP(network, address string, timeout time.Duration) (net.IP, error) {
	conn, err := net.DialTimeout(network, address, timeout)
	if err != nil {
		return nil, err
	}
	a := conn.LocalAddr()
	ipAddr, err := net.ResolveIPAddr(a.Network(), a.String())
	defer conn.Close()
	if err != nil {
		return nil, err
	}
	return ipAddr.IP, nil
}

// ServeIP returns the IP to bind to in Listen. It tries to find an IP that can be used
// by other machines to reach this machine.
// Order is by DialIP and ListenIP
func ServeIP(networks, addresses []string, timeout time.Duration) (net.IP, error) {
	for _, network := range networks {
		for _, address := range addresses {
			ip, err := DialIP(network, address, timeout)
			if err != nil {
				continue
			}
			return ip, nil
		}
	}
	return ListenIP()
}

// If the first octet's second least-significant-bit is set, then it's local.
// https://en.wikipedia.org/wiki/MAC_address#Universal_vs._local
func isLocalMacAddr(addr net.HardwareAddr) bool {
	if len(addr) == 0 {
		return false
	}
	return addr[0]&2 == 2
}

func isRoutableIP(network string, ip net.IP) (net.IP, bool) {
	if !ip.IsLoopback() && !ip.IsLinkLocalUnicast() && !ip.IsGlobalUnicast() {
		return nil, false
	}
	switch network {
	case "ip4":
		if ip := ip.To4(); ip != nil {
			return ip, true
		}
	case "ip6":
		if ip.IsLoopback() { // addressing scope of the loopback address depends on each implementation
			return nil, false
		}
		if ip := ip.To16(); ip != nil && ip.To4() == nil {
			return ip, true
		}
	default:
		if ip := ip.To4(); ip != nil {
			return ip, true
		}
		if ip := ip.To16(); ip != nil {
			return ip, true
		}
	}
	return nil, false
}
