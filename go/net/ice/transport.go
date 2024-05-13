// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// https://tools.ietf.org/html/rfc7064 3.2.  URI Scheme Semantics
// https://tools.ietf.org/html/rfc7065 3.2.  URI Scheme Semantics
package ice

import (
	"errors"
)

// Transport indicates the transport protocol type that is used in the ice.URL
// structure.
type Transport string

const (
	// TransportUDP indicates the URL uses:
	// a UDP transport for turn|stun.
	// a DTLS-over-UDP transport for turns|stuns.
	TransportUDP Transport = "udp"

	// TransportTCP indicates the URL uses:
	// a TCP transport for turn|stun.
	// a TLS-over-TCP transport for turns|stuns.
	TransportTCP = "tcp"
)

func ParseProto(s string) (Transport, error) {
	if s == "" {
		return "", errors.New("empty proto")
	}
	return Transport(s), nil
}

// https://tools.ietf.org/html/rfc7065
// transport-ext = 1*unreserved
func (t Transport) String() string {
	if t == "" {
		return errors.New("empty proto").Error()
	}
	return string(t)
}
