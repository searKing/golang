// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package expvar_test

import (
	"expvar"
	"testing"

	expvar_ "github.com/searKing/golang/go/expvar"
)

func TestLeak(t *testing.T) {
	reqs := expvar_.NewLeak("requests")
	if i := reqs.Value(); i != 0 {
		t.Errorf("reqs.Value() = %v, want 0", i)
	}
	if reqs != expvar.Get("requests").(*expvar_.Leak) {
		t.Errorf("Get() failed.")
	}

	reqs.Add(1)
	reqs.Add(3)
	if i := reqs.Value(); i != 4 {
		t.Errorf("reqs.Value() = %v, want 4", i)
	}

	if s := reqs.String(); s != "[4 4]" {
		t.Errorf("reqs.String() = %q, want %q", "[4 4]", s)
	}

	reqs.Add(-4)
	if i := reqs.Value(); i != 0 {
		t.Errorf("reqs.Value() = %v, want 4", i)
	}

	if s := reqs.String(); s != "[0 4]" {
		t.Errorf("reqs.String() = %q, want %q", "[0 4]", s)
	}
	reqs.Add(1)
	reqs.Done()
	if i := reqs.Value(); i != 0 {
		t.Errorf("reqs.Value() = %v, want 4", i)
	}

	if s := reqs.String(); s != "[0 5]" {
		t.Errorf("reqs.String() = %q, want %q", "[0 5]", s)
	}
}
