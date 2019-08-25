package http2

import (
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/hpack"
	"io"
)

const (
	// http://http2.github.io/http2-spec/#SettingValues
	initialHeaderTableSize = 4096
)

func MatchHTTP2Header(w io.Writer, r io.Reader, names map[string]struct{}, matches func(headerFields map[string]hpack.HeaderField) (matched bool)) (matched bool) {
	// filter http2 only
	if !HasClientPreface(w, r) {
		return false
	}
	done := false
	framer := http2.NewFramer(w, r)
	var filteredHeaderFields map[string]hpack.HeaderField
	readMetaHeaders := hpack.NewDecoder(initialHeaderTableSize, func(f hpack.HeaderField) {
		if _, has := names[f.Name]; has {
			filteredHeaderFields[f.Name] = f
		}
	})
	for {
		frame, err := framer.ReadFrame()
		if err != nil {
			return false
		}

		switch frame := frame.(type) {
		case *http2.SettingsFrame:
			// Sender acknowledged the SETTINGS frame. No need to write
			// SETTINGS again.
			if frame.IsAck() {
				break
			}
			if err := framer.WriteSettings(); err != nil {
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

		if done || len(filteredHeaderFields) == len(names) {
			return matches(filteredHeaderFields)
		}
	}
}
