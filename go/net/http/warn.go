// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// HTTP warn codes as registered with IANA.
// See: https://www.iana.org/assignments/http-warn-codes/http-warn-codes.xhtml
// See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Warning#warning_codes
// Deprecated: The "Warning" header field was used to carry additional information about
// the status or transformation of a message that might not be reflected in the status code.
// This specification obsoletes it, as it is not widely generated or surfaced to users.
// The information it carried can be gleaned from examining other header fields, such as Age.
// See: https://www.rfc-editor.org/rfc/rfc9111.html#name-warning
const (
	WarnResponseIsStale                = 110 // RFC 7234, 5.5.1
	WarnRevalidationFailed             = 111 // RFC 7234, 5.5.2
	WarnDisconnectedOperation          = 112 // RFC 7234, 5.5.3
	WarnHeuristicExpiration            = 113 // RFC 7234, 5.5.4
	WarnMiscellaneousWarning           = 199 // RFC 7234, 5.5.5
	WarnTransformationApplied          = 214 // RFC 7234, 5.5.6
	WarnMiscellaneousPersistentWarning = 299 // RFC 7234, 5.5.7
)

// WarnText returns a text for the HTTP warn code. It returns the empty
// string if the code is unknown.
func WarnText(code int) string {
	switch code {
	case WarnResponseIsStale:
		return "Response is Stale"
	case WarnRevalidationFailed:
		return "Revalidation Failed"
	case WarnDisconnectedOperation:
		return "Disconnected Operation"
	case WarnHeuristicExpiration:
		return "Heuristic Expiration"
	case WarnMiscellaneousWarning:
		return "Miscellaneous Warning"
	case WarnTransformationApplied:
		return "Transformation Applied"
	case WarnMiscellaneousPersistentWarning:
		return "Miscellaneous Persistent Warning"
	default:
		return ""
	}
}

type Warn struct {
	Warn     string // e.g. "200 OK"
	WarnCode int    // e.g. 200
	Agent    string // e.g. "-"
	Date     time.Time
}

func (r Warn) String() string {
	var buf strings.Builder
	_ = r.Write(&buf)
	return buf.String()
}

func (r Warn) Write(w io.Writer) error {
	// Status line
	text := r.Warn
	if text == "" {
		text = WarnText(r.WarnCode)
		if text == "" {
			text = "warn code " + strconv.Itoa(r.WarnCode)
		}
	} else {
		// Just to reduce stutter, if user set w.Warn to "112 Disconnected Operation" and WarnCode to 112.
		// Not important.
		text = strings.TrimPrefix(text, strconv.Itoa(r.WarnCode)+" ")
	}

	agent := r.Agent
	if agent == "" {
		agent = "-"
	}

	// See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Warning#syntax
	// Warning: <warn-code> <warn-agent> <warn-text> [<warn-date>]
	// Warning: 110 anderson/1.3.37 "Response is stale"
	// Date: Wed, 21 Oct 2015 07:28:00 GMT
	// Warning: 112 - "cache down" "Wed, 21 Oct 2015 07:28:00 GMT"
	if _, err := fmt.Fprintf(w, "%03d %s %q", r.WarnCode, agent, text); err != nil {
		return err
	}
	r.Date.IsZero()
	if !r.Date.IsZero() {
		if _, err := fmt.Fprintf(w, " %q", r.Date.UTC().Format(http.TimeFormat)); err != nil {
			return err
		}
	}
	return nil
}
