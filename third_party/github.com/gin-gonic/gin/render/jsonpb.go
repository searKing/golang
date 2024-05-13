// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package render

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/golang/protobuf/jsonpb"
	protov1 "github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/encoding/protojson"
	protov2 "google.golang.org/protobuf/proto"
)

// JSONPB contains the given interface object.
type JSONPB struct {
	Data any
}

var jsonpbContentType = []string{"application/json; charset=utf-8"}

// Render (JSON) writes data with custom ContentType.
func (r JSONPB) Render(w http.ResponseWriter) (err error) {
	if err = WriteJSONPB(w, r.Data); err != nil {
		panic(err)
	}
	return
}

// WriteContentType (JSON) writes JSON ContentType.
func (r JSONPB) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, jsonpbContentType)
}

// WriteJSONPB marshals the given interface object and writes it with custom ContentType.
func WriteJSONPB(w http.ResponseWriter, obj any) error {
	writeContentType(w, jsonpbContentType)
	var jsonBytes []byte
	var err error
	switch ee := obj.(type) {
	case protov1.Message:
		mm := jsonpb.Marshaler{}
		var buf bytes.Buffer
		err := mm.Marshal(&buf, ee)
		if err != nil {
			// This may fail for proto.Anys, e.g. for xDS v2, LDS, the v2
			// messages are not imported, and this will fail because the message
			// is not found.
			return err
		}
		jsonBytes = buf.Bytes()
	case protov2.Message:
		mm := protojson.MarshalOptions{AllowPartial: true}
		jsonBytes, err = mm.Marshal(ee)
		if err != nil {
			// This may fail for proto.Anys, e.g. for xDS v2, LDS, the v2
			// messages are not imported, and this will fail because the message
			// is not found.
			return err
		}
	default:
		jsonBytes, err = json.Marshal(ee)
	}

	if err != nil {
		return err
	}
	_, err = w.Write(jsonBytes)
	return err
}

func writeContentType(w http.ResponseWriter, value []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	}
}
