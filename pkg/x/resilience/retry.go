// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package resilience provides helpers for dealing with resilience.
// This code is borrowed from https://github.com/ory/x/tree/master/resilience/retry.go
package resilience

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Retry executes a f until no error is returned or failAfter is reached.
// A failAfter Timeout of zero means no timeout.
// maxWait max interval of two f
func Retry(ctx context.Context, logger logrus.FieldLogger, maxWait time.Duration, failAfter time.Duration, f func() error) (err error) {
	if logger == nil {
		logger = logrus.StandardLogger()
	}

	err = errors.New("did not connect")
	waitReform := func(wait time.Duration) time.Duration {
		if wait > maxWait {
			wait = maxWait
		}
		return wait
	}

	loopWait := waitReform(time.Millisecond * 100)
	if failAfter != 0 {
		cancelCtx, cancelFn := context.WithTimeout(ctx, failAfter)
		defer cancelFn()
		ctx = cancelCtx
	} else {
		cancelCtx, cancelFn := context.WithCancel(ctx)
		defer cancelFn()
		ctx = cancelCtx
	}
L:
	for {
		start := time.Now()

		if err = f(); err == nil {
			return nil
		}
		logger.WithError(err).Warnf("retrying in %s seconds...", loopWait)
		select {
		case <-ctx.Done():
			break L
		case <-time.After(loopWait):
		}

		// task takes too much time, keep the step
		if time.Now().Before(start.Add(maxWait * 2)) {
			loopWait = waitReform(loopWait * 2)
		}

	}
	return err
}
