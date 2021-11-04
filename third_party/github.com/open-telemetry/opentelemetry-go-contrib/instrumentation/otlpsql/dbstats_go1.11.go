// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package otlpsql

import (
	"context"
	"database/sql"
	"sync"
	"time"
)

// RecordStats records database statistics for provided sql.DB at the provided
// interval.
func RecordStats(db *sql.DB, interval time.Duration) (fnStop func()) {
	var (
		closeOnce sync.Once
		ctx       = context.Background()
		ticker    = time.NewTicker(interval)
		done      = make(chan struct{})
	)

	go func() {
		for {
			select {
			case <-ticker.C:
				dbStats := db.Stats()
				MeasureOpenConnections.Record(ctx, int64(dbStats.OpenConnections))
				MeasureIdleConnections.Record(ctx, int64(dbStats.Idle))
				MeasureActiveConnections.Record(ctx, int64(dbStats.InUse))
				MeasureWaitCount.Record(ctx, dbStats.WaitCount)
				MeasureWaitDuration.Record(ctx, dbStats.WaitDuration.Milliseconds())
				MeasureIdleClosed.Record(ctx, dbStats.MaxIdleClosed)
				MeasureLifetimeClosed.Record(ctx, dbStats.MaxLifetimeClosed)
			case <-done:
				ticker.Stop()
				return
			}
		}
	}()

	return func() {
		closeOnce.Do(func() {
			close(done)
		})
	}
}
