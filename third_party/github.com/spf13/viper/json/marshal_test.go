// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package json_test

import (
	"bytes"
	"testing"

	json_ "github.com/searKing/golang/go/encoding/json"
)

func TestMarshal(t *testing.T) {
	tests := []struct {
		v         any
		wantJson  []byte
		wantError bool
	}{
		{
			v: map[int]string{
				0: "alice",
				1: "bob",
			},
			wantJson:  []byte(`{"0":"alice","1":"bob"}`),
			wantError: false,
		},
		{
			v: map[any]string{
				"0": "alice",
				"1": "bob",
			},
			wantJson:  []byte(`{"0":"alice","1":"bob"}`),
			wantError: false,
		},
		{
			v: map[any]string{
				0: "alice",
				1: "bob",
			},
			wantJson:  []byte{},
			wantError: true,
		},
	}

	for idx, tt := range tests {
		gotJson, err := json_.Marshal(tt.v)
		if (err != nil) != (tt.wantError) {
			t.Errorf("#%d: got err %s", idx, err)
			continue
		}
		if !bytes.Equal(tt.wantJson, gotJson) {
			t.Errorf("#%d: expected %s got %s", idx, tt.wantJson, gotJson)
		}
	}
}
