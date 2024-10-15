// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package proto_test

import (
	testpb "github.com/searKing/golang/third_party/github.com/spf13/viper/proto/internal/testprotos/test"
	test3pb "github.com/searKing/golang/third_party/github.com/spf13/viper/proto/internal/testprotos/test3"
	testeditionspb "github.com/searKing/golang/third_party/github.com/spf13/viper/proto/internal/testprotos/testeditions"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func templateMessages(messages ...proto.Message) []protoreflect.MessageType {
	if len(messages) == 0 {
		messages = []proto.Message{
			(*testpb.TestAllTypes)(nil),
			(*test3pb.TestAllTypes)(nil),
			(*testpb.TestAllExtensions)(nil),
			(*testeditionspb.TestAllTypes)(nil),
		}
	}
	var out []protoreflect.MessageType
	for _, m := range messages {
		out = append(out, m.ProtoReflect().Type())
	}
	return out

}
