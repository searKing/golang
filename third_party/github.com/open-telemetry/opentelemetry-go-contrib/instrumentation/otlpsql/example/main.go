// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"os"
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
		logger := slog.Default()
		if span.SpanContext().HasTraceID() {
			logger = logger.With("trace_id", span.SpanContext().TraceID())
		}
		if span.SpanContext().HasSpanID() {
			logger = logger.With("spacn_id", span.SpanContext().SpanID())
		}
		logger.With(slog.Any("attrs", attrs)).With(slog.Duration("cost", elapse)).Info("")
	}
}

func TempFilename() string {
	f, err := os.CreateTemp("", "go-sqlite3-test-*.db")
	if err != nil {
		log.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func main() {
	{
		// Register our sqlite3-otlp wrapper for the provided SQLite3 driver.
		// "sqlite3-otlp" must not be registered, set in func init(){} as recommended.
		sql.Register("sqlite3-otlp", otlpsql.Wrap(&sqlite3.SQLiteDriver{}, otlpsql.WithAllWrapperOptions()))
	}

	_, n, _ := sqlite3.Version()
	if n < 3024000 {
		log.Fatal("UPSERT requires sqlite3 >= 3.24.0")
	}
	tempFilename := TempFilename()
	defer os.Remove(tempFilename)

	// Connect to a SQLite3 database using the otlpsql driver wrapper.
	db, err := sql.Open("sqlite3-otlp", tempFilename)
	if err != nil {
		log.Fatal(err)
	}

	{
		res, err := db.Exec("CREATE TABLE users (name string primary key , age integer)")
		if err != nil {
			log.Fatal(err)
		}
		_ = res
	}

	{
		res, err := db.Exec("insert into users(name, age) values('key', 3)")
		if err != nil {
			log.Fatal("Failed to insert record:", err)
		}
		affected, _ := res.RowsAffected()
		if affected != 1 {
			log.Fatalf("Expected %d for affected rows, but %d:", 1, affected)
		}
	}
	{
		res, err := db.Exec("insert into users(name, age) values('kid', 3)")
		if err != nil {
			log.Fatal("Failed to insert record:", err)
		}
		affected, _ := res.RowsAffected()
		if affected != 1 {
			log.Fatalf("Expected %d for affected rows, but %d:", 1, affected)
		}
	}
	{
		res, err := db.Exec("insert into users(name, age) values('adult', 27)")
		if err != nil {
			log.Fatal("Failed to insert record:", err)
		}
		affected, _ := res.RowsAffected()
		if affected != 1 {
			log.Fatalf("Expected %d for affected rows, but %d:", 1, affected)
		}
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
	fmt.Printf("%q are %d years old", strings.Join(names, ", "), age)

	// Output:
	// 2023/08/15 19:41:32 INFO  attrs="[{Key:db.type Value:{vtype:4 numeric:0 stringly:sql slice:<nil>}} {Key:db.statement Value:{vtype:4 numeric:0 stringly:CREATE TABLE users (name string primary key , age integer) slice:<nil>}}]" cost=878.5µs
	// 2023/08/15 19:41:32 INFO  attrs="[{Key:db.type Value:{vtype:4 numeric:0 stringly:sql slice:<nil>}} {Key:db.statement Value:{vtype:4 numeric:0 stringly:insert into users(name, age) values('key', 3) slice:<nil>}}]" cost=488.333µs
	// 2023/08/15 19:41:32 INFO  attrs="[{Key:db.type Value:{vtype:4 numeric:0 stringly:sql slice:<nil>}}]" cost=1.041µs
	// 2023/08/15 19:41:32 INFO  attrs="[{Key:db.type Value:{vtype:4 numeric:0 stringly:sql slice:<nil>}} {Key:db.statement Value:{vtype:4 numeric:0 stringly:insert into users(name, age) values('kid', 3) slice:<nil>}}]" cost=410.709µs
	// 2023/08/15 19:41:32 INFO  attrs="[{Key:db.type Value:{vtype:4 numeric:0 stringly:sql slice:<nil>}}]" cost=708ns
	// 2023/08/15 19:41:32 INFO  attrs="[{Key:db.type Value:{vtype:4 numeric:0 stringly:sql slice:<nil>}} {Key:db.statement Value:{vtype:4 numeric:0 stringly:insert into users(name, age) values('adult', 27) slice:<nil>}}]" cost=491.958µs
	// 2023/08/15 19:41:32 INFO  attrs="[{Key:db.type Value:{vtype:4 numeric:0 stringly:sql slice:<nil>}}]" cost=1.041µs
	// 2023/08/15 19:41:32 INFO  attrs="[{Key:db.type Value:{vtype:4 numeric:0 stringly:sql slice:<nil>}} {Key:db.statement Value:{vtype:4 numeric:0 stringly:SELECT name FROM users WHERE age=? slice:<nil>}} {Key:sql.arg.1 Value:{vtype:2 numeric:27 stringly: slice:<nil>}}]" cost=26.875µs
	// 2023/08/15 19:41:32 INFO  attrs="[{Key:db.type Value:{vtype:4 numeric:0 stringly:sql slice:<nil>}}]" cost=24.75µs
	// 2023/08/15 19:41:32 INFO  attrs="[{Key:db.type Value:{vtype:4 numeric:0 stringly:sql slice:<nil>}}]" cost=3.25µs
	// 2023/08/15 19:41:32 INFO  attrs="[{Key:db.type Value:{vtype:4 numeric:0 stringly:sql slice:<nil>}}]" cost=1.625µs
	// "adult" are 27 years old
}
