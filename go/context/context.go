// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package context

import "fmt"

type key string

func (k key) String() string { return fmt.Sprintf("context key(%s-%p)", string(k), &k) }

func Key(name string) *key {
	return (*key)(&name)
}
