// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package reflect

import (
	"testing"
	"time"
)

func TestIsNil(t *testing.T) {
	var nilTime *time.Time
	tests := []struct {
		a    any
		want bool
	}{
		{
			a:    nil,
			want: true,
		},
		{
			a:    true,
			want: false,
		},
		{
			a:    0,
			want: false,
		},
		{
			a:    "",
			want: false,
		},
		{
			a:    time.Now(),
			want: false,
		},
		{
			a:    nilTime,
			want: true,
		},
	}
	for i, in := range tests {
		if got := IsNil(in.a); got != in.want {
			t.Errorf("case %d: IsNil(%v) = %t, want %t， type is %T", i, in.a, got, in.want, in.a)
		}
	}
}
