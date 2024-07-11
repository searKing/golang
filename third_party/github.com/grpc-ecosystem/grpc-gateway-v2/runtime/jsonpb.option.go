// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import "google.golang.org/protobuf/encoding/protojson"

func WithMarshalOptions(option protojson.MarshalOptions) JSONPbOption {
	return JSONPbOptionFunc(func(pb *JSONPb) {
		pb.MarshalOptions = option
	})
}

func WithUnmarshalOptions(option protojson.UnmarshalOptions) JSONPbOption {
	return JSONPbOptionFunc(func(pb *JSONPb) {
		pb.UnmarshalOptions = option
	})
}

// WithUseEnumNumbers Whether to render enum values as integers, as opposed to string values.
func WithUseEnumNumbers(useEnumNumbers bool) JSONPbOption {
	return JSONPbOptionFunc(func(pb *JSONPb) {
		pb.UseEnumNumbers = useEnumNumbers
	})
}

// WithEmitUnpopulated Whether to render fields with zero values.
func WithEmitUnpopulated(emitUnpopulated bool) JSONPbOption {
	return JSONPbOptionFunc(func(pb *JSONPb) {
		pb.EmitUnpopulated = emitUnpopulated
	})
}

// WithIndent A string to indent each level by. The presence of this field will
// also cause a space to appear between the field separator and
// value, and for newlines to be appear between fields and array
// elements.
func WithIndent(indent string) JSONPbOption {
	return JSONPbOptionFunc(func(pb *JSONPb) {
		pb.Indent = indent
	})
}

// WithUseProtoNames Whether to use the original (.proto) name for fields.
func WithUseProtoNames(useProtoNames bool) JSONPbOption {
	return JSONPbOptionFunc(func(pb *JSONPb) {
		pb.UseProtoNames = useProtoNames
	})
}

// Deprecated: Use WithUseEnumNumbers instead.
func WithEmitAsInts(enumsAsInts bool) JSONPbOption {
	return WithUseEnumNumbers(enumsAsInts)
}

// Deprecated: Use WithEmitUnpopulated instead.
func WithEmitDefaults(emitDefaults bool) JSONPbOption {
	return WithEmitUnpopulated(emitDefaults)
}

// Deprecated: Use WithUseProtoNames instead.
func WithOrigName(origName bool) JSONPbOption {
	return WithUseProtoNames(origName)
}

// WithUnmarshalAllowPartial If AllowPartial is set, input for messages that will result in missing
// required fields will not return an error.
func WithUnmarshalAllowPartial(allowPartial bool) JSONPbOption {
	return JSONPbOptionFunc(func(pb *JSONPb) {
		pb.UnmarshalOptions.AllowPartial = allowPartial
	})
}

// WithDiscardUnknown If DiscardUnknown is set, unknown fields are ignored.
func WithDiscardUnknown(discardUnknown bool) JSONPbOption {
	return JSONPbOptionFunc(func(pb *JSONPb) {
		pb.DiscardUnknown = discardUnknown
	})
}
