// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sqlx

// WithCompileQueryOptionAliasWithSelect Generate Alias
// `SELECT t.a, b`
// TO
// `select t.a as t_a, b as b`,
func WithCompileQueryOptionAliasWithSelect(aliasWithSelect bool) CompileQueryOption {
	return CompileQueryOptionFunc(func(opt *compileQuery) {
		opt.AliasWithSelect = aliasWithSelect
	})
}
