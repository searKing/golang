// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package distribution

import (
	"os/exec"
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
	}
	return "Butt"
}

func GetOSVersion() Distribution {
	if runtime.GOOS == "Windows" {
		return Windows
	}

	params := []string{
		"cat",
		"/proc/version",
	}
	data, err := exec.Command(params[0], params[1:]...).CombinedOutput()
	if err != nil {
		return Butt
	}

	versionInfo := strings.ToLower(string(data))

	switch {
	case strings.Contains(versionInfo, "red hat"):
		return Redhat
	case strings.Contains(versionInfo, "centos"):
		return Centos
	case strings.Contains(versionInfo, "ubuntu"):
		return Ubuntu
	}
	return Butt
}
