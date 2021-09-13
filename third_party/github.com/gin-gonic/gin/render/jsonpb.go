// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package render

import (
	"encoding/json"
	"net/http"

	protov1 "github.com/golang/protobuf/proto"
	"github.com/searKing/golang/third_party/github.com/golang/protobuf/jsonpb"
	"github.com/searKing/golang/third_party/google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// JSONPB contains the given interface object.
type JSONPB struct {
	Data interface{}
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
func WriteJSONPB(w http.ResponseWriter, obj interface{}) error {
	writeContentType(w, jsonpbContentType)
	var jsonBytes []byte
	var err error
	if msg, ok := obj.(proto.Message); ok {
		jsonBytes, err = protojson.Marshal(msg)
	} else if msg, ok := obj.(protov1.Message); ok {
		jsonBytes, err = jsonpb.Marshal(msg)
	} else {
		jsonBytes, err = json.Marshal(obj)
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
