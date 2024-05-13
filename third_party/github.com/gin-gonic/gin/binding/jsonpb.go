// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package binding

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/golang/protobuf/jsonpb"
	protov1 "github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/encoding/protojson"
	protov2 "google.golang.org/protobuf/proto"
)

// JSONPB encode json to proto.Message
var JSONPB = jsonpbBinding{}

type jsonpbBinding struct{}

func (jsonpbBinding) Name() string {
	return "jsonpb"
}

func (b jsonpbBinding) Bind(req *http.Request, obj any) error {
	buf, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}
	return b.BindBody(buf, obj)
}

func (jsonpbBinding) BindBody(body []byte, obj any) error {
	switch msg := obj.(type) {
	case protov1.Message:
		mm := jsonpb.Unmarshaler{AllowUnknownFields: true}
		if err := mm.Unmarshal(bytes.NewBuffer(body), msg); err != nil {
			return err
		}
	case protov2.Message:
		mm := protojson.UnmarshalOptions{
			AllowPartial:   true,
			DiscardUnknown: true,
		}
		if err := mm.Unmarshal(body, msg); err != nil {
			return err
		}
	default:
		if err := json.Unmarshal(body, obj); err != nil {
			return err
		}
	}
	// Here it's same to return validate(obj), but util now we can't add
	// `binding:""` to the struct which automatically generate by gen-proto
	return nil
	// return validate(obj)
}
