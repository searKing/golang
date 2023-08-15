// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !go1.9
// +build !go1.9

package otlpsql

import (
	"database/sql/driver"
	"errors"
)

// Dummy error for setSpanStatus (does exist as sql.ErrConnDone in 1.9+)
var errConnDone = errors.New("database/sql: connection is already closed")

// otlpDriver implements driver.Driver
type otlpDriver struct {
	parent  driver.Driver
	options wrapper
}

func wrapDriver(d driver.Driver, o wrapper) driver.Driver {
	return otlpDriver{parent: d, options: o}
}

func wrapConn(c driver.Conn, options wrapper) driver.Conn {
	return &otlpConn{parent: c, options: options}
}

func wrapStmt(stmt driver.Stmt, query string, options wrapper) driver.Stmt {
	s := otlpStmt{parent: stmt, query: query, options: options}
	_, hasExeCtx := stmt.(driver.StmtExecContext)
	_, hasQryCtx := stmt.(driver.StmtQueryContext)
	c, hasColCnv := stmt.(driver.ColumnConverter)
	switch {
	case !hasExeCtx && !hasQryCtx && !hasColCnv:
		return struct {
			driver.Stmt
		}{s}
	case !hasExeCtx && hasQryCtx && !hasColCnv:
		return struct {
			driver.Stmt
			driver.StmtQueryContext
		}{s, s}
	case hasExeCtx && !hasQryCtx && !hasColCnv:
		return struct {
			driver.Stmt
			driver.StmtExecContext
		}{s, s}
	case hasExeCtx && hasQryCtx && !hasColCnv:
		return struct {
			driver.Stmt
			driver.StmtExecContext
			driver.StmtQueryContext
		}{s, s, s}
	case !hasExeCtx && !hasQryCtx && hasColCnv:
		return struct {
			driver.Stmt
			driver.ColumnConverter
		}{s, c}
	case !hasExeCtx && hasQryCtx && hasColCnv:
		return struct {
			driver.Stmt
			driver.StmtQueryContext
			driver.ColumnConverter
		}{s, s, c}
	case hasExeCtx && !hasQryCtx && hasColCnv:
		return struct {
			driver.Stmt
			driver.StmtExecContext
			driver.ColumnConverter
		}{s, s, c}
	case hasExeCtx && hasQryCtx && hasColCnv:
		return struct {
			driver.Stmt
			driver.StmtExecContext
			driver.StmtQueryContext
			driver.ColumnConverter
		}{s, s, s, c}
	}
	panic("unreachable")
}
