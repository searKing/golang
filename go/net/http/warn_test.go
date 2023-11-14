// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http

import (
	"strings"
	"testing"
	"time"
)

type respWriteTest struct {
	Resp Warn
	Raw  string
}

func TestWarn_Write(t *testing.T) {
	respWriteTests := []respWriteTest{
		{
			Warn{
				WarnCode: WarnTransformationApplied,
			},

			`214 - "Transformation Applied"`,
		},
		{
			Warn{
				Warn:     "cache down",
				WarnCode: WarnDisconnectedOperation,
				Date:     time.Date(2015, 10, 21, 7, 28, 0, 0, time.UTC),
			},

			`112 - "cache down" "Wed, 21 Oct 2015 07:28:00 GMT"`,
		},
		{
			Warn{
				Warn:     "Arbitrary information that should be presented to a user or logged.",
				WarnCode: WarnMiscellaneousWarning,
				Agent:    "Go-http-client/1.1",
			},

			`199 Go-http-client/1.1 "Arbitrary information that should be presented to a user or logged."`,
		},
		{
			Warn{
				Warn: "Arbitrary information that should be presented to a user or logged. " +
					"This warn-code is similar to the warn-code 199 and additionally indicates a persistent warning.",
				WarnCode: WarnMiscellaneousPersistentWarning,
				Date:     time.Date(2011, 11, 23, 1, 5, 3, 0, time.UTC),
			},

			`299 - "Arbitrary information that should be presented to a user or logged. This warn-code is similar to the warn-code 199 and additionally indicates a persistent warning." "Wed, 23 Nov 2011 01:05:03 GMT"`,
		},
	}

	for i := range respWriteTests {
		tt := &respWriteTests[i]
		var braw strings.Builder
		err := tt.Resp.Write(&braw)
		if err != nil {
			t.Errorf("error writing #%d: %s", i, err)
			continue
		}
		sraw := braw.String()
		if sraw != tt.Raw {
			t.Errorf("Test %d, expecting:\n%q\nGot:\n%q\n", i, tt.Raw, sraw)
			continue
		}
	}
}
