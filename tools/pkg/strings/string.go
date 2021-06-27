// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package strings

func LoadElse(loaded bool, loadValue string, elseValue string) string {
	if loaded {
		return loadValue
	}
	return elseValue
}

func LoadElseGet(loaded bool, loadValue string, elseValueGetter func() string) string {
	if loaded {
		return loadValue
	}
	if elseValueGetter == nil {
		return ""
	}
	return elseValueGetter()
}
