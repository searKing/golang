// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package version records versioning information about this module.

// borrowed from https://github.com/protocolbuffers/protobuf-go/blob/v1.25.0/internal/version/version.go

package version

import (
	"fmt"
	"io"
	"runtime"
	"strings"
)

// These constants determine the current version of this module.
//
//
// For our release process, we enforce the following rules:
//	* Tagged releases use a tag that is identical to String.
//	* Tagged releases never reference a commit where the String
//	contains "devel".
//	* The set of all commits in this repository where String
//	does not contain "devel" must have a unique String.
//
//
// Steps for tagging a new release:
//	1. Create a new CL.
//
//	2. Update Minor, Patch, and/or PreRelease as necessary.
//	PreRelease must not contain the string "devel".
//
//	3. Since the last released minor version, have there been any changes to
//	generator that relies on new functionality in the runtime?
//	If yes, then increment RequiredGenerated.
//
//	4. Since the last released minor version, have there been any changes to
//	the runtime that removes support for old .pb.go source code?
//	If yes, then increment SupportMinimum.
//
//	5. Send out the CL for review and submit it.
//	Note that the next CL in step 8 must be submitted after this CL
//	without any other CLs in-between.
//
//	6. Tag a new version, where the tag is is the current String.
//
//	7. Write release notes for all notable changes
//	between this release and the last release.
//
//	8. Create a new CL.
//
//	9. Update PreRelease to include the string "devel".
//	For example: "" -> "devel" or "rc.1" -> "rc.1.devel"
//
//	10. Send out the CL for review and submit it.

type Version struct {
	Major      int
	Minor      int
	Patch      int
	PreRelease string

	// NOTE: The $Format strings are replaced during 'git archive' thanks to the
	// companion .gitattributes file containing 'export-subst' in this same
	// directory.  See also https://git-scm.com/docs/gitattributes
	// git describe --long --tags --dirty --tags --always
	// "v0.0.0-master+$Format:%h$"
	RawVersion string

	// build date in ISO8601 format, output of $(date -u +'%Y-%m-%dT%H:%M:%SZ')
	// "1970-01-01T00:00:00Z"
	BuildTime string
	// sha1 from git, output of $(git rev-parse HEAD)
	// "$Format:%H$"
	GitHash string

	MetaData []string // take effect when PreRelease contains devel
}

// String formats the version string for this module in semver format.
//
// Examples:
//
//	v1.20.1
//	v1.21.0-rc.1
func (ver Version) String() string {
	if ver.RawVersion != "" {
		return ver.RawVersion
	}
	v := fmt.Sprintf("v%d.%d.%d", ver.Major, ver.Minor, ver.Patch)
	if ver.PreRelease != "" {
		v += "-" + ver.PreRelease
	}
	if ver.GitHash != "" {
		v += "(" + ver.GitHash + ")"
	}
	// TODO: Add metadata about the commit or build hash.
	// See https://golang.org/issue/29814
	// See https://golang.org/issue/33533
	var metadata = strings.Join(ver.MetaData, ".")
	if strings.Contains(ver.PreRelease, "devel") && metadata != "" {
		v += "+" + metadata
	}
	return v
}

func (ver Version) BuildInfo() string {
	//	GoVersion = runtime.Version()
	//	Compiler  = runtime.Compiler
	//	Platform  = fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
	return fmt.Sprintf("%s-%s-%s",
		runtime.Compiler, runtime.Version(),
		fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH))
}

// Format Examples:
// v1.2.3-fix, Build #gc-go1.15.6-darwin/amd64, built on NOW
func (ver Version) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = io.WriteString(s, ver.String())
			if buildInfo := ver.BuildInfo(); buildInfo != "" {
				_, _ = io.WriteString(s, ", Build #"+buildInfo)
			}

			if ver.BuildTime != "" {
				_, _ = io.WriteString(s, ", built on "+ver.BuildTime)
			}
			return
		}
		fallthrough
	case 's', 'q':
		_, _ = io.WriteString(s, ver.String())
	}
}
