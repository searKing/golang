// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package version_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/searKing/golang/go/version"
)

func TestVersion_String(t *testing.T) {
	testCases := []struct {
		ver    version.Version
		expect string
	}{
		{
			ver: version.Version{
				Major:      1,
				Minor:      2,
				Patch:      3,
				PreRelease: "fix",
				BuildTime:  time.Now().String(),
				GitHash:    "0xFFFF",
				MetaData:   []string{"M", "E", "T", "A"},
			},
			expect: "v1.2.3-fix(0xFFFF)",
		},
		{
			ver: version.Version{
				Major:      1,
				Minor:      2,
				Patch:      3,
				PreRelease: "devel",
				BuildTime:  "NOW",
				GitHash:    "0xFFFF",
				MetaData:   []string{"M", "E", "T", "A"},
			},
			expect: "v1.2.3-devel(0xFFFF)+M.E.T.A",
		},
		{
			ver: version.Version{
				Major:      1,
				Minor:      2,
				Patch:      3,
				RawVersion: "v7.8.9",
				PreRelease: "fix",
				BuildTime:  "NOW",
				GitHash:    "0xFFFF",
				MetaData:   []string{"M", "E", "T", "A"},
			},
			expect: "v7.8.9",
		},
	}
	for i, tc := range testCases {
		if tc.ver.String() != tc.expect {
			t.Errorf("#%d expect %s, got %s", i, tc.expect, tc.ver.String())
		}
	}
}

func TestVersion_Format(t *testing.T) {
	testCases := []struct {
		ver    version.Version
		fmt    string
		expect string
	}{
		{
			ver: version.Version{
				Major:      1,
				Minor:      2,
				Patch:      3,
				PreRelease: "fix",
				BuildTime:  "NOW",
				GitHash:    "0xFFFF",
				MetaData:   []string{"M", "E", "T", "A"},
			},
			fmt:    "%s",
			expect: "v1.2.3-fix(0xFFFF)",
		},
		{
			ver: version.Version{
				Major:      1,
				Minor:      2,
				Patch:      3,
				PreRelease: "fix",
				BuildTime:  "NOW",
				GitHash:    "0xFFFF",
				MetaData:   []string{"M", "E", "T", "A"},
			},
			fmt:    "%q",
			expect: "v1.2.3-fix(0xFFFF)",
		},
		{
			ver: version.Version{
				Major:      1,
				Minor:      2,
				Patch:      3,
				PreRelease: "fix",
				BuildTime:  "NOW",
				GitHash:    "0xFFFF",
				MetaData:   []string{"M", "E", "T", "A"},
			},
			fmt:    "%v",
			expect: "v1.2.3-fix(0xFFFF)",
		},
		{
			ver: version.Version{
				Major:      1,
				Minor:      2,
				Patch:      3,
				PreRelease: "fix",
				BuildTime:  "NOW",
				GitHash:    "0xFFFF",
				MetaData:   []string{"M", "E", "T", "A"},
			},
			fmt:    "%+v",
			expect: "v1.2.3-fix(0xFFFF), Build #gc-go1.15.6-darwin/amd64, built on NOW",
		},
	}
	for i, tc := range testCases {
		ver := fmt.Sprintf(tc.fmt, tc.ver)
		if ver != tc.expect {
			t.Errorf("#%d expect %s, got %s", i, tc.expect, ver)
		}

	}
}
