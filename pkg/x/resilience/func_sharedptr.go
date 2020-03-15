// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resilience

import (
	"context"

	"github.com/sirupsen/logrus"
)

func NewSharedPtrFunc(ctx context.Context,
	new func() (interface{}, error),
	ready func(x interface{}) error,
	close func(x interface{}), l logrus.FieldLogger) *SharedPtr {
	return NewSharedPtr(ctx, WithFuncNewer(new, ready, close), l)
}
