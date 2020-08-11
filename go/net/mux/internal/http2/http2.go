// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http2

import (
	"io"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/hpack"
)

const (
	// http://http2.github.io/http2-spec/#SettingValues
	initialHeaderTableSize = 4096
)

// MatchHTTP2Header matches all headerFields if len(filterNames) == 0
func MatchHTTP2Header(w io.Writer, r io.Reader, filterNames map[string]struct{}, matches func(headerFields map[string]hpack.HeaderField) (matched bool)) (matched bool) {
	// filter http2 only
	if !HasClientPreface(r) {
		return false
	}

	readAll := len(filterNames) > 0
	done := false
	framer := http2.NewFramer(w, r)
	var filteredHeaderFields map[string]hpack.HeaderField
	readMetaHeaders := hpack.NewDecoder(initialHeaderTableSize, func(f hpack.HeaderField) {
		if _, has := filterNames[f.Name]; has || readAll {
			filteredHeaderFields[f.Name] = f
		}
	})

	frame, err := framer.ReadFrame()
	if err != nil {
		return false
	}

	sf, ok := frame.(*http2.SettingsFrame)
	if !ok {
		return false
	}
	if err := handleSettings(framer, sf); err != nil {
		return false
	}

	for {
		frame, err := framer.ReadFrame()
		if err != nil {
			return false
		}

		switch frame := frame.(type) {
		case *http2.SettingsFrame:
			if err := handleSettings(framer, sf); err != nil {
				return false
			}
		case *http2.ContinuationFrame:
			if _, err := readMetaHeaders.Write(frame.HeaderBlockFragment()); err != nil {
				return false
			}
			done = frame.HeadersEnded()
		case *http2.HeadersFrame:
			if _, err := readMetaHeaders.Write(frame.HeaderBlockFragment()); err != nil {
				return false
			}
			done = frame.HeadersEnded()
		}

		if done || (!readAll && len(filteredHeaderFields) == len(filterNames)) {
			return matches(filteredHeaderFields)
		}
	}
}
func handleSettings(framer *http2.Framer, frame *http2.SettingsFrame) error {
	// Sender acknowledged the SETTINGS frame. No need to write
	// SETTINGS again.
	if frame.IsAck() {
		return nil
	}
	return framer.WriteSettings()
}
