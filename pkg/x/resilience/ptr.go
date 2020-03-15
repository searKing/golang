// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resilience

type Ptr interface {
	Value() interface{} //actual instance
	Ready() error
	Close()
}
