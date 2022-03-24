// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package otlpsql

import "database/sql/driver"

// Compile time validation that our types implement the expected interfaces
var _ driver.NamedValueChecker = otlpStmt{}

func (s otlpStmt) CheckNamedValue(v *driver.NamedValue) error {
	if checker, ok := s.parent.(driver.NamedValueChecker); ok {
		return checker.CheckNamedValue(v)
	}

	return driver.ErrSkip
}
