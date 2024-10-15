// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package proto_test

import (
	"fmt"
	"maps"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/testing/protopack"

	maps_ "github.com/searKing/golang/go/exp/maps"
	slices_ "github.com/searKing/golang/go/exp/slices"
	proto_ "github.com/searKing/golang/third_party/github.com/spf13/viper/proto"
	"github.com/searKing/golang/third_party/github.com/spf13/viper/proto/internal/protobuild"
	legacypb "github.com/searKing/golang/third_party/github.com/spf13/viper/proto/internal/testprotos/legacy"
	testpb "github.com/searKing/golang/third_party/github.com/spf13/viper/proto/internal/testprotos/test"
	test3pb "github.com/searKing/golang/third_party/github.com/spf13/viper/proto/internal/testprotos/test3"
)

type testMerge struct {
	desc  string
	dst   map[any]any
	src   protobuild.Message
	want  map[any]any // if dst and want are nil, want = src
	types []proto.Message
	opts  []proto_.MergeOption
}

var testMerges = []testMerge{{
	desc: "clone a large message",
	src: protobuild.Message{
		"optional_int32":          1001,
		"optional_int64":          1002,
		"optional_uint32":         1003,
		"optional_uint64":         1004,
		"optional_sint32":         1005,
		"optional_sint64":         1006,
		"optional_fixed32":        1007,
		"optional_fixed64":        1008,
		"optional_sfixed32":       1009,
		"optional_sfixed64":       1010,
		"optional_float":          1011.5,
		"optional_double":         1012.5,
		"optional_bool":           true,
		"optional_string":         "string",
		"optional_bytes":          []byte("bytes"),
		"optional_nested_enum":    1,
		"optional_nested_message": protobuild.Message{"a": 100},
		"repeated_int32":          []int32{1001, 2001},
		"repeated_int64":          []int64{1002, 2002},
		"repeated_uint32":         []uint32{1003, 2003},
		"repeated_uint64":         []uint64{1004, 2004},
		"repeated_sint32":         []int32{1005, 2005},
		"repeated_sint64":         []int64{1006, 2006},
		"repeated_fixed32":        []uint32{1007, 2007},
		"repeated_fixed64":        []uint64{1008, 2008},
		"repeated_sfixed32":       []int32{1009, 2009},
		"repeated_sfixed64":       []int64{1010, 2010},
		"repeated_float":          []float32{1011.5, 2011.5},
		"repeated_double":         []float64{1012.5, 2012.5},
		"repeated_bool":           []bool{true, false},
		"repeated_string":         []string{"foo", "bar"},
		"repeated_bytes":          []string{"FOO", "BAR"},
		"repeated_nested_enum":    []string{"FOO", "BAR"},
		"repeated_nested_message": []protobuild.Message{
			{"a": 200},
			{"a": 300},
		},
	},
	want: map[any]any{
		"optional_int32":          int32(1001),
		"optional_int64":          int64(1002),
		"optional_uint32":         uint32(1003),
		"optional_uint64":         uint64(1004),
		"optional_sint32":         int32(1005),
		"optional_sint64":         int64(1006),
		"optional_fixed32":        uint32(1007),
		"optional_fixed64":        uint64(1008),
		"optional_sfixed32":       int32(1009),
		"optional_sfixed64":       int64(1010),
		"optional_float":          float32(1011.5),
		"optional_double":         float64(1012.5),
		"optional_bool":           true,
		"optional_string":         "string",
		"optional_bytes":          []byte("bytes"),
		"optional_nested_enum":    "BAR",
		"optional_nested_message": map[any]any{"a": int32(100)},
		"repeated_int32":          []any{int32(1001), int32(2001)},
		"repeated_int64":          []any{int64(1002), int64(2002)},
		"repeated_uint32":         []any{uint32(1003), uint32(2003)},
		"repeated_uint64":         []any{uint64(1004), uint64(2004)},
		"repeated_sint32":         []any{int32(1005), int32(2005)},
		"repeated_sint64":         []any{int64(1006), int64(2006)},
		"repeated_fixed32":        []any{uint32(1007), uint32(2007)},
		"repeated_fixed64":        []any{uint64(1008), uint64(2008)},
		"repeated_sfixed32":       []any{int32(1009), int32(2009)},
		"repeated_sfixed64":       []any{int64(1010), int64(2010)},
		"repeated_float":          []any{float32(1011.5), float32(2011.5)},
		"repeated_double":         []any{float64(1012.5), float64(2012.5)},
		"repeated_bool":           []any{true, false},
		"repeated_string":         []any{"foo", "bar"},
		"repeated_bytes":          []any{[]byte("FOO"), []byte("BAR")},
		"repeated_nested_enum":    []any{"FOO", "BAR"},
		"repeated_nested_message": []any{
			map[any]any{"a": int32(200)},
			map[any]any{"a": int32(300)},
		},
	},
	types: nil,
}, {
	desc: "clone maps",
	src: protobuild.Message{
		"map_int32_int32":       map[int32]int32{1056: 1156, 2056: 2156},
		"map_int64_int64":       map[int64]int64{1057: 1157, 2057: 2157},
		"map_uint32_uint32":     map[uint32]uint32{1058: 1158, 2058: 2158},
		"map_uint64_uint64":     map[uint64]uint64{1059: 1159, 2059: 2159},
		"map_sint32_sint32":     map[int32]int32{1060: 1160, 2060: 2160},
		"map_sint64_sint64":     map[int64]int64{1061: 1161, 2061: 2161},
		"map_fixed32_fixed32":   map[uint32]uint32{1062: 1162, 2062: 2162},
		"map_fixed64_fixed64":   map[uint64]uint64{1063: 1163, 2063: 2163},
		"map_sfixed32_sfixed32": map[int32]int32{1064: 1164, 2064: 2164},
		"map_sfixed64_sfixed64": map[int64]int64{1065: 1165, 2065: 2165},
		"map_int32_float":       map[int32]float32{1066: 1166.5, 2066: 2166.5},
		"map_int32_double":      map[int32]float64{1067: 1167.5, 2067: 2167.5},
		"map_bool_bool":         map[bool]bool{true: false, false: true},
		"map_string_string":     map[string]string{"69.1.key": "69.1.val", "69.2.key": "69.2.val"},
		"map_string_bytes":      map[string][]byte{"70.1.key": []byte("70.1.val"), "70.2.key": []byte("70.2.val")},
		"map_string_nested_message": map[string]protobuild.Message{
			"71.1.key": {"a": 1171},
			"71.2.key": {"a": 2171},
		},
		"map_string_nested_enum": map[string]string{"73.1.key": "FOO", "73.2.key": "BAR"},
	}, want: map[any]any{
		"map_int32_int32":       map[any]any{int32(1056): int32(1156), int32(2056): int32(2156)},
		"map_int64_int64":       map[any]any{int64(1057): int64(1157), int64(2057): int64(2157)},
		"map_uint32_uint32":     map[any]any{uint32(1058): uint32(1158), uint32(2058): uint32(2158)},
		"map_uint64_uint64":     map[any]any{uint64(1059): uint64(1159), uint64(2059): uint64(2159)},
		"map_sint32_sint32":     map[any]any{int32(1060): int32(1160), int32(2060): int32(2160)},
		"map_sint64_sint64":     map[any]any{int64(1061): int64(1161), int64(2061): int64(2161)},
		"map_fixed32_fixed32":   map[any]any{uint32(1062): uint32(1162), uint32(2062): uint32(2162)},
		"map_fixed64_fixed64":   map[any]any{uint64(1063): uint64(1163), uint64(2063): uint64(2163)},
		"map_sfixed32_sfixed32": map[any]any{int32(1064): int32(1164), int32(2064): int32(2164)},
		"map_sfixed64_sfixed64": map[any]any{int64(1065): int64(1165), int64(2065): int64(2165)},
		"map_int32_float":       map[any]any{int32(1066): float32(1166.5), int32(2066): float32(2166.5)},
		"map_int32_double":      map[any]any{int32(1067): float64(1167.5), int32(2067): float64(2167.5)},
		"map_bool_bool":         map[any]any{true: false, false: true},
		"map_string_string":     map[any]any{"69.1.key": "69.1.val", "69.2.key": "69.2.val"},
		"map_string_bytes":      map[any]any{"70.1.key": []byte("70.1.val"), "70.2.key": []byte("70.2.val")},
		"map_string_nested_message": map[any]any{
			"71.1.key": map[any]any{"a": int32(1171)},
			"71.2.key": map[any]any{"a": int32(2171)},
		},
		"map_string_nested_enum": map[any]any{"73.1.key": "FOO", "73.2.key": "BAR"},
	},
	types: []proto.Message{&testpb.TestAllTypes{}, &test3pb.TestAllTypes{}},
}, {
	desc: "clone oneof uint32",
	src: protobuild.Message{
		"oneof_uint32": 1111,
	},
	want:  map[any]any{"oneof_uint32": uint32(1111)},
	types: []proto.Message{&testpb.TestAllTypes{}, &test3pb.TestAllTypes{}},
}, {
	desc: "clone oneof string",
	src: protobuild.Message{
		"oneof_string": "string",
	},
	want:  map[any]any{"oneof_string": "string"},
	types: []proto.Message{&testpb.TestAllTypes{}, &test3pb.TestAllTypes{}},
}, {
	desc: "clone oneof bytes",
	src: protobuild.Message{
		"oneof_bytes": "bytes",
	},
	want:  map[any]any{"oneof_bytes": []byte("bytes")},
	types: []proto.Message{&testpb.TestAllTypes{}, &test3pb.TestAllTypes{}},
}, {
	desc: "clone oneof bool",
	src: protobuild.Message{
		"oneof_bool": true,
	},
	want:  map[any]any{"oneof_bool": true},
	types: []proto.Message{&testpb.TestAllTypes{}, &test3pb.TestAllTypes{}},
}, {
	desc: "clone oneof uint64",
	src: protobuild.Message{
		"oneof_uint64": 100,
	},
	want:  map[any]any{"oneof_uint64": uint64(100)},
	types: []proto.Message{&testpb.TestAllTypes{}, &test3pb.TestAllTypes{}},
}, {
	desc: "clone oneof float",
	src: protobuild.Message{
		"oneof_float": 100,
	},
	want:  map[any]any{"oneof_float": float32(100)},
	types: []proto.Message{&testpb.TestAllTypes{}, &test3pb.TestAllTypes{}},
}, {
	desc: "clone oneof double",
	src: protobuild.Message{
		"oneof_double": 1111,
	},
	want:  map[any]any{"oneof_double": float64(1111)},
	types: []proto.Message{&testpb.TestAllTypes{}, &test3pb.TestAllTypes{}},
}, {
	desc: "clone oneof enum",
	src: protobuild.Message{
		"oneof_enum": 1,
	},
	want:  map[any]any{"oneof_enum": string("BAR")},
	types: []proto.Message{&testpb.TestAllTypes{}, &test3pb.TestAllTypes{}},
}, {
	desc: "clone oneof message",
	src: protobuild.Message{
		"oneof_nested_message": protobuild.Message{
			"a": 1,
		},
	},
	want:  map[any]any{"oneof_nested_message": map[any]any{"a": int32(1)}},
	types: []proto.Message{&testpb.TestAllTypes{}, &test3pb.TestAllTypes{}},
}, {
	desc: "clone oneof group",
	src: protobuild.Message{
		"oneofgroup": protobuild.Message{
			"a": 1,
		},
	},
	want:  map[any]any{"OneofGroup": map[any]any{"a": int32(1)}},
	types: []proto.Message{&testpb.TestAllTypes{}},
}, {
	desc: "merge bytes",
	dst: map[any]any{
		"optional_bytes":   []byte{1, 2, 3},
		"repeated_bytes":   [][]byte{{1, 2}, {3, 4}},
		"map_string_bytes": map[string][]byte{"alpha": {1, 2, 3}},
	},
	src: protobuild.Message{
		"optional_bytes":   []byte{4, 5, 6},
		"repeated_bytes":   [][]byte{{5, 6}, {7, 8}},
		"map_string_bytes": map[string][]byte{"alpha": {4, 5, 6}, "bravo": {1, 2, 3}},
	},
	want: map[any]any{
		"optional_bytes":   []byte{4, 5, 6},
		"repeated_bytes":   []any{[]byte{1, 2}, []byte{3, 4}, []byte{5, 6}, []byte{7, 8}},
		"map_string_bytes": map[any]any{"alpha": []byte{4, 5, 6}, "bravo": []byte{1, 2, 3}},
	},
	types: []proto.Message{&testpb.TestAllTypes{}, &test3pb.TestAllTypes{}},
}, {
	desc: "merge singular fields",
	dst: map[any]any{
		"optional_int32":       1,
		"optional_int64":       1,
		"optional_uint32":      1,
		"optional_uint64":      1,
		"optional_sint32":      1,
		"optional_sint64":      1,
		"optional_fixed32":     1,
		"optional_fixed64":     1,
		"optional_sfixed32":    1,
		"optional_sfixed64":    1,
		"optional_float":       1,
		"optional_double":      1,
		"optional_bool":        false,
		"optional_string":      "1",
		"optional_bytes":       "1",
		"optional_nested_enum": 1,
		"optional_nested_message": map[any]any{
			"a": int64(1),
			"corecursive": map[any]any{
				"optional_int64": int64(1),
			},
		},
	},
	src: protobuild.Message{
		"optional_int32":       2,
		"optional_int64":       2,
		"optional_uint32":      2,
		"optional_uint64":      2,
		"optional_sint32":      2,
		"optional_sint64":      2,
		"optional_fixed32":     2,
		"optional_fixed64":     2,
		"optional_sfixed32":    2,
		"optional_sfixed64":    2,
		"optional_float":       2,
		"optional_double":      2,
		"optional_bool":        true,
		"optional_string":      "2",
		"optional_bytes":       "2",
		"optional_nested_enum": 2,
		"optional_nested_message": protobuild.Message{
			"a": 2,
			"corecursive": protobuild.Message{
				"optional_int64": 2,
			},
		},
	},
	want: map[any]any{
		"optional_int32":       int32(2),
		"optional_int64":       int64(2),
		"optional_uint32":      uint32(2),
		"optional_uint64":      uint64(2),
		"optional_sint32":      int32(2),
		"optional_sint64":      int64(2),
		"optional_fixed32":     uint32(2),
		"optional_fixed64":     uint64(2),
		"optional_sfixed32":    int32(2),
		"optional_sfixed64":    int64(2),
		"optional_float":       float32(2),
		"optional_double":      float64(2),
		"optional_bool":        true,
		"optional_string":      "2",
		"optional_bytes":       []byte("2"),
		"optional_nested_enum": "BAZ",
		"optional_nested_message": map[any]any{
			"a": int32(2),
			"corecursive": map[any]any{
				"optional_int64": int64(2),
			},
		},
	},
}, {
	desc: "no merge of empty singular fields",
	dst: map[any]any{
		"optional_int32":       int32(1),
		"optional_int64":       int64(1),
		"optional_uint32":      uint32(1),
		"optional_uint64":      uint64(1),
		"optional_sint32":      int32(1),
		"optional_sint64":      int64(1),
		"optional_fixed32":     uint32(1),
		"optional_fixed64":     uint64(1),
		"optional_sfixed32":    int32(1),
		"optional_sfixed64":    int64(1),
		"optional_float":       float32(1),
		"optional_double":      float64(1),
		"optional_bool":        false,
		"optional_string":      "1",
		"optional_bytes":       []byte("1"),
		"optional_nested_enum": 1,
		"optional_nested_message": map[any]any{
			"a": int32(1),
			"corecursive": map[any]any{
				"optional_int64": int64(1),
			},
		},
	},
	src: protobuild.Message{
		"optional_nested_message": protobuild.Message{
			"a": 1,
			"corecursive": protobuild.Message{
				"optional_int32": 2,
			},
		},
	},
	want: map[any]any{
		"optional_int32":       int32(1),
		"optional_int64":       int64(1),
		"optional_uint32":      uint32(1),
		"optional_uint64":      uint64(1),
		"optional_sint32":      int32(1),
		"optional_sint64":      int64(1),
		"optional_fixed32":     uint32(1),
		"optional_fixed64":     uint64(1),
		"optional_sfixed32":    int32(1),
		"optional_sfixed64":    int64(1),
		"optional_float":       float32(1),
		"optional_double":      float64(1),
		"optional_bool":        false,
		"optional_string":      "1",
		"optional_bytes":       []byte("1"),
		"optional_nested_enum": 1,
		"optional_nested_message": map[any]any{
			"a": int32(1),
			"corecursive": map[any]any{
				"optional_int32": int32(2),
				"optional_int64": int64(1),
			},
		},
	},
}, {
	desc: "merge list fields",
	dst: map[any]any{
		"repeated_int32":       toSlicesOfAny([]int32{1, 2, 3}),
		"repeated_int64":       toSlicesOfAny([]int64{1, 2, 3}),
		"repeated_uint32":      toSlicesOfAny([]uint32{1, 2, 3}),
		"repeated_uint64":      toSlicesOfAny([]uint64{1, 2, 3}),
		"repeated_sint32":      toSlicesOfAny([]int32{1, 2, 3}),
		"repeated_sint64":      toSlicesOfAny([]int64{1, 2, 3}),
		"repeated_fixed32":     toSlicesOfAny([]uint32{1, 2, 3}),
		"repeated_fixed64":     toSlicesOfAny([]uint64{1, 2, 3}),
		"repeated_sfixed32":    toSlicesOfAny([]int32{1, 2, 3}),
		"repeated_sfixed64":    toSlicesOfAny([]int64{1, 2, 3}),
		"repeated_float":       toSlicesOfAny([]float32{1, 2, 3}),
		"repeated_double":      toSlicesOfAny([]float64{1, 2, 3}),
		"repeated_bool":        toSlicesOfAny([]bool{true}),
		"repeated_string":      toSlicesOfAny([]string{"a", "b", "c"}),
		"repeated_bytes":       toSlicesOfAny([][]byte{[]byte("a"), []byte("b"), []byte("c")}),
		"repeated_nested_enum": []int32{int32(1), int32(2), int32(3)},
		"repeated_nested_message": []any{
			map[any]any{"a": int32(100)},
			map[any]any{"a": int32(200)},
		},
	},
	src: protobuild.Message{
		"repeated_int32":       []int32{4, 5, 6},
		"repeated_int64":       []int64{4, 5, 6},
		"repeated_uint32":      []uint32{4, 5, 6},
		"repeated_uint64":      []uint64{4, 5, 6},
		"repeated_sint32":      []int32{4, 5, 6},
		"repeated_sint64":      []int64{4, 5, 6},
		"repeated_fixed32":     []uint32{4, 5, 6},
		"repeated_fixed64":     []uint64{4, 5, 6},
		"repeated_sfixed32":    []int32{4, 5, 6},
		"repeated_sfixed64":    []int64{4, 5, 6},
		"repeated_float":       []float32{4, 5, 6},
		"repeated_double":      []float64{4, 5, 6},
		"repeated_bool":        []bool{false},
		"repeated_string":      []string{"d", "e", "f"},
		"repeated_bytes":       []string{"d", "e", "f"},
		"repeated_nested_enum": []int{4, 5, 6},
		"repeated_nested_message": []protobuild.Message{
			{"a": 300},
			{"a": 400},
		},
	},
	want: map[any]any{
		"repeated_int32":       toSlicesOfAny([]int32{1, 2, 3, 4, 5, 6}),
		"repeated_int64":       toSlicesOfAny([]int64{1, 2, 3, 4, 5, 6}),
		"repeated_uint32":      toSlicesOfAny([]uint32{1, 2, 3, 4, 5, 6}),
		"repeated_uint64":      toSlicesOfAny([]uint64{1, 2, 3, 4, 5, 6}),
		"repeated_sint32":      toSlicesOfAny([]int32{1, 2, 3, 4, 5, 6}),
		"repeated_sint64":      toSlicesOfAny([]int64{1, 2, 3, 4, 5, 6}),
		"repeated_fixed32":     toSlicesOfAny([]uint32{1, 2, 3, 4, 5, 6}),
		"repeated_fixed64":     toSlicesOfAny([]uint64{1, 2, 3, 4, 5, 6}),
		"repeated_sfixed32":    toSlicesOfAny([]int32{1, 2, 3, 4, 5, 6}),
		"repeated_sfixed64":    toSlicesOfAny([]int64{1, 2, 3, 4, 5, 6}),
		"repeated_float":       toSlicesOfAny([]float32{1, 2, 3, 4, 5, 6}),
		"repeated_double":      toSlicesOfAny([]float64{1, 2, 3, 4, 5, 6}),
		"repeated_bool":        toSlicesOfAny([]bool{true, false}),
		"repeated_string":      toSlicesOfAny([]string{"a", "b", "c", "d", "e", "f"}),
		"repeated_bytes":       toSlicesOfAny([][]byte{[]byte("a"), []byte("b"), []byte("c"), []byte("d"), []byte("e"), []byte("f")}),
		"repeated_nested_enum": []any{int32(1), int32(2), int32(3), int32(4), int32(5), int32(6)},
		"repeated_nested_message": []any{
			map[any]any{"a": int32(100)},
			map[any]any{"a": int32(200)},
			map[any]any{"a": int32(300)},
			map[any]any{"a": int32(400)},
		},
	},
}, {
	desc: "merge map fields",
	dst: map[any]any{
		"map_int32_int32":       toMapsOfAny(map[int32]int32{1: 1, 3: 1}),
		"map_int64_int64":       toMapsOfAny(map[int64]int64{1: 1, 3: 1}),
		"map_uint32_uint32":     toMapsOfAny(map[uint32]uint32{1: 1, 3: 1}),
		"map_uint64_uint64":     toMapsOfAny(map[uint64]uint64{1: 1, 3: 1}),
		"map_sint32_sint32":     toMapsOfAny(map[int32]int32{1: 1, 3: 1}),
		"map_sint64_sint64":     toMapsOfAny(map[int64]int64{1: 1, 3: 1}),
		"map_fixed32_fixed32":   toMapsOfAny(map[uint32]uint32{1: 1, 3: 1}),
		"map_fixed64_fixed64":   toMapsOfAny(map[uint64]uint64{1: 1, 3: 1}),
		"map_sfixed32_sfixed32": toMapsOfAny(map[int32]int32{1: 1, 3: 1}),
		"map_sfixed64_sfixed64": toMapsOfAny(map[int64]int64{1: 1, 3: 1}),
		"map_int32_float":       toMapsOfAny(map[int32]float32{1: 1, 3: 1}),
		"map_int32_double":      toMapsOfAny(map[int32]float64{1: 1, 3: 1}),
		"map_bool_bool":         toMapsOfAny(map[bool]bool{true: true}),
		"map_string_string":     toMapsOfAny(map[string]string{"a": "1", "ab": "1"}),
		"map_string_bytes":      toMapsOfAny(map[string][]byte{"a": []byte("1"), "ab": []byte("1")}),
		"map_string_nested_message": map[any]any{
			"a": map[any]any{"a": int32(1)},
			"ab": map[any]any{
				"a": int32(1),
				"corecursive": map[any]any{
					"map_int32_int32": toMapsOfAny(map[int32]int32{1: 1, 3: 1}),
				},
			},
		},
		"map_string_nested_enum": map[string]int{"a": 1, "ab": 1},
	},
	src: protobuild.Message{
		"map_int32_int32":       map[int]int{2: 2, 3: 2},
		"map_int64_int64":       map[int]int{2: 2, 3: 2},
		"map_uint32_uint32":     map[int]int{2: 2, 3: 2},
		"map_uint64_uint64":     map[int]int{2: 2, 3: 2},
		"map_sint32_sint32":     map[int]int{2: 2, 3: 2},
		"map_sint64_sint64":     map[int]int{2: 2, 3: 2},
		"map_fixed32_fixed32":   map[int]int{2: 2, 3: 2},
		"map_fixed64_fixed64":   map[int]int{2: 2, 3: 2},
		"map_sfixed32_sfixed32": map[int]int{2: 2, 3: 2},
		"map_sfixed64_sfixed64": map[int]int{2: 2, 3: 2},
		"map_int32_float":       map[int]int{2: 2, 3: 2},
		"map_int32_double":      map[int]int{2: 2, 3: 2},
		"map_bool_bool":         map[bool]bool{false: false},
		"map_string_string":     map[string]string{"b": "2", "ab": "2"},
		"map_string_bytes":      map[string]string{"b": "2", "ab": "2"},
		"map_string_nested_message": map[string]protobuild.Message{
			"b": {"a": 2},
			"ab": {
				"a": 2,
				"corecursive": protobuild.Message{
					"map_int32_int32": map[int]int{2: 2, 3: 2},
				},
			},
		},
		"map_string_nested_enum": map[string]int{"b": 2, "ab": 2},
	},
	want: map[any]any{
		"map_int32_int32":       toMapsOfAny(map[int32]int32{1: 1, 2: 2, 3: 2}),
		"map_int64_int64":       toMapsOfAny(map[int64]int64{1: 1, 2: 2, 3: 2}),
		"map_uint32_uint32":     toMapsOfAny(map[uint32]uint32{1: 1, 2: 2, 3: 2}),
		"map_uint64_uint64":     toMapsOfAny(map[uint64]uint64{1: 1, 2: 2, 3: 2}),
		"map_sint32_sint32":     toMapsOfAny(map[int32]int32{1: 1, 2: 2, 3: 2}),
		"map_sint64_sint64":     toMapsOfAny(map[int64]int64{1: 1, 2: 2, 3: 2}),
		"map_fixed32_fixed32":   toMapsOfAny(map[uint32]uint32{1: 1, 2: 2, 3: 2}),
		"map_fixed64_fixed64":   toMapsOfAny(map[uint64]uint64{1: 1, 2: 2, 3: 2}),
		"map_sfixed32_sfixed32": toMapsOfAny(map[int32]int32{1: 1, 2: 2, 3: 2}),
		"map_sfixed64_sfixed64": toMapsOfAny(map[int64]int64{1: 1, 2: 2, 3: 2}),
		"map_int32_float":       toMapsOfAny(map[int32]float32{1: 1, 2: 2, 3: 2}),
		"map_int32_double":      toMapsOfAny(map[int32]float64{1: 1, 2: 2, 3: 2}),
		"map_bool_bool":         toMapsOfAny(map[bool]bool{true: true, false: false}),
		"map_string_string":     toMapsOfAny(map[string]string{"a": "1", "b": "2", "ab": "2"}),
		"map_string_bytes":      toMapsOfAny(map[string][]byte{"a": []byte("1"), "b": []byte("2"), "ab": []byte("2")}),
		"map_string_nested_message": map[any]any{
			"a": map[any]any{"a": int32(1)},
			"b": map[any]any{"a": int32(2)},
			"ab": map[any]any{
				"a": int32(2),
				"corecursive": map[any]any{
					// The map item "ab" was entirely replaced, so
					// this does not contain 1:1 from dst.
					"map_int32_int32": toMapsOfAny(map[int32]int32{2: 2, 3: 2}),
				},
			},
		},
		"map_string_nested_enum": map[any]any{"a": 1, "b": "BAZ", "ab": "BAZ"},
	},
	types: []proto.Message{&testpb.TestAllTypes{}, &test3pb.TestAllTypes{}},
}, {
	desc: "merge oneof message fields",
	dst: map[any]any{
		"oneof_nested_message": map[any]any{
			"a": 100,
		},
	},
	src: protobuild.Message{
		"oneof_nested_message": protobuild.Message{
			"corecursive": protobuild.Message{
				"optional_int64": 1000,
			},
		},
	},
	want: map[any]any{
		"oneof_nested_message": map[any]any{
			"a": 100,
			"corecursive": map[any]any{
				"optional_int64": int64(1000),
			},
		},
	},
	types: []proto.Message{&testpb.TestAllTypes{}, &test3pb.TestAllTypes{}},
}, {
	desc: "merge oneof scalar fields",
	dst: map[any]any{
		"oneof_uint32": uint32(100),
	},
	src: protobuild.Message{
		"oneof_float": 3.14152,
	},
	want: map[any]any{
		"oneof_float": float32(3.14152),
	},
	types: []proto.Message{&testpb.TestAllTypes{}, &test3pb.TestAllTypes{}},
}, {
	desc: "merge unknown fields",
	opts: []proto_.MergeOption{proto_.WithMergeEmitUnknown(true)},
	dst: map[any]any{
		protobuild.Unknown: protopack.Message{
			protopack.Tag{Number: 50000, Type: protopack.VarintType}, protopack.Svarint(-5),
		}.Marshal(),
	},
	src: protobuild.Message{
		protobuild.Unknown: protopack.Message{
			protopack.Tag{Number: 500000, Type: protopack.VarintType}, protopack.Svarint(-50),
		}.Marshal(),
	},
	want: map[any]any{
		protobuild.Unknown: protopack.Message{
			protopack.Tag{Number: 50000, Type: protopack.VarintType}, protopack.Svarint(-5),
			protopack.Tag{Number: 500000, Type: protopack.VarintType}, protopack.Svarint(-50),
		}.Marshal(),
	},
}, {
	desc: "clone legacy message",
	src: protobuild.Message{"f1": protobuild.Message{
		"optional_int32":        1,
		"optional_int64":        1,
		"optional_uint32":       1,
		"optional_uint64":       1,
		"optional_sint32":       1,
		"optional_sint64":       1,
		"optional_fixed32":      1,
		"optional_fixed64":      1,
		"optional_sfixed32":     1,
		"optional_sfixed64":     1,
		"optional_float":        1,
		"optional_double":       1,
		"optional_bool":         true,
		"optional_string":       "string",
		"optional_bytes":        "bytes",
		"optional_sibling_enum": 1,
		"optional_sibling_message": protobuild.Message{
			"f1": "value",
		},
		"repeated_int32":        []int32{1},
		"repeated_int64":        []int64{1},
		"repeated_uint32":       []uint32{1},
		"repeated_uint64":       []uint64{1},
		"repeated_sint32":       []int32{1},
		"repeated_sint64":       []int64{1},
		"repeated_fixed32":      []uint32{1},
		"repeated_fixed64":      []uint64{1},
		"repeated_sfixed32":     []int32{1},
		"repeated_sfixed64":     []int64{1},
		"repeated_float":        []float32{1},
		"repeated_double":       []float64{1},
		"repeated_bool":         []bool{true},
		"repeated_string":       []string{"string"},
		"repeated_bytes":        []string{"bytes"},
		"repeated_sibling_enum": []int{1},
		"repeated_sibling_message": []protobuild.Message{
			{"f1": "1"},
		},
		"map_bool_int32":    map[bool]int{true: 1},
		"map_bool_int64":    map[bool]int{true: 1},
		"map_bool_uint32":   map[bool]int{true: 1},
		"map_bool_uint64":   map[bool]int{true: 1},
		"map_bool_sint32":   map[bool]int{true: 1},
		"map_bool_sint64":   map[bool]int{true: 1},
		"map_bool_fixed32":  map[bool]int{true: 1},
		"map_bool_fixed64":  map[bool]int{true: 1},
		"map_bool_sfixed32": map[bool]int{true: 1},
		"map_bool_sfixed64": map[bool]int{true: 1},
		"map_bool_float":    map[bool]int{true: 1},
		"map_bool_double":   map[bool]int{true: 1},
		"map_bool_bool":     map[bool]bool{true: false},
		"map_bool_string":   map[bool]string{true: "1"},
		"map_bool_bytes":    map[bool]string{true: "1"},
		"map_bool_sibling_message": map[bool]protobuild.Message{
			true: {"f1": "1"},
		},
		"map_bool_sibling_enum": map[bool]int{true: 1},
		"oneof_sibling_message": protobuild.Message{
			"f1": "1",
		},
	}},
	want: map[any]any{
		"f1": map[any]any{
			"optional_int32":        int32(1),
			"optional_int64":        int64(1),
			"optional_uint32":       uint32(1),
			"optional_uint64":       uint64(1),
			"optional_sint32":       int32(1),
			"optional_sint64":       int64(1),
			"optional_fixed32":      uint32(1),
			"optional_fixed64":      uint64(1),
			"optional_sfixed32":     int32(1),
			"optional_sfixed64":     int64(1),
			"optional_float":        float32(1),
			"optional_double":       float64(1),
			"optional_bool":         true,
			"optional_string":       "string",
			"optional_bytes":        []byte("bytes"),
			"optional_sibling_enum": int32(1),
			"optional_sibling_message": map[any]any{
				"f1": "value",
			},
			"repeated_int32":        toSlicesOfAny([]int32{1}),
			"repeated_int64":        toSlicesOfAny([]int64{1}),
			"repeated_uint32":       toSlicesOfAny([]uint32{1}),
			"repeated_uint64":       toSlicesOfAny([]uint64{1}),
			"repeated_sint32":       toSlicesOfAny([]int32{1}),
			"repeated_sint64":       toSlicesOfAny([]int64{1}),
			"repeated_fixed32":      toSlicesOfAny([]uint32{1}),
			"repeated_fixed64":      toSlicesOfAny([]uint64{1}),
			"repeated_sfixed32":     toSlicesOfAny([]int32{1}),
			"repeated_sfixed64":     toSlicesOfAny([]int64{1}),
			"repeated_float":        toSlicesOfAny([]float32{1}),
			"repeated_double":       toSlicesOfAny([]float64{1}),
			"repeated_bool":         toSlicesOfAny([]bool{true}),
			"repeated_string":       toSlicesOfAny([]string{"string"}),
			"repeated_bytes":        toSlicesOfAny([][]byte{[]byte("bytes")}),
			"repeated_sibling_enum": toSlicesOfAny([]int32{1}),
			"repeated_sibling_message": []any{
				map[any]any{"f1": "1"},
			},
			"map_bool_int32":    toMapsOfAny(map[bool]int32{true: 1}),
			"map_bool_int64":    toMapsOfAny(map[bool]int64{true: 1}),
			"map_bool_uint32":   toMapsOfAny(map[bool]uint32{true: 1}),
			"map_bool_uint64":   toMapsOfAny(map[bool]uint64{true: 1}),
			"map_bool_sint32":   toMapsOfAny(map[bool]int32{true: 1}),
			"map_bool_sint64":   toMapsOfAny(map[bool]int64{true: 1}),
			"map_bool_fixed32":  toMapsOfAny(map[bool]uint32{true: 1}),
			"map_bool_fixed64":  toMapsOfAny(map[bool]uint64{true: 1}),
			"map_bool_sfixed32": toMapsOfAny(map[bool]int32{true: 1}),
			"map_bool_sfixed64": toMapsOfAny(map[bool]int64{true: 1}),
			"map_bool_float":    toMapsOfAny(map[bool]float32{true: 1}),
			"map_bool_double":   toMapsOfAny(map[bool]float64{true: 1}),
			"map_bool_bool":     toMapsOfAny(map[bool]bool{true: false}),
			"map_bool_string":   toMapsOfAny(map[bool]string{true: "1"}),
			"map_bool_bytes":    toMapsOfAny(map[bool][]byte{true: []byte("1")}),
			"map_bool_sibling_message": toMapsOfAny(map[bool]map[any]any{
				true: {"f1": "1"},
			}),
			"map_bool_sibling_enum": toMapsOfAny(map[bool]int32{true: 1}),
			"oneof_sibling_message": toMapsOfAny(map[string]string{
				"f1": "1",
			}),
		},
	},
	types: []proto.Message{&legacypb.Legacy{}},
}}

func TestMerge(t *testing.T) {
	for _, tt := range testMerges {
		for _, mt := range templateMessages(tt.types...) {
			t.Run(fmt.Sprintf("%s (%v)", tt.desc, mt.Descriptor().FullName()), func(t *testing.T) {
				dst := deepCopy(tt.dst).(map[any]any)

				src := mt.New().Interface()
				tt.src.Build(src.ProtoReflect())

				want := tt.want

				if mt.Descriptor().FullName() == "goproto.proto.test.TestAllExtensions" {
					want = extensionMap(want, mt)
					dst = extensionMap(dst, mt)
				}
				proto_.Merge(dst, src, tt.opts...)

				if !maps.EqualFunc(dst, want, reflect.DeepEqual) {
					t.Fatalf("proto_.Merge() mismatch:\ngot %v\nwant %v\ndiff (-want,+got):\n%v", dst, want, cmp.Diff(want, dst, protocmp.Transform()))
				}
				mutateValue(protoreflect.ValueOfMessage(src.ProtoReflect()))
				if !maps.EqualFunc(dst, want, reflect.DeepEqual) {
					t.Fatalf("mutation observed after modifying source:\n got %v\nwant %v\ndiff (-want,+got):\n%v", dst, want, cmp.Diff(want, dst, protocmp.Transform()))
				}
			})
		}
	}
}

func TestMergeFromNil(t *testing.T) {
	dst := map[any]any{}
	proto_.Merge(dst, (*testpb.TestAllTypes)(nil))
	if !maps.EqualFunc(dst, map[any]any{}, reflect.DeepEqual) {
		t.Errorf("destination should be empty after merging from nil message; got:\n%v", dst)
	}
}

// TestMergeAberrant tests inputs that are beyond the protobuf data model.
// Just because there is a test for the current behavior does not mean that
// this will behave the same way in the future.
func TestMergeAberrant(t *testing.T) {
	tests := []struct {
		label string
		dst   map[any]any
		src   proto.Message
		check func(map[any]any) bool
	}{{
		label: "Proto2EmptyBytes",
		dst:   map[any]any{"optional_bytes": nil},
		src:   &testpb.TestAllTypes{OptionalBytes: []byte{}},
		check: func(m map[any]any) bool { return m["optional_bytes"] != nil },
	}, {
		label: "Proto3EmptyBytes",
		dst:   map[any]any{"singular_bytes": nil},
		src:   &test3pb.TestAllTypes{SingularBytes: []byte{}},
		check: func(m map[any]any) bool { return m["singular_bytes"] == nil },
	}, {
		label: "EmptyList",
		dst:   map[any]any{"repeated_int32": nil},
		src:   &testpb.TestAllTypes{RepeatedInt32: []int32{}},
		check: func(m map[any]any) bool { return m["repeated_int32"] == nil },
	}, {
		label: "ListWithNilBytes",
		dst:   map[any]any{"repeated_bytes": nil},
		src:   &testpb.TestAllTypes{RepeatedBytes: [][]byte{nil}},
		check: func(m map[any]any) bool { return reflect.DeepEqual(m["repeated_bytes"], toSlicesOfAny([][]byte{{}})) },
	}, {
		label: "ListWithEmptyBytes",
		dst:   map[any]any{"repeated_bytes": nil},
		src:   &testpb.TestAllTypes{RepeatedBytes: [][]byte{{}}},
		check: func(m map[any]any) bool { return reflect.DeepEqual(m["repeated_bytes"], toSlicesOfAny([][]byte{{}})) },
	}, {
		label: "ListWithNilMessage",
		dst:   map[any]any{"repeated_nested_message": nil},
		src:   &testpb.TestAllTypes{RepeatedNestedMessage: []*testpb.TestAllTypes_NestedMessage{nil}},
		check: func(m map[any]any) bool { return m["repeated_nested_message"] != nil },
	}, {
		label: "EmptyMap",
		dst:   map[any]any{"map_string_string": nil},
		src:   &testpb.TestAllTypes{MapStringString: map[string]string{}},
		check: func(m map[any]any) bool { return m["map_string_string"] == nil },
	}, {
		label: "MapWithNilBytes",
		dst:   map[any]any{"map_string_bytes": nil},
		src:   &testpb.TestAllTypes{MapStringBytes: map[string][]byte{"k": nil}},
		check: func(m map[any]any) bool {
			return reflect.DeepEqual(m["map_string_bytes"], toMapsOfAny(map[string][]byte{"k": {}}))
		},
	}, {
		label: "MapWithEmptyBytes",
		dst:   map[any]any{"map_string_bytes": nil},
		src:   &testpb.TestAllTypes{MapStringBytes: map[string][]byte{"k": {}}},
		check: func(m map[any]any) bool {
			return reflect.DeepEqual(m["map_string_bytes"], toMapsOfAny(map[string][]byte{"k": {}}))
		},
	}, {
		label: "MapWithNilMessage",
		dst:   map[any]any{"map_string_nested_message": nil},
		src:   &testpb.TestAllTypes{MapStringNestedMessage: map[string]*testpb.TestAllTypes_NestedMessage{"k": nil}},
		check: func(m map[any]any) bool { return m["map_string_nested_message"].(map[any]any)["k"] != nil },
	}, {
		label: "OneofWithTypedNilWrapper",
		dst:   map[any]any{},
		src:   &testpb.TestAllTypes{OneofField: (*testpb.TestAllTypes_OneofNestedMessage)(nil)},
		check: func(m map[any]any) bool { return m["one_of_field"] == nil },
	}, {
		label: "OneofWithNilMessage",
		dst:   map[any]any{},
		src:   &testpb.TestAllTypes{OneofField: &testpb.TestAllTypes_OneofNestedMessage{OneofNestedMessage: nil}},
		check: func(m map[any]any) bool { return m["oneof_nested_message"] != nil },
		// TODO: extension, nil message
		// TODO: repeated extension, nil
		// TODO: extension bytes
		// TODO: repeated extension, nil message
	}}

	for _, tt := range tests {
		t.Run(tt.label, func(t *testing.T) {
			var pass bool
			func() {
				defer func() { recover() }()
				proto_.Merge(tt.dst, tt.src)
				pass = tt.check(tt.dst)
			}()
			if !pass {
				t.Errorf("check failed, got: %v", tt.dst)
			}
		})
	}
}

func TestMergeRace(t *testing.T) {
	dst := map[any]any{}
	srcs := []*testpb.TestAllTypes{
		{OptionalInt32: proto.Int32(1)},
		{OptionalString: proto.String("hello")},
		{RepeatedInt32: []int32{2, 3, 4}},
		{RepeatedString: []string{"goodbye"}},
		{MapStringString: map[string]string{"key": "value"}},
		{OptionalNestedMessage: &testpb.TestAllTypes_NestedMessage{
			A: proto.Int32(5),
		}},
		func() *testpb.TestAllTypes {
			m := new(testpb.TestAllTypes)
			m.ProtoReflect().SetUnknown(protopack.Message{
				protopack.Tag{Number: 50000, Type: protopack.VarintType}, protopack.Svarint(-5),
			}.Marshal())
			return m
		}(),
	}

	// It should be safe to concurrently merge non-overlapping fields.
	var wg sync.WaitGroup
	defer wg.Wait()
	for _, src := range srcs {
		wg.Add(1)
		go func(src proto.Message) {
			defer wg.Done()
			proto_.Merge(dst, src)
		}(src)
	}
}

func TestMergeSelf(t *testing.T) {
	got := &testpb.TestAllTypes{
		OptionalInt32:   proto.Int32(1),
		OptionalString:  proto.String("hello"),
		RepeatedInt32:   []int32{2, 3, 4},
		RepeatedString:  []string{"goodbye"},
		MapStringString: map[string]string{"key": "value"},
		OptionalNestedMessage: &testpb.TestAllTypes_NestedMessage{
			A: proto.Int32(5),
		},
	}
	got.ProtoReflect().SetUnknown(protopack.Message{
		protopack.Tag{Number: 50000, Type: protopack.VarintType}, protopack.Svarint(-5),
	}.Marshal())
	gotMapAnyAny := map[any]any{}
	proto_.Merge(gotMapAnyAny, got)
	proto_.Merge(gotMapAnyAny, got)

	// The main impact of merging to self is that repeated fields and
	// unknown fields are doubled.
	want := &testpb.TestAllTypes{
		OptionalInt32:   proto.Int32(1),
		OptionalString:  proto.String("hello"),
		RepeatedInt32:   []int32{2, 3, 4, 2, 3, 4},
		RepeatedString:  []string{"goodbye", "goodbye"},
		MapStringString: map[string]string{"key": "value"},
		OptionalNestedMessage: &testpb.TestAllTypes_NestedMessage{
			A: proto.Int32(5),
		},
	}
	want.ProtoReflect().SetUnknown(protopack.Message{
		protopack.Tag{Number: 50000, Type: protopack.VarintType}, protopack.Svarint(-5),
		protopack.Tag{Number: 50000, Type: protopack.VarintType}, protopack.Svarint(-5),
	}.Marshal())

	wantMapAnyAny := map[any]any{}
	proto_.Merge(wantMapAnyAny, want)

	if !maps.EqualFunc(gotMapAnyAny, wantMapAnyAny, reflect.DeepEqual) {
		t.Fatalf("proto_.Merge() mismatch:\ngot %v\nwant %v\ndiff (-want,+got):\n%v", gotMapAnyAny, wantMapAnyAny, cmp.Diff(wantMapAnyAny, gotMapAnyAny, protocmp.Transform()))
	}
}

func TestClone(t *testing.T) {
	want := &testpb.TestAllTypes{
		OptionalInt32: proto.Int32(1),
	}
	got := proto.Clone(want).(*testpb.TestAllTypes)
	if !proto.Equal(got, want) {
		t.Errorf("Clone(src) != src:\n got %v\nwant %v", got, want)
	}
}

// mutateValue changes a Value, returning a new value.
//
// For scalar values, it returns a value different from the input.
// For Message, List, and Map values, it mutates the input and returns it.
func mutateValue(v protoreflect.Value) protoreflect.Value {
	switch v := v.Interface().(type) {
	case bool:
		return protoreflect.ValueOfBool(!v)
	case protoreflect.EnumNumber:
		return protoreflect.ValueOfEnum(v + 1)
	case int32:
		return protoreflect.ValueOfInt32(v + 1)
	case int64:
		return protoreflect.ValueOfInt64(v + 1)
	case uint32:
		return protoreflect.ValueOfUint32(v + 1)
	case uint64:
		return protoreflect.ValueOfUint64(v + 1)
	case float32:
		return protoreflect.ValueOfFloat32(v + 1)
	case float64:
		return protoreflect.ValueOfFloat64(v + 1)
	case []byte:
		for i := range v {
			v[i]++
		}
		return protoreflect.ValueOfBytes(v)
	case string:
		return protoreflect.ValueOfString("_" + v)
	case protoreflect.Message:
		v.Range(func(fd protoreflect.FieldDescriptor, val protoreflect.Value) bool {
			v.Set(fd, mutateValue(val))
			return true
		})
		return protoreflect.ValueOfMessage(v)
	case protoreflect.List:
		for i := 0; i < v.Len(); i++ {
			v.Set(i, mutateValue(v.Get(i)))
		}
		return protoreflect.ValueOfList(v)
	case protoreflect.Map:
		v.Range(func(mk protoreflect.MapKey, mv protoreflect.Value) bool {
			v.Set(mk, mutateValue(mv))
			return true
		})
		return protoreflect.ValueOfMap(v)
	default:
		panic(fmt.Sprintf("unknown value type %T", v))
	}
}

func extensionMap(m map[any]any, mt protoreflect.MessageType) map[any]any {
	if m == nil || mt == nil {
		return m
	}
	res := make(map[any]any)
	par := mt.Descriptor().Parent().FullName()
	for k, v := range m {
		if vv, ok := v.(map[any]any); ok {
			v = extensionMap(vv, mt)
		}
		kk := fmt.Sprintf("%s", k)
		// all extensions is started with prefix "optional_"
		// See: ./internal/testprotos/test/test.proto, defined in extend TestAllExtensions
		if strings.HasPrefix(kk, "optional_") || strings.HasPrefix(kk, "repeated_") || strings.HasPrefix(kk, "default_") {
			kk = fmt.Sprintf("[%s.%s]", par, k)
		}
		res[kk] = v
	}
	return res
}

func deepCopy(item any) any {
	if item == nil {
		return nil
	}
	typ := reflect.TypeOf(item)
	val := reflect.ValueOf(item)
	if typ.Kind() == reflect.Ptr {
		newVal := reflect.New(typ.Elem())
		newVal.Elem().Set(reflect.ValueOf(deepCopy(val.Elem().Interface())))
		return newVal.Interface()
	} else if typ.Kind() == reflect.Map {
		newMap := reflect.MakeMap(typ)
		for _, k := range val.MapKeys() {
			newMap.SetMapIndex(k, reflect.ValueOf(deepCopy(val.MapIndex(k).Interface())))
		}
		return newMap.Interface()
	} else if typ.Kind() == reflect.Slice {
		newSlice := reflect.MakeSlice(typ, val.Len(), val.Cap())
		for i := 0; i < val.Len(); i++ {
			newSlice.Index(i).Set(reflect.ValueOf(deepCopy(val.Index(i).Interface())))
		}
		return newSlice.Interface()
	}
	return item
}

func toSlicesOfAny[S ~[]E, E any](s S) []any {
	return slices_.MapFunc(s, func(e E) any { return e })
}

func toMapsOfAny[M ~map[K]V, K comparable, V any](m M) map[any]any {
	return maps_.MapFunc(m, func(k K, v V) (any, any) { return k, v })
}
