// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package addr_test

import (
	"net"
	"testing"

	"github.com/searKing/golang/go/net/addr"
)

func TestScoreAddr(t *testing.T) {
	ipv4 := net.ParseIP("10.0.1.2")
	ipv6 := net.ParseIP("2001:db8:a0b:12f0::1")

	tests := []struct {
		msg       string
		iface     net.Interface
		addr      net.Addr
		wantScore int
		wantIP    net.IP
	}{
		{
			msg:       "non-local up ipv4 IPNet address",
			iface:     net.Interface{Flags: net.FlagUp},
			addr:      &net.IPNet{IP: ipv4},
			wantScore: 475,
			wantIP:    ipv4,
		},
		{
			msg:       "non-local up ipv4 IPAddr address",
			iface:     net.Interface{Flags: net.FlagUp},
			addr:      &net.IPAddr{IP: ipv4},
			wantScore: 475,
			wantIP:    ipv4,
		},
		{
			msg: "non-local up ipv4 IPAddr address, docker interface",
			iface: net.Interface{
				Flags:        net.FlagUp,
				HardwareAddr: mustParseMAC("02:42:ac:11:56:af"),
			},
			addr:      &net.IPNet{IP: ipv4},
			wantScore: 425,
			wantIP:    ipv4,
		},
		{
			msg: "non-local up ipv4 address, local MAC address",
			iface: net.Interface{
				Flags:        net.FlagUp,
				HardwareAddr: mustParseMAC("02:42:9c:52:fc:86"),
			},
			addr:      &net.IPNet{IP: ipv4},
			wantScore: 425,
			wantIP:    ipv4,
		},
		{
			msg:       "non-local down ipv4 address",
			iface:     net.Interface{},
			addr:      &net.IPNet{IP: ipv4},
			wantScore: 375,
			wantIP:    ipv4,
		},
		{
			msg:       "non-local down ipv6 address",
			iface:     net.Interface{},
			addr:      &net.IPAddr{IP: ipv6},
			wantScore: 75,
			wantIP:    ipv6,
		},
		{
			msg:       "unknown address type",
			iface:     net.Interface{},
			addr:      &net.UnixAddr{Name: "/tmp/socket"},
			wantScore: -1,
		},
	}

	for i, tt := range tests {
		gotScore, gotIP := addr.ScoreAddr(tt.iface, tt.addr)
		if tt.wantScore != gotScore {
			t.Errorf("#%d, %s: expected %d got %d", i, tt.msg, tt.wantScore, gotScore)
		}
		if tt.wantIP.String() != gotIP.String() {
			t.Errorf("#%d, %s: expected %q got %q", i, tt.msg, tt.wantIP, gotIP)
		}
	}
}

func mustParseMAC(s string) net.HardwareAddr {
	addr, err := net.ParseMAC(s)
	if err != nil {
		panic(err)
	}
	return addr
}
