// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package metadata

import (
	"strings"

	"google.golang.org/grpc/metadata"
)

func New(k string, vals ...string) metadata.MD {
	md := metadata.MD{}
	key := strings.ToLower(k)
	for _, val := range vals {
		md[key] = append(md[key], val)
	}
	return md
}
