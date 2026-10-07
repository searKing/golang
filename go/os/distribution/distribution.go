// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package distribution

import (
	"bufio"
	"os"
	"runtime"
	"strings"
)

type Distribution int

const (
	Windows Distribution = iota
	Ubuntu
	Redhat
	Centos
	Butt
)

func (d Distribution) String() string {
	switch d {
	case Windows:
		return "Windows"
	case Ubuntu:
		return "Ubuntu"
	case Redhat:
		return "Redhat"
	case Centos:
		return "Centos"
	default:
		return "Butt"
	}
}

// GetOSVersion returns the distribution of the running OS, Butt if unknown.
// On Linux, it reads ID in os-release(5), and falls back to /proc/version,
// which describes the kernel build, that of the host in a container.
func GetOSVersion() Distribution {
	if runtime.GOOS == "windows" {
		return Windows
	}

	for _, name := range []string{"/etc/os-release", "/usr/lib/os-release"} {
		data, err := os.ReadFile(name)
		if err == nil {
			return parseOSRelease(string(data))
		}
	}

	data, err := os.ReadFile("/proc/version")
	if err != nil {
		return Butt
	}
	return parseProcVersion(string(data))
}

// parseOSRelease parses the distribution from ID in os-release(5).
func parseOSRelease(data string) Distribution {
	s := bufio.NewScanner(strings.NewReader(data))
	for s.Scan() {
		id, ok := strings.CutPrefix(strings.TrimSpace(s.Text()), "ID=")
		if !ok {
			continue
		}
		switch strings.ToLower(strings.Trim(id, `"'`)) {
		case "ubuntu":
			return Ubuntu
		case "rhel":
			return Redhat
		case "centos":
			return Centos
		}
		return Butt
	}
	return Butt
}

// parseProcVersion parses the distribution from /proc/version.
func parseProcVersion(data string) Distribution {
	versionInfo := strings.ToLower(data)

	switch {
	// CentOS kernels are built by gcc of "Red Hat", so match centos first
	case strings.Contains(versionInfo, "centos"):
		return Centos
	case strings.Contains(versionInfo, "red hat"):
		return Redhat
	case strings.Contains(versionInfo, "ubuntu"):
		return Ubuntu
	}
	return Butt
}
