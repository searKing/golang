// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package addr

import (
	"errors"
	"net"
	"time"
)

// This code is borrowed from https://github.com/uber/tchannel-go/blob/dev/localip.go

// This code is borrowed from https://github.com/uber/tchannel-go/blob/dev/localip.go

// ScoreAddr scores how likely the given addr is to be a remote address and returns the
// IP to use when listening. Any address which receives a negative score should not be used.
// Scores are calculated as:
// -1 for any unknown IP addreseses.
// +300 for IPv4 addresses
// +100 for non-local addresses, extra +100 for "up" interaces.
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
	if isLocalMacAddr(iface.HardwareAddr) {
		score -= 50
	}
	return score, ip
}

func listenIP(interfaces []net.Interface) (net.IP, error) {
	bestScore := -1
	var bestIP net.IP
	// Select the highest scoring IP as the best IP.
	for _, iface := range interfaces {
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
			}
		}
	}

	if bestScore == -1 {
		return nil, errors.New("no addresses to listen on")
	}

	return bestIP, nil
}

// ListenIP returns the IP to bind to in Listen. It tries to find an IP that can be used
// by other machines to reach this machine.
func ListenIP() (net.IP, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	return listenIP(interfaces)
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
