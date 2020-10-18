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

// WithCompileQueryOptionTrimByColumn trim node by column name
// take effect in WHERE|INSERT|UPDATE, ignore if multi rows
// `SELECT t.a, b FROM t WHERE first_name= :hehe AND middle_name=? OR last_name=?`
// TO
// `select t.a as t_a, b as b from t where first_name = :first_name or last_name = :last_name`,
func WithCompileQueryOptionTrimByColumn(trimByColumn map[string]interface{}) CompileQueryOption {
	return CompileQueryOptionFunc(func(opt *compileQuery) {
		opt.TrimByColumn = trimByColumn
	})
}
