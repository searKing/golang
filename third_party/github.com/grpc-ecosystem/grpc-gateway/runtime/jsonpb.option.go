// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

// WithEmitAsInts specifies whether to render enum values as integers, as opposed to string values.
func WithEmitAsInts(enumsAsInts bool) JSONPbOption {
	return JSONPbOptionFunc(func(pb *JSONPb) {
		pb.JSONPb.EnumsAsInts = enumsAsInts
	})
}

// WithEmitDefaults specifies whether to render fields with zero values.
func WithEmitDefaults(emitDefaults bool) JSONPbOption {
	return JSONPbOptionFunc(func(pb *JSONPb) {
		pb.JSONPb.EmitDefaults = emitDefaults
	})
}

// WithIndent controls whether the output is compact or not.
func WithIndent(indent string) JSONPbOption {
	return JSONPbOptionFunc(func(pb *JSONPb) {
		pb.JSONPb.Indent = indent
	})
}

// Whether to use the original (.proto) name for fields.
func WithOrigName(origName bool) JSONPbOption {
	return JSONPbOptionFunc(func(pb *JSONPb) {
		pb.JSONPb.OrigName = origName
	})
}
