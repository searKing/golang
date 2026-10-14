// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package structpb_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/searKing/golang/third_party/google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/encoding/protojson"
)

type Human struct {
	Name      string
	Friends   []string
	Strangers []Human
}

type ToProtoStructTests struct {
	input Human
	want  string
}

var (
	toProtoStructTests = []ToProtoStructTests{{
		input: Human{
			Name:    "Alice",
			Friends: []string{"Bob", "Carol", "Dave"},
			Strangers: []Human{{
				Name:    "Eve",
				Friends: []string{"Oscar"},
				Strangers: []Human{{
					Name:    "Isaac",
					Friends: []string{"Justin", "Trent", "Pat", "Victor", "Walter"},
				}},
			}},
		},
		want: `{
 "Friends": [
  "Bob",
  "Carol",
  "Dave"
 ],
 "Name": "Alice",
 "Strangers": [
  {
   "Friends": [
    "Oscar"
   ],
   "Name": "Eve",
   "Strangers": [
    {
     "Friends": [
      "Justin",
      "Trent",
      "Pat",
      "Victor",
      "Walter"
     ],
     "Name": "Isaac",
     "Strangers": null
    }
   ]
  }
 ]
}`,
	},
	}
)

func TestToProtoStruct(t *testing.T) {
	for m, tt := range toProtoStructTests {
		humanStructpb, err := structpb.ToProtoStruct(tt.input)
		if err != nil {
			t.Errorf("#%d: ToProtoStruct(%+v): got: _, %v exp: _, nil", m, tt.input, err)
		}

		marshal := protojson.MarshalOptions{EmitUnpopulated: false, Indent: " ", UseProtoNames: true}
		humanBytes, err := marshal.Marshal(humanStructpb)

		if err != nil {
			t.Errorf("#%d: json.Marshal(%+v): got: _, %v exp: _, nil", m, tt.input, err)
		}
		// protojson adds spaces randomly, against depending on its output byte by byte, so compare JSON values.
		var got, want any
		if err := json.Unmarshal(humanBytes, &got); err != nil {
			t.Fatalf("#%d: json.Unmarshal(%s): %v", m, humanBytes, err)
		}
		if err := json.Unmarshal([]byte(tt.want), &want); err != nil {
			t.Fatalf("#%d: json.Unmarshal(%s): %v", m, tt.want, err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("#%d: json.Marshal(%+v): got(%dB): %v want(%dB): %v", m, tt.input, len(humanBytes), string(humanBytes), len(tt.want), tt.want)
		}
	}
}
