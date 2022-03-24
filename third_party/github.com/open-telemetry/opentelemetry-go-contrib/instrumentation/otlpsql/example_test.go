// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package otlpsql_test

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/mattn/go-sqlite3"
	"github.com/searKing/golang/third_party/github.com/open-telemetry/opentelemetry-go-contrib/instrumentation/otlpsql"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	_ "github.com/mattn/go-sqlite3"
)

func init() {
	otlpsql.PostCall = func(ctx context.Context, err error, elapse time.Duration, attrs ...attribute.KeyValue) {
		span := trace.SpanFromContext(ctx)
		log.Printf("trace_id: %s,space_id: %s, %v cost: %s",
			span.SpanContext().TraceID().String(), span.SpanContext().SpanID().String(), attrs, elapse)
	}
}

func ExampleDB_QueryContext() {
	{
		// Register our sqlite3-otlp wrapper for the provided SQLite3 driver.
		// "sqlite3-otlp" must not be registered, set in func init(){} as recommended.
		sql.Register("sqlite3-otlp", otlpsql.Wrap(&sqlite3.SQLiteDriver{}, otlpsql.WithAllWrapperOptions()))
	}

	// Connect to a SQLite3 database using the otlpsql driver wrapper.
	db, err := sql.Open("sqlite3-otlp", "resource.db")
	if err != nil {
		log.Fatal(err)
	}

	age := 27
	rows, err := db.QueryContext(context.Background(), "SELECT name FROM users WHERE age=?", age)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	names := make([]string, 0)

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			// Check for a scan error.
			// Query rows will be closed with defer.
			log.Fatal(err)
		}
		names = append(names, name)
	}
	// If the database is being written to ensure to check for Close
	// errors that may be returned from the driver. The query may
	// encounter an auto-commit error and be forced to rollback changes.
	rerr := rows.Close()
	if rerr != nil {
		log.Fatal(rerr)
	}

	// Rows.Err will report the last error encountered by Rows.Scan.
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s are %d years old", strings.Join(names, ", "), age)
}
