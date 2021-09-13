// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package binding

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	protov1 "github.com/golang/protobuf/proto"
	"github.com/searKing/golang/third_party/github.com/golang/protobuf/jsonpb"
	"github.com/searKing/golang/third_party/google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// JSONPB encode json to proto.Message
var JSONPB = jsonpbBinding{}

type jsonpbBinding struct{}

func (jsonpbBinding) Name() string {
	return "jsonpb"
}

func (b jsonpbBinding) Bind(req *http.Request, obj interface{}) error {
	buf, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}
	return b.BindBody(buf, obj)
}

func (jsonpbBinding) BindBody(body []byte, obj interface{}) error {
	if msg, ok := obj.(proto.Message); ok {
		if err := protojson.Unmarshal(body, msg); err != nil {
			return err
		}
	} else if msg, ok := obj.(protov1.Message); ok {
		if err := jsonpb.Unmarshal(body, msg); err != nil {
			return err
		}
	} else {
		if err := json.Unmarshal(body, obj); err != nil {
			return err
		}
	}
	// Here it's same to return validate(obj), but util now we can't add
	// `binding:""` to the struct which automatically generate by gen-proto
	return nil
	// return validate(obj)
}
