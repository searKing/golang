// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// https://tools.ietf.org/html/rfc7065
package ice

import (
	"fmt"
	"net/url"
)

// https://tools.ietf.org/html/rfc7065 3.  Definitions of the "turn" and "turns" URI.
// turnURI       = scheme ":" host [ ":" port ]
//
//	[ "?transport=" transport ]
//
// scheme        = "turn" / "turns"
// transport     = "udp" / "tcp" / transport-ext
// transport-ext = 1*unreserved
func parseTurnProto(scheme Scheme, rawQuery string) (Transport, error) {
	form, err := url.ParseQuery(rawQuery)
	if err != nil {
		return "", err
	}
	if len(form) > 1 {
		return "", fmt.Errorf("malformed query:%v", form)
	}

	if proto := form.Get("transport"); proto != "" {
		return Transport(proto), nil
	}
	form.Del("transport")
	if len(form) > 0 {
		return "", fmt.Errorf("malformed query:%v", form)
	}

	switch scheme {
	case SchemeTURN:
		return TransportUDP, nil
	case SchemeTURNS:
		return TransportTCP, nil
	default:
		return "", fmt.Errorf("malformed scheme %s ", scheme.String())
	}

}
