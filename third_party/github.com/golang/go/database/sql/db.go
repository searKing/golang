// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	"github.com/searKing/golang/go/database/dsn"
	time_ "github.com/searKing/golang/go/time"
)

var (
	ErrRegister = errors.New("register db driver")
)

// DB represents a connection to a SQL database.
//go:generate go-option -type "DB"
type DB struct {
	DSN           string
	db            *sqlx.DB
	driverName    string
	driverPackage string

	logger logrus.FieldLogger

	// options
	opts struct {
		// UseTracedDriver will make it so that a wrapped driver is used that supports the opentracing API.
		// Deprecated: remove trace options.
		UseTracedDriver bool
		// if TraceOrphans is set to true, then spans with no parent will be traced anyway, if false, they will not be.
		// Deprecated: remove trace options.
		TraceOrphans bool
		// if OmitArgs is set to true, then query arguments are omitted from tracing spans.
		// Deprecated: remove trace options.
		OmitArgs bool
		// ForcedDriverName is specifically for writing tests as you can't register a driver with the same name more than once.
		// Deprecated: remove trace options.
		ForcedDriverName string
	}
}

// Open returns a new DB.
func Open(dsn string, opts ...DBOption) (*DB, error) {
	connection := &DB{DSN: dsn}
	connection.ApplyOptions(opts...)

	return connection, nil
}

func cleanURLQuery(in url.Values) (out url.Values) {
	out, _ = url.ParseQuery(in.Encode())
	out.Del("max_conns")
	out.Del("max_idle_conns")
	out.Del("max_conn_lifetime")
	out.Del("parseTime")
	return out
}

func (db *DB) fieldLogger() logrus.FieldLogger {
	if db.logger == nil {
		return logrus.StandardLogger()
	}
	return db.logger

}

// GetDatabaseRetry tries to connect to a database and fails after failAfter.
func (db *DB) GetDatabaseRetry(ctx context.Context, maxWait time.Duration, failAfter time.Duration) (*sqlx.DB, error) {
	// how long to sleep on retry failure
	var tempDelay = time_.NewDefaultExponentialBackOff(
		time_.WithExponentialBackOffOptionMaxInterval(maxWait),
		time_.WithExponentialBackOffOptionMaxElapsedDuration(failAfter))
	var err error
	for {
		db.db, err = db.GetDatabase()
		if err != nil {
			if errors.Is(err, ErrRegister) {
				db.fieldLogger().WithError(err).Errorf("database: Register")
				return nil, err
			}

			delay, ok := tempDelay.NextBackOff()
			if !ok {
				db.fieldLogger().WithError(err).
					Errorf("database: Connect; retried canceled as time exceed(%v)",
						tempDelay.GetMaxElapsedDuration())
				return nil, err
			}
			db.fieldLogger().WithError(err).Errorf("database: Connect; retrying in (%v)", delay)
			time.Sleep(delay)
			continue
		}
		return db.db, nil
	}
}

// GetDatabase returns a database instance.
func (db *DB) GetDatabase() (*sqlx.DB, error) {
	if db.db != nil {
		return db.db, nil
	}

	driverName, driverPackage, err := db.registerDriver()
	if err != nil {
		return nil, fmt.Errorf("could not register driver: %w", err)
	}

	dsn_, err := dsnForSqlOpen(db.DSN)
	if err != nil {
		return nil, err
	}

	classifiedDSN := dsn.Masking(dsn_)
	db.fieldLogger().WithField("dsn", classifiedDSN).Info("Establishing connection with SQL database backend")

	db_, err := sql.Open(driverName, dsn_)
	if err != nil {
		db.fieldLogger().WithError(err).WithField("dsn", classifiedDSN).Error("Unable to open SQL connection")
		return nil, fmt.Errorf("could not open SQL connection: %w", err)
	}

	db.db = sqlx.NewDb(db_, driverPackage) // This must be clean.Scheme otherwise things like `Rebind()` won't work
	if err := db.db.Ping(); err != nil {
		db.fieldLogger().WithError(err).WithField("dsn", classifiedDSN).Error("Unable to ping SQL database backend")
		return nil, fmt.Errorf("could not ping SQL connection: %w", err)
	}

	db.fieldLogger().WithField("dsn", classifiedDSN).Info("Successfully connected to SQL database backend")

	_, _, query := dsn.Split(dsn_)

	maxConns := maxParallelism() * 2
	if v := query.Get("max_conns"); v != "" {
		s, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			db.fieldLogger().WithError(err).Warnf(`Query parameter "max_conns" value %v could not be parsed to int, falling back to default value %d`, v, maxConns)
		} else {
			maxConns = int(s)
		}
	}

	maxIdleConns := maxParallelism()
	if v := query.Get("max_idle_conns"); v != "" {
		s, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			db.fieldLogger().Warnf("max_idle_conns value %s could not be parsed to int: %s", v, err)
			db.fieldLogger().WithError(err).Warnf(`Query parameter "max_idle_conns" value %v could not be parsed to int, falling back to default value %d`, v, maxIdleConns)
		} else {
			maxIdleConns = int(s)
		}
	}

	maxConnLifetime := time.Duration(0)
	if v := query.Get("max_conn_lifetime"); v != "" {
		s, err := time.ParseDuration(v)
		if err != nil {
			db.fieldLogger().WithError(err).Warnf(`Query parameter "max_conn_lifetime" value %v could not be parsed to int, falling back to default value %d`, v, maxConnLifetime)
		} else {
			maxConnLifetime = s
		}
	}

	db.db.SetMaxOpenConns(maxConns)
	db.db.SetMaxIdleConns(maxIdleConns)
	db.db.SetConnMaxLifetime(maxConnLifetime)

	return db.db, nil
}

func maxParallelism() int {
	maxProcs := runtime.GOMAXPROCS(0)
	numCPU := runtime.NumCPU()
	if maxProcs < numCPU {
		return maxProcs
	}
	return numCPU
}

func dsnForSqlOpen(dsn_ string) (string, error) {
	scheme, connect, query := dsn.Split(dsn_)

	query = cleanURLQuery(query)
	// special case, remove scheme for mysql*
	if strings.HasPrefix(scheme, "mysql") {
		query.Set("parseTime", "true")
		scheme = ""
	}
	if scheme == "cockroach" {
		scheme = "postgres"
	}

	return dsn.Join(scheme, connect, query), nil
}

// registerDriver checks if tracing is enabled and registers a custom "instrumented-sql-driver" driver that internally
// wraps the proper driver (mysql/postgres) with an instrumented driver.
func (db *DB) registerDriver() (string, string, error) {
	if db.driverName != "" {
		return db.driverName, db.driverPackage, nil
	}

	scheme, _, _ := dsn.Split(db.DSN)
	driverName := scheme
	driverPackage := scheme

	if driverName == "cockroach" || driverName == "postgres" {
		// If we're not using the instrumented driver, we need to replace "cockroach" with "postgres"
		driverName = "pgx"
	}

	switch scheme {
	case "cockroach":
		driverPackage = "postgres"
	}

	db.driverName = driverName
	db.driverPackage = driverPackage
	if driverName == "" {
		return "", "", fmt.Errorf("unsupported scheme (%s) in DSN : %w", scheme, ErrRegister)
	}
	return driverName, driverPackage, nil
}
