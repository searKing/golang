// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package net

import (
	"encoding/binary"
	"fmt"
	"math/big"
	"net"
)

// INetAToN returns the numeric value of an IP address
var INetAToN = IP4toUInt32

// INetNToA returns the IP address from a numeric value
var INetNToA = ParseIPV4

// INet6AToN returns the numeric value of an IPv6 address
var INet6AToN = IPtoBigInt

// INet6NToA returns the IPv6 address from a numeric value
var INet6NToA = ParseIPV6

type IPFormat net.IP

func (ip IPFormat) Format(s fmt.State, verb rune) {
	switch verb {
	case 'b', 'o', 'O', 'd', 'x', 'X':
		IPtoBigInt(net.IP(ip)).Format(s, verb)
		return
	default:
		_, _ = fmt.Fprintf(s, "%"+string(verb), net.IP(ip))
	}
}

// ParseIPV4 parses i as an IP address, returning the result.
// If i is not a valid integer representation of an IP address,
// ParseIP returns nil.
func ParseIPV4(i uint32) net.IP {
	ip := make(net.IP, net.IPv4len)
	binary.BigEndian.PutUint32(ip, i)
	return ip.To4()
}

func ParseIPV6(s string, base int) net.IP {
	ipi, ok := big.NewInt(0).SetString(s, base)
	if !ok {
		return nil
	}
	ipb := ipi.Bytes()
	if len(ipb) == net.IPv4len {
		return net.ParseIP(net.IP(ipb).To4().String())
	}

	ip := make(net.IP, net.IPv6len)
	for i := 0; i < len(ip); i++ {
		j := len(ip) - 1 - i
		k := len(ipb) - 1 - i
		if k < 0 {
			break
		}
		ip[j] = ipb[k]
	}
	return ip.To16()
}

func ParseIP(s string) net.IP {
	return ParseIPV6(s, 10)
}

func IPtoBigInt(ip net.IP) *big.Int {
	ipInt := big.NewInt(0)
	// If IPv4, use dotted notation.
	if p4 := ip.To4(); len(p4) == net.IPv4len {
		ip = ip.To4()
	} else {
		ip = ip.To16()
	}
	ipInt.SetBytes(ip)
	return ipInt
}

func IP4toUInt32(ip net.IP) uint32 {
	ipInt := big.NewInt(0)
	// If IPv4, use dotted notation.
	if p4 := ip.To4(); len(p4) == net.IPv4len {
		ip = ip.To4()
	} else {
		return 0
	}
	ipInt.SetBytes(ip)
	return uint32(ipInt.Uint64())
}
