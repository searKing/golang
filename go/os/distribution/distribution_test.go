// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package distribution

import "testing"

func TestParseOSRelease(t *testing.T) {
	tests := []struct {
		name string
		data string
		want Distribution
	}{
		{name: "ubuntu", data: "NAME=\"Ubuntu\"\nVERSION_ID=\"22.04\"\nID=ubuntu\nID_LIKE=debian\n", want: Ubuntu},
		{name: "centos", data: "NAME=\"CentOS Linux\"\nID=\"centos\"\nID_LIKE=\"rhel fedora\"\n", want: Centos},
		{name: "rhel", data: "NAME=\"Red Hat Enterprise Linux\"\nID=\"rhel\"\nID_LIKE=\"fedora\"\n", want: Redhat},
		{name: "unknown", data: "NAME=\"Debian GNU/Linux\"\nID=debian\n", want: Butt},
		{name: "no id", data: "NAME=\"Linux\"\n", want: Butt},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseOSRelease(tt.data); got != tt.want {
				t.Errorf("parseOSRelease() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseProcVersion(t *testing.T) {
	tests := []struct {
		name string
		data string
		want Distribution
	}{
		{
			name: "centos",
			data: "Linux version 3.10.0-1160.el7.x86_64 (mockbuild@kbuilder.bsys.centos.org) (gcc version 4.8.5 20150623 (Red Hat 4.8.5-44) (GCC) ) #1 SMP",
			want: Centos,
		},
		{
			name: "rhel",
			data: "Linux version 4.18.0-372.9.1.el8.x86_64 (mockbuild@x86-vm-07.build.eng.bos.redhat.com) (gcc version 8.5.0 20210514 (Red Hat 8.5.0-10) (GCC)) #1 SMP",
			want: Redhat,
		},
		{
			name: "ubuntu",
			data: "Linux version 5.15.0-91-generic (buildd@lcy02-amd64-045) (gcc (Ubuntu 11.4.0-1ubuntu1~22.04) 11.4.0) #101-Ubuntu SMP",
			want: Ubuntu,
		},
		{name: "unknown", data: "Linux version 6.1.0 (gcc 12.2.0)", want: Butt},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseProcVersion(tt.data); got != tt.want {
				t.Errorf("parseProcVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}
