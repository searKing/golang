// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package context_test

import (
	"context"
	"errors"
	"fmt"

	context_ "github.com/searKing/golang/go/context"
)

// This example ties the lifetime of a request to both the server's and the
// client's contexts: the work stops as soon as either of them is done, and
// the cause reports which one it was.
func ExampleWithCancelCauseContexts() {
	server, stopServer := context.WithCancelCause(context.Background())
	defer stopServer(nil)
	client, disconnect := context.WithCancelCause(context.Background())
	defer disconnect(nil)

	ctx, cancel := context_.WithCancelCauseContexts(server, client)
	defer cancel(nil)

	disconnect(errors.New("client disconnected"))

	<-ctx.Done()
	fmt.Println(ctx.Err())
	fmt.Println(context.Cause(ctx))

	// Output:
	// context canceled
	// client disconnected
}
