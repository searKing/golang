// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cause_test

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/searKing/golang/go/error/cause"
)

func ExampleWithError() {
	cause_ := errors.New("whoops")
	err := fmt.Errorf("oh noes")
	err = cause.WithError(cause_, err)
	fmt.Println(err)

	// Output: oh noes: whoops
}
