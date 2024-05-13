// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package context

import (
	"context"
	"net/textproto"
)

// ExtractMIMIETags returns a pre-existing Tags object in the Context.
// If the context wasn't set in a tag interceptor, a no-op Tag storage is returned that will *not* be propagated in context.
func ExtractMIMIETags(ctx context.Context, key any) (tags textproto.MIMEHeader, has bool) {
	t, ok := ctx.Value(key).(textproto.MIMEHeader)
	if !ok {
		return textproto.MIMEHeader{}, false
	}

	return t, true
}

// ExtractOrCreateMIMETags extracts or create tags from context by key
func ExtractOrCreateMIMETags(ctx context.Context, key any) (
	ctx_ context.Context, stags textproto.MIMEHeader) {
	tags, has := ExtractMIMIETags(ctx, key)
	if has {
		return ctx, tags
	}
	tags = textproto.MIMEHeader{}
	return WithMIMETags(ctx, key, tags), tags
}

// WithMIMETags create tags from context by key
func WithMIMETags(ctx context.Context, key any, tags textproto.MIMEHeader) context.Context {
	return context.WithValue(ctx, key, tags)
}
