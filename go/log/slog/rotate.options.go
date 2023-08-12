// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slog

import (
	time_ "github.com/searKing/golang/go/time"
)

// WithRotateFilePathRotateStrftime sets time layout in strftime format to format rotate file.
func WithRotateFilePathRotateStrftime(layout string) RotateOption {
	return RotateOptionFunc(func(r *rotate) {
		r.FilePathRotateLayout = time_.LayoutStrftimeToSimilarTime(layout)
	})
}
