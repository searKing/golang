// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package net_test

import (
	"fmt"
	"net"
	"reflect"
	"testing"

	net_ "github.com/searKing/golang/go/net"
)

var parseIPTests = []struct {
	ipIntStr string
	base     int
	ip       net.IP
	ipFmtStr string
}{
	{ipIntStr: "2130706690", base: 10, ipFmtStr: "127.0.1.2", ip: net.IPv4(127, 0, 1, 2)},
	{ipIntStr: "2130706433", base: 10, ipFmtStr: "127.0.0.1", ip: net.IPv4(127, 0, 0, 1)},
	{ipIntStr: "2130772483", base: 10, ipFmtStr: "127.001.002.003", ip: net.IPv4(127, 1, 2, 3)},
	{ipIntStr: "2130772483", base: 10, ipFmtStr: "::ffff:127.1.2.3", ip: net.IPv4(127, 1, 2, 3)},
	{ipIntStr: "2130772483", base: 10, ipFmtStr: "::ffff:127.001.002.003", ip: net.IPv4(127, 1, 2, 3)},
	{ipIntStr: "2130772483", base: 10, ipFmtStr: "::ffff:7f01:0203", ip: net.IPv4(127, 1, 2, 3)},
	{ipIntStr: "2130772483", base: 10, ipFmtStr: "0:0:0:0:0000:ffff:127.1.2.3", ip: net.IPv4(127, 1, 2, 3)},
	{ipIntStr: "2130772483", base: 10, ipFmtStr: "0:0:0:0:000000:ffff:127.1.2.3", ip: net.IPv4(127, 1, 2, 3)},
	{ipIntStr: "2130772483", base: 10, ipFmtStr: "0:0:0:0::ffff:127.1.2.3", ip: net.IPv4(127, 1, 2, 3)},

	{ipIntStr: "20014860000020010000000000000068", base: 16, ipFmtStr: "2001:4860:0:2001::68", ip: net.IP{0x20, 0x01, 0x48, 0x60, 0, 0, 0x20, 0x01, 0, 0, 0, 0, 0, 0, 0x00, 0x68}},
	{ipIntStr: "20014860000020010000000000000068", base: 16, ipFmtStr: "2001:4860:0000:2001:0000:0000:0000:0068", ip: net.IP{0x20, 0x01, 0x48, 0x60, 0, 0, 0x20, 0x01, 0, 0, 0, 0, 0, 0, 0x00, 0x68}},
}

func TestParseIPV6(t *testing.T) {
	for i, tt := range parseIPTests {
		if got := net_.ParseIPV6(tt.ipIntStr, tt.base); !reflect.DeepEqual(got, tt.ip) {
			t.Errorf("#%d, ParseIP(%s) = %v, want %v, %s", i, tt.ipIntStr, got, tt.ip, tt.ipFmtStr)
		}
	}
}

func TestIPtoInt64(t *testing.T) {
	for i, tt := range parseIPTests {
		got := net_.IPtoBigInt(tt.ip)
		var format string
		switch tt.base {
		case 2:
			format = "%b"
		case 8:
			format = "%o"
		case 16:
			format = "%x"
		default:
			format = "%d"
		}
		gotString := fmt.Sprintf(format, got)

		if gotString != tt.ipIntStr {
			t.Errorf("#%d, IPtoBigInt(%s) = %d, want %s", i, tt.ip.String(), got, tt.ipIntStr)
		}
	}
}

func TestIPFormat_Format(t *testing.T) {
	for i, tt := range parseIPTests {
		var format string
		switch tt.base {
		case 2:
			format = "%b"
		case 8:
			format = "%o"
		case 16:
			format = "%x"
		default:
			format = "%d"
		}
		gotString := fmt.Sprintf(format, net_.IPFormat(tt.ip))

		if gotString != tt.ipIntStr {
			t.Errorf("#%d, fmt.Sprintf(%q) = %s, want %s", i, format, gotString, tt.ipIntStr)
		}
	}
}
