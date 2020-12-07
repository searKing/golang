// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package time_test

import (
	"context"
	"testing"
	"time"

	time_ "github.com/searKing/golang/go/time"
)

func TestUntil(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	time_.Until(ctx, func(ctx context.Context) {
		t.Fatal("should not have been invoked")
	}, 0)

	ctx, cancel = context.WithCancel(context.Background())
	called := make(chan struct{})
	go func() {
		time_.Until(ctx, func(ctx context.Context) {
			called <- struct{}{}
		}, 0)
		close(called)
	}()
	<-called
	cancel()
	<-called
}

func TestNonSlidingUntil(t *testing.T) {
	ctx, cancel := context.WithCancel(context.TODO())
	cancel()
	time_.NonSlidingUntil(ctx, func(context.Context) {
		t.Fatal("should not have been invoked")
	}, 0)

	ctx, cancel = context.WithCancel(context.TODO())
	called := make(chan struct{})
	go func() {
		time_.NonSlidingUntil(ctx, func(context.Context) {
			called <- struct{}{}
		}, 0)
		close(called)
	}()
	<-called
	cancel()
	<-called
}

func TestUntilReturnsImmediately(t *testing.T) {
	ctx, cancel := context.WithCancel(context.TODO())

	now := time.Now()
	time_.Until(ctx, func(ctx context.Context) {
		cancel()
	}, 30*time.Second)
	if now.Add(25 * time.Second).Before(time.Now()) {
		t.Errorf("Until did not return immediately when the stop chan was closed inside the func")
	}
}

func TestJitterUntil(t *testing.T) {
	ctx, cancel := context.WithCancel(context.TODO())
	cancel()
	time_.JitterUntil(ctx, func(context.Context) {
		t.Fatal("should not have been invoked")
	}, true,
		time_.WithExponentialBackOffOptionRandomizationFactor(0.5),
		time_.WithExponentialBackOffOptionMultiplier(1),
		time_.WithExponentialBackOffOptionInitialInterval(0))

	ctx, cancel = context.WithCancel(context.TODO())
	called := make(chan struct{})
	go func() {
		time_.JitterUntil(ctx, func(context.Context) {
			called <- struct{}{}
		}, true,
			time_.WithExponentialBackOffOptionRandomizationFactor(0.5),
			time_.WithExponentialBackOffOptionMultiplier(1),
			time_.WithExponentialBackOffOptionInitialInterval(0))
		close(called)
	}()
	<-called
	cancel()
	<-called
}

func TestJitterUntilReturnsImmediately(t *testing.T) {
	ctx, cancel := context.WithCancel(context.TODO())

	now := time.Now()
	time_.JitterUntil(ctx, func(ctx context.Context) {
		cancel()
	}, true,
		time_.WithExponentialBackOffOptionRandomizationFactor(0.5),
		time_.WithExponentialBackOffOptionMultiplier(1),
		time_.WithExponentialBackOffOptionInitialInterval(30*time.Second))
	if now.Add(25 * time.Second).Before(time.Now()) {
		t.Errorf("JitterUntil did not return immediately when the stop chan was closed inside the func")
	}
}
