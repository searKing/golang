// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package url

import (
	"encoding/json"
	"net/url"
	"time"
)

func ParseBoolFromValues(q url.Values, key string) (bool, error) {
	if !q.Has(key) {
		return false, nil
	}
	s := q.Get(key)
	if s == "" {
		return true, nil
	}
	var b bool
	err := json.Unmarshal([]byte(s), &b)
	if err != nil {
		return false, err
	}
	return b, nil
}

func ParseTimeDurationFromValues(q url.Values, key string) (time.Duration, error) {
	if !q.Has(key) {
		return 0, nil
	}
	s := q.Get(key)
	if s == "" {
		return 0, nil
	}
	var b time.Duration
	err := json.Unmarshal([]byte(s), &b)
	if err != nil {
		return 0, err
	}
	return b, nil
}
