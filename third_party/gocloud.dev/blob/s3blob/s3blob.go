// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package s3blob

import (
	"context"
	"net/url"

	"gocloud.dev/blob"
)

var _ blob.BucketURLOpener = (BucketURLOpenerFunc)(nil)

// The BucketURLOpenerFunc type is an adapter to allow the use of
// ordinary functions as HTTP handlers. If f is a function
// with the appropriate signature, BucketURLOpenerFunc(f) is a
// Handler that calls f.
type BucketURLOpenerFunc func(ctx context.Context, u *url.URL) (*blob.Bucket, error)

// OpenBucketURL calls f(ctx, u).
func (f BucketURLOpenerFunc) OpenBucketURL(ctx context.Context, u *url.URL) (*blob.Bucket, error) {
	return f(ctx, u)
}
