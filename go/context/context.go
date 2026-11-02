// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package context

import "context"

// WithCancelCauseContexts behaves like [context.WithCancelCause] but also
// cancels the returned context when any of causes is done.
//
// The returned context's Done channel is closed when the returned cancel
// function is called, when the parent context's Done channel is closed, or
// when the Done channel of any of causes is closed, whichever happens first.
// [context.Cause] of the returned context reports the cause of whichever of
// these happened first. Nil causes are ignored.
//
// If any of causes is already done, the returned context is canceled before
// WithCancelCauseContexts returns.
//
// Canceling this context releases resources associated with it, including
// those registered with causes, so code should call cancel as soon as the
// operations running in this [context.Context] complete.
func WithCancelCauseContexts(parent context.Context, causes ...context.Context) (ctx context.Context, cancel context.CancelCauseFunc) {
	ctx, cancel = context.WithCancelCause(parent)

	var stops []func() bool
	for _, cause := range causes {
		if cause == nil {
			continue
		}
		if cause.Err() != nil {
			cancel(context.Cause(cause))
			break
		}
		stops = append(stops, context.AfterFunc(cause, func() { cancel(context.Cause(cause)) }))
	}
	if len(stops) > 0 {
		context.AfterFunc(ctx, func() {
			for _, stop := range stops {
				stop()
			}
		})
	}
	return ctx, cancel
}
