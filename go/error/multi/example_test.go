// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package multi_test

import (
	"fmt"

	"github.com/searKing/golang/go/error/multi"
)

func ExampleNew() {
	err := multi.New(fmt.Errorf("whoops"), fmt.Errorf("foo"))
	fmt.Println(err)
	fmt.Println("--------")
	fmt.Printf("%+v", err)

	// Output: whoops|foo
	// --------
	// Multiple errors occurred:
	//	whoops|foo
}
