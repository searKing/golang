// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package time

import (
	"strings"
	"testing"
	"time"
)

// formatByChunks formats t chunk by chunk, as split by nextStdChunk.
// It equals t.Format(layout) only if nextStdChunk splits layout as package time does.
func formatByChunks(t time.Time, layout string) string {
	var b strings.Builder
	for layout != "" {
		prefix, std, suffix := nextStdChunk(layout)
		b.WriteString(prefix)
		if std == 0 {
			break
		}
		b.WriteString(t.Format(layout[len(prefix) : len(layout)-len(suffix)]))
		layout = suffix
	}
	return b.String()
}

// FuzzNextStdChunk detects drift between the copied nextStdChunk and package time of the running Go.
func FuzzNextStdChunk(f *testing.F) {
	for _, layout := range []string{
		time.Layout, time.ANSIC, time.UnixDate, time.RubyDate, time.RFC822, time.RFC822Z,
		time.RFC850, time.RFC1123, time.RFC1123Z, time.RFC3339, time.RFC3339Nano, time.Kitchen,
		time.Stamp, time.StampMilli, time.StampMicro, time.StampNano, time.DateTime, time.DateOnly, time.TimeOnly,
		"January Jan Janx Monday Mon Monx MST",
		"_2006 _2 __2 002 _x __x",
		"1 15 2 2006 3 4 5 01 02 03 04 05 06 07 00",
		"PM pm Pm pM",
		"-070000 -07:00:00 -0700 -07:00 -07 -0",
		"Z070000 Z07:00:00 Z0700 Z07:00 Z07 Z0",
		".0 .000 .999999999 ,000 ,999 .0001x 05.0001 .9a ,",
	} {
		f.Add(layout)
	}
	// Every std chunk formats ref to a text other than the chunk itself, e.g. "PM" to "AM".
	ref := time.Date(2021, 3, 12, 7, 49, 8, 236589123, time.FixedZone("CST", 8*3600+30*60+15))
	f.Fuzz(func(t *testing.T, layout string) {
		if got, want := formatByChunks(ref, layout), ref.Format(layout); got != want {
			t.Errorf("layout %q: formatted by chunks %q, want %q", layout, got, want)
		}
	})
}
