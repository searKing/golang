// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package any_test

import (
	"strings"
	"testing"

	"github.com/golang/protobuf/jsonpb"
	"github.com/searKing/golang/third_party/github.com/golang/protobuf/ptypes/any"
)

type Human struct {
	Name      string
	Friends   []string
	Strangers []Human
}

type ToProtoAnyTests struct {
	input  Human
	output string
}

var (
	toProtoAnyTests = []ToProtoAnyTests{{
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
		output: `{
 "@type": "type.googleapis.com/google.protobuf.Struct",
 "value": {
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
  }
}`,
	},
	}
)

func TestToProtoAny(t *testing.T) {
	for m, test := range toProtoAnyTests {
		humanStructpb, err := any.ToProtoAny(test.input)
		if err != nil {
			t.Errorf("#%d: ToProtoAny(%+v): got: _, %v exp: _, nil", m, test.input, err)
		}

		marshal := jsonpb.Marshaler{EmitDefaults: false, Indent: " ", OrigName: true}
		humanStr, err := marshal.MarshalToString(humanStructpb)

		if err != nil {
			t.Errorf("#%d: json.Marshal(%+v): got: _, %v exp: _, nil", m, test.input, err)
		}

		if strings.Compare(humanStr, test.output) != 0 {
			t.Errorf("#%d: json.Marshal(%+v): got: %v exp: %v", m, test.input, humanStr, test.output)
		}
	}
}
