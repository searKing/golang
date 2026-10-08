// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package time

import (
	"strings"
)

// LayoutTimeToStrftime converts a layout of package time, such as "2006-01-02", to a strftime layout, such as "%Y-%m-%d".
// Chunks of layout without exact equivalents in strftime are kept as is.
func LayoutTimeToStrftime(layout string) string {
	return convertLayout(layout, nextStdChunk, std.StrftimeString)
}

// LayoutTimeToSimilarStrftime is like LayoutTimeToStrftime,
// but converts chunks of layout without exact equivalents in strftime to similar ones,
// and drops fractional seconds.
func LayoutTimeToSimilarStrftime(layout string) string {
	return convertLayout(layout, nextStdChunk, std.SimilarStrftimeString)
}

// LayoutStrftimeToTime converts a strftime layout, such as "%Y-%m-%d", to a layout of package time, such as "2006-01-02".
// Conversion specifications without exact equivalents in package time, such as "%U", are kept as is.
func LayoutStrftimeToTime(layout string) string {
	return convertLayout(layout, nextStrftimeChunk, std.String)
}

// LayoutStrftimeToSimilarTime is like LayoutStrftimeToTime,
// but converts conversion specifications without exact equivalents in package time to similar ones,
// such as "%U" to "Mon".
func LayoutStrftimeToSimilarTime(layout string) string {
	return convertLayout(layout, nextStrftimeChunk, std.SimilarString)
}

// convertLayout splits layout into chunks by next, and converts each std chunk by name.
func convertLayout(layout string, next func(layout string) (prefix string, std int, suffix string), name func(std) string) string {
	var buf strings.Builder
	// Each iteration generates one std value.
	for layout != "" {
		prefix, std_, suffix := next(layout)
		buf.WriteString(prefix)
		if std_ == 0 {
			break
		}
		buf.WriteString(name(std(std_)))
		layout = suffix
	}
	return buf.String()
}

// nextStrftimeChunk finds the first occurrence of a strftime conversion specification in
// layout and returns the text before, the std value, and the text after.
// Unknown or incomplete conversion specifications, such as "%q" or a trailing "%", are kept as text.
func nextStrftimeChunk(layout string) (prefix string, std int, suffix string) {
	for i := 0; i < len(layout); i++ {
		if layout[i] != '%' {
			continue
		}
		j := i + 1
		chunks := &strftimeChunks
		if j < len(layout) {
			switch layout[j] {
			case 'E':
				chunks, j = &strftimeEChunks, j+1
			case 'O':
				chunks, j = &strftimeOChunks, j+1
			}
		}
		if j < len(layout) {
			if std := chunks[layout[j]]; std != 0 {
				return layout[:i], std, layout[j+1:]
			}
		}
	}
	return layout, 0, ""
}
