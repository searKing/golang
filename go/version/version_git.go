// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package version

import (
	"os"
	"path/filepath"
)

var (
	// GitTag
	// NOTE: The $Format strings are replaced during 'git archive' thanks to the
	// companion .gitattributes file containing 'export-subst' in this same
	// directory.  See also https://git-scm.com/docs/gitattributes
	GitTag    = "v0.0.0-master+$Format:%h$" // git describe --long --tags --dirty --tags --always
	BuildTime = "1970-01-01T00:00:00Z"      // build date in ISO8601 format, output of $(date -u +'%Y-%m-%dT%H:%M:%SZ')
	GitHash   = "$Format:%H$"               // sha1 from git, output of $(git rev-parse HEAD)

	ServiceName        = filepath.Base(os.Args[0])
	ServiceDisplayName = filepath.Base(os.Args[0])
	ServiceDescription = ""
	ServiceId          = ""
)

// Example
// git_tag=$(shell git describe --long --tags --dirty --tags --always)
// git_commit=$(shell git rev-parse HEAD)
// git_build_time=$(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
// go build -gcflags=all="-N -l" \
// -ldflags "-s -X 'github.com/searKing/golang/go/version.GitTag=${git_tag}' \
// -X 'github.com/searKing/golang/go/version.BuildTime=${git_build_time}' \
// -X 'github.com/searKing/golang/go/version.GitHash=${git_commit}'"

// Get returns GitTag as version
func Get() Version {
	return Version{
		RawVersion: GitTag,
		BuildTime:  BuildTime,
		GitHash:    GitHash,
	}
}
