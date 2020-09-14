// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package addr

import (
	"fmt"
	"net"
	"strings"
	"time"
)

// http://:80
func BaseUrl(rawurl string) (addr string, err error) {
	return rawurl, nil
}
func LocalIPAddrByUDPMulticast() (string, error) {
	return "", nil
}

func LocalIPAddrByTcp(addr string) (string, error) {
	conn, err := net.DialTimeout("tcp", addr, time.Second)
	if err != nil {
		fmt.Println("Failed to get our IP address", err.Error())
	}
	defer conn.Close()

	return strings.Split(conn.LocalAddr().String(), ":")[0], nil
}

func parseAddr(url string) (addr string) {
	paths := strings.Split(url, "//")
	if len(paths) > 1 {
		tmp := strings.Split(paths[1], "/")
		addr = tmp[0]
		if !strings.Contains(addr, ":") {
			addr += ":80"
		}
	}
	return
}
