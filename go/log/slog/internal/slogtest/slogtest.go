// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slogtest

import "log/slog"

// RemoveTime removes the top-level time attribute.
// It is intended to be used as a ReplaceAttr function,
// to make example output deterministic.
func RemoveTime(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.TimeKey && len(groups) == 0 {
		return slog.Attr{}
	}
	return a
}
