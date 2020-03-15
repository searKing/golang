// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package url

import (
	"fmt"
	"reflect"
	"testing"
)

type protoTests struct {
	in  string
	out *Proto
}

var prototests = []protoTests{
	{
		"RTSP/1.0",
		&Proto{
			"RTSP",
			1,
			0,
		},
	},
}

func TestParseProto(t *testing.T) {
	for _, tt := range prototests {
		u, err := ParseProto(tt.in)
		if err != nil {
			t.Errorf("Parse(%q) returned error %v", tt.in, err)
			continue
		}
		if !reflect.DeepEqual(u, tt.out) {
			t.Errorf("Parse(%q):\n\tgot  %v\n\twant %v\n", tt.in, pfmt(u), pfmt(tt.out))
		}
	}
}

var stringProtoTests = []struct {
	proto Proto
	want  string
}{
	// No leading slash on path should prepend slash on String() call
	{
		proto: Proto{
			Type:  "RTSP",
			Major: 1,
			Minor: 0,
		},
		want: "RTSP/1.0",
	},
}

func TestProtoString(t *testing.T) {

	for _, tt := range prototests {
		u, err := ParseProto(tt.in)
		if err != nil {
			t.Errorf("Parse(%q) returned error %s", tt.in, err)
			continue
		}
		expected := tt.in
		s := u.String()
		if s != expected {
			t.Errorf("Parse(%q).String() == %q (expected %q)", tt.in, s, expected)
		}
	}

	for _, tt := range stringProtoTests {
		if got := tt.proto.String(); got != tt.want {
			t.Errorf("%+v.String() = %q; want %q", tt.proto, got, tt.want)
		}
	}
}

// more useful string for debugging than fmt's struct printer
func pfmt(p *Proto) string {
	return fmt.Sprintf("type=%q, major=%q, minor=%#v",
		p.Type, p.Major, p.Minor)
}
