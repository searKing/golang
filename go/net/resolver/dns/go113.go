// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// +build go1.13

package dns

import "net"

func init() {
	filterError = func(err error) error {
		if dnsErr, ok := err.(*net.DNSError); ok && dnsErr.IsNotFound {
			// The name does not exist; not an error.
			return nil
		}
		return err
	}
}
