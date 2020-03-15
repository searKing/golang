// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package json

import (
	"encoding/json"
	"io/ioutil"
)

func ReadConfigFile(name string, v interface{}) error {
	data, err := ioutil.ReadFile(name)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}
