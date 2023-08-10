// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package prettyjson

import "io"

// An Encoder writes JSON values to an output stream.
type Encoder struct {
	w          io.Writer
	err        error
	escapeHTML bool

	indentBuf    []byte
	indentPrefix string
	indentValue  string

	truncateBytes  int
	truncateString int
	truncateMap    int
	truncateSlice  int
	truncateArray  int
}

// NewEncoder returns a new encoder that writes to w.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w, escapeHTML: true}
}

// Encode writes the JSON encoding of v to the stream,
// followed by a newline character.
//
// See the documentation for Marshal for details about the
// conversion of Go values to JSON.
func (enc *Encoder) Encode(v any) error {
	if enc.err != nil {
		return enc.err
	}

	e := newEncodeState()
	defer encodeStatePool.Put(e)

	err := e.marshal(v, encOpts{escapeHTML: enc.escapeHTML,
		truncateBytes:        enc.truncateBytes,
		truncateString:       enc.truncateString,
		truncateMap:          enc.truncateMap,
		truncateSliceOrArray: enc.truncateSlice})
	if err != nil {
		return err
	}

	// Terminate each value with a newline.
	// This makes the output look a little nicer
	// when debugging, and some kind of space
	// is required if the encoded value was a number,
	// so that the reader knows there aren't more
	// digits coming.
	e.WriteByte('\n')

	b := e.Bytes()
	if enc.indentPrefix != "" || enc.indentValue != "" {
		enc.indentBuf, err = appendIndent(enc.indentBuf[:0], b, enc.indentPrefix, enc.indentValue)
		if err != nil {
			return err
		}
		b = enc.indentBuf
	}
	if _, err = enc.w.Write(b); err != nil {
		enc.err = err
	}
	return err
}

// SetIndent instructs the encoder to format each subsequent encoded
// value as if indented by the package-level function Indent(dst, src, prefix, indent).
// Calling SetIndent("", "") disables indentation.
func (enc *Encoder) SetIndent(prefix, indent string) {
	enc.indentPrefix = prefix
	enc.indentValue = indent
}

// SetEscapeHTML specifies whether problematic HTML characters
// should be escaped inside JSON quoted strings.
// The default behavior is to escape &, <, and > to \u0026, \u003c, and \u003e
// to avoid certain safety problems that can arise when embedding JSON in HTML.
//
// In non-HTML settings where the escaping interferes with the readability
// of the output, SetEscapeHTML(false) disables this behavior.
func (enc *Encoder) SetEscapeHTML(on bool) {
	enc.escapeHTML = on
}

// SetTruncate shrinks bytes, string, map, slice, array 's len to n at most
func (enc *Encoder) SetTruncate(n int) {
	enc.truncateBytes = n
	enc.truncateString = n
	enc.truncateMap = n
	enc.truncateSlice = n
	enc.truncateArray = n
}
