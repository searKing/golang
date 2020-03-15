// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// errNoOverlap is returned by serveContent's parseContentRange if first-byte-pos of
// all of the byte-range-spec values is greater than the content size.
var errNoOverlap = errors.New("invalid range: failed to overlap")

// httpContentRange specifies the byte range to be sent to the client.
type httpContentRange struct {
	firstBytePos, lastBytePos, completeLength int64
}

func (r httpContentRange) String() string {
	var b strings.Builder
	b.WriteString("bytes ")
	if r.firstBytePos == -1 && r.lastBytePos == -1 {
		b.WriteString("*")
	} else {
		if r.firstBytePos == -1 {
			b.WriteString("*")
		} else {
			b.WriteString(fmt.Sprintf("%d", r.firstBytePos))
		}
		b.WriteString("-")
		if r.lastBytePos == -1 {
			b.WriteString("*")
		} else {
			b.WriteString(fmt.Sprintf("%d", r.lastBytePos))
		}
	}
	b.WriteString("/")

	if r.completeLength < 0 {
		b.WriteString("*")
	} else {
		b.WriteString(fmt.Sprintf("%d", r.completeLength))
	}
	return b.String()
}

func parseContentRanges(s []string) ([]httpContentRange, error) {
	var ranges []*httpContentRange

	for _, ra := range s {
		ra = strings.TrimSpace(ra)
		if ra == "" {
			continue
		}
		r, err := parseContentRange(ra)
		if err != nil {
			return nil, err
		}
		ranges = append(ranges, r)
	}

	var size int64 = -1
	for _, ra := range ranges {
		if size != -1 && ra.completeLength != -1 && size != ra.completeLength {
			return nil, errors.New("invalid range")
		}
		if ra.completeLength >= 0 {
			size = ra.completeLength
		}
	}

	for _, ra := range ranges {
		ra.completeLength = size
		if ra.firstBytePos < 0 {
			ra.firstBytePos = 0
		}
		if ra.lastBytePos < 0 {
			if size < 0 {
				return nil, errors.New("invalid range")
			}
			ra.lastBytePos = size
		}
	}

	var totalSize int64
	for _, ra := range ranges {
		totalSize += ra.lastBytePos - ra.firstBytePos + 1
	}

	if size >= 0 && size != totalSize {
		return nil, errors.New("invalid range")
	}

	var outRanges []httpContentRange
	for _, ra := range ranges {
		ra.completeLength = totalSize
		outRanges = append(outRanges, *ra)
	}

	return outRanges, nil
}

// parseContentRange parses a Range header string as per RFC 7233.
// errNoOverlap is returned if none of the ranges overlap.
func parseContentRange(s string) (contentRange *httpContentRange, err error) {
	// bytes 0-499/1234
	if s == "" {
		return nil, nil // header not present
	}
	const b = "bytes "
	if !strings.HasPrefix(s, b) {
		return nil, errors.New("invalid range")
	}

	ra := strings.TrimSpace(s[len(b):])
	if ra == "" {
		return nil, nil
	}

	i := strings.Index(ra, "/")
	if i < 0 {
		return nil, errors.New("invalid range")
	}
	byteRange := strings.TrimSpace(ra[:i])
	completeLength := strings.TrimSpace(ra[i+1:])

	if byteRange == "*" && completeLength == "*" {
		return nil, errors.New("invalid range")
	}
	if byteRange == "" && completeLength == "" {
		return nil, errors.New("invalid range")
	}
	if byteRange == "*" && completeLength == "" {
		return nil, errors.New("invalid range")
	}
	if byteRange == "" && completeLength == "*" {
		return nil, errors.New("invalid range")
	}

	var r = httpContentRange{
		firstBytePos:   -1,
		lastBytePos:    -1,
		completeLength: -1,
	}
	if byteRange != "*" {
		i := strings.Index(byteRange, "-")
		if i < 0 {
			return nil, errors.New("invalid range")
		}
		firstBytePos := strings.TrimSpace(byteRange[:i])
		lastBytePos := strings.TrimSpace(byteRange[i+1:])
		if firstBytePos != "" {
			i, err := strconv.ParseInt(firstBytePos, 10, 64)
			if err != nil {
				return nil, errors.New("invalid range")
			}
			r.firstBytePos = i
		}
		if lastBytePos != "" {
			i, err := strconv.ParseInt(lastBytePos, 10, 64)
			if err != nil {
				return nil, errors.New("invalid range")
			}
			r.lastBytePos = i
		}
	}

	if completeLength != "*" {
		i, err := strconv.ParseInt(completeLength, 10, 64)
		if err != nil {
			return nil, errors.New("invalid range")
		}
		r.completeLength = i
	}

	if r.firstBytePos < 0 && r.lastBytePos < 0 && r.completeLength < 0 {
		return nil, errors.New("invalid range")
	}

	if r.firstBytePos >= 0 && r.lastBytePos >= 0 && r.firstBytePos > r.lastBytePos {
		return nil, errors.New("invalid range")
	}

	if r.firstBytePos >= 0 && r.lastBytePos >= 0 && r.completeLength >= 0 && (r.lastBytePos-r.firstBytePos+1 < r.completeLength) {
		// The specified ranges did not overlap with the content.
		return nil, errNoOverlap
	}
	return &r, nil
}

// countingWriter counts how many bytes have been written to it.
type countingWriter int64

func (w *countingWriter) Write(p []byte) (n int, err error) {
	*w += countingWriter(len(p))
	return len(p), nil
}

func writeContentRanges(w http.ResponseWriter, ranges []httpContentRange) {
	for _, ra := range ranges {
		w.Header().Add("Content-Range", ra.String())
	}
}
