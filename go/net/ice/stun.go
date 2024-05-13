// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// https://tools.ietf.org/html/rfc7064
package ice

import (
	"fmt"
	"net/url"
)

// https://tools.ietf.org/html/rfc7064 3.  Definition of the "stun" or "stuns" URI
// stunURI       = scheme ":" host [ ":" port ]
// scheme        = "stun" / "stuns"
func parseStunProto(scheme Scheme, rawQuery string) (Transport, error) {
	// stunURI       = scheme ":" host [ ":" port ]
	// scheme        = "stun" / "stuns"
	qArgs, err := url.ParseQuery(rawQuery)
	if err != nil {
		return "", err
	}
	if len(qArgs) > 0 {
		return "", fmt.Errorf("malformed query %v ", qArgs)
	}
	switch scheme {
	case SchemeSTUN:
		return TransportUDP, nil
	case SchemeSTUNS:
		return TransportTCP, nil
	default:
		return "", fmt.Errorf("malformed scheme %s ", scheme.String())
	}
}
