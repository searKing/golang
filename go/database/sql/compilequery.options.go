// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sql

// WithCompileQueryOptionAliasWithSelect Generate Alias
// `SELECT t.a, b`
// TO
// `select t.a as t_a, b as b`,
func WithCompileQueryOptionAliasWithSelect(aliasWithSelect bool) CompileQueryOption {
	return CompileQueryOptionFunc(func(opt *compileQuery) {
		opt.AliasWithSelect = aliasWithSelect
	})
}

// WithCompileQueryOptionArgument keep column if argument by column name is not zero
// take effect in WHERE|INSERT|UPDATE, ignore if multi rows
// nil: keep all
// []string: keep if exist
// map[string]{{value}} : keep if exist and none zero
// struct{} tag is `db:"{{col_name}}"`: keep if exist and none zero
//
// `SELECT t.a, b FROM t WHERE first_name= :hehe AND middle_name=? OR last_name=?`
// TO
// `select t.a as t_a, b as b from t where first_name = :first_name or last_name = :last_name`,
func WithCompileQueryOptionArgument(arg interface{}) CompileQueryOption {
	return CompileQueryOptionFunc(func(opt *compileQuery) {
		opt.Argument = arg
	})
}
