// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ice

import "fmt"

var portMap = map[string]string{
	"stun":  "3478",
	"turn":  "3478",
	"stuns": "5349",
	"turns": "5349",
}
var getDefaultPort = func(schema string) (string, error) {
	port, ok := portMap[schema]
	if ok {
		return port, nil
	}
	return "", fmt.Errorf("malformed schema:%s", schema)
}
