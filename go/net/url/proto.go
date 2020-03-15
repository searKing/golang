// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package url

import (
	"errors"
	"fmt"
	"strings"
)

type Proto struct {
	Type  string // "RTSP"
	Major int    // 1
	Minor int    // 0
}

func ParseProto(rawproto string) (*Proto, error) {
	if rawproto == "" {
		return nil, errors.New("empty url_")
	}
	proto := new(Proto)
	var protoType string
	var protoMajor, protoMinor int

	protoType, protoVersion := split(rawproto, "/", true)
	if protoType == "" {
		return nil, errors.New("invalid PROTO for request : empty protoType")
	}
	n, err := fmt.Sscanf(protoVersion, "%d.%d", &protoMajor, &protoMinor)
	if n == 2 && err == nil {
		proto.Type = protoType
		proto.Major = protoMajor
		proto.Minor = protoMinor
		return proto, nil
	}
	return nil, errors.New("invalid PROTO for request")
}

// Maybe s is of the form t c u.
// If so, return t, c u (or t, u if cutc == true).
// If not, return s, "".
func split(s string, c string, cutc bool) (string, string) {
	i := strings.Index(s, c)
	if i < 0 {
		return s, ""
	}
	if cutc {
		return s[:i], s[i+len(c):]
	}
	return s[:i], s[i:]
}
func (p *Proto) String() string {
	return fmt.Sprintf("%s/%d.%d", p.Type, p.Major, p.Minor)
}
func (p *Proto) ProtoAtLeast(major, minor int) bool {
	return p.Major > major ||
		p.Major == major && p.Minor >= minor
}
