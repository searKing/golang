// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package encoding_test

import (
	"testing"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/sourcecontextpb"

	"github.com/searKing/golang/pkg/webserver/pkg/encoding"
)

func TestNewJSONPb(t *testing.T) {
	tests := []struct {
		name             string
		opts             []encoding.JSONPbOption
		wantMarshal      string
		wantUnmarshalErr bool
	}{
		{name: "default", wantMarshal: `{}`},
		{name: "nil option", opts: []encoding.JSONPbOption{nil}, wantMarshal: `{}`},
		{
			name: "custom",
			opts: []encoding.JSONPbOption{
				encoding.WithMarshalOptions(protojson.MarshalOptions{EmitUnpopulated: true}),
				encoding.WithUnmarshalOptions(protojson.UnmarshalOptions{}),
			},
			wantMarshal:      `{"fileName":""}`,
			wantUnmarshalErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := encoding.NewJSONPb(tt.opts...)
			got, err := j.Marshal(&sourcecontextpb.SourceContext{})
			if err != nil {
				t.Fatal(err)
			}
			if string(got) != tt.wantMarshal {
				t.Errorf("Marshal = %s, want %s", got, tt.wantMarshal)
			}

			err = j.Unmarshal([]byte(`{"unknown":1}`), &sourcecontextpb.SourceContext{})
			if gotErr := err != nil; gotErr != tt.wantUnmarshalErr {
				t.Errorf("Unmarshal unknown field error = %v, want error %v", err, tt.wantUnmarshalErr)
			}
		})
	}
}
