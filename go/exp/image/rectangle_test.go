// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package image_test

import (
	"fmt"
	"strconv"
	"testing"

	image_ "github.com/searKing/golang/go/exp/image"
)

func TestRectangle(t *testing.T) {
	// in checks that every point in f is in g.
	in := func(f, g image_.Rectangle[float32]) error {
		if !f.In(g) {
			return fmt.Errorf("f=%s, f.In(%s): got false, want true", f, g)
		}
		for y := f.Min.Y; y < f.Max.Y; y++ {
			for x := f.Min.X; x < f.Max.X; x++ {
				p := image_.Pt[float32](x, y)
				if !p.In(g) {
					return fmt.Errorf("p=%s, p.In(%s): got false, want true", p, g)
				}
			}
		}
		return nil
	}

	rects := []image_.Rectangle[float32]{
		image_.Rect[float32](0, 0, 10, 10),
		image_.Rect[float32](10, 0, 20, 10),
		image_.Rect[float32](1, 2, 3, 4),
		image_.Rect[float32](4, 6, 10, 10),
		image_.Rect[float32](2, 3, 12, 5),
		image_.Rect[float32](-1, -2, 0, 0),
		image_.Rect[float32](-1, -2, 4, 6),
		image_.Rect[float32](-10, -20, 30, 40),
		image_.Rect[float32](8, 8, 8, 8),
		image_.Rect[float32](88, 88, 88, 88),
		image_.Rect[float32](6, 5, 4, 3),
	}

	// r.Eq(s) should be equivalent to every point in r being in s, and every
	// point in s being in r.
	for _, r := range rects {
		for _, s := range rects {
			got := r.Eq(s)
			want := in(r, s) == nil && in(s, r) == nil
			if got != want {
				t.Errorf("Eq: r=%s, s=%s: got %t, want %t", r, s, got, want)
			}
		}
	}

	// The intersection should be the largest rectangle a such that every point
	// in a is both in r and in s.
	for _, r := range rects {
		for _, s := range rects {
			a := r.Intersect(s)
			if err := in(a, r); err != nil {
				t.Errorf("Intersect: r=%s, s=%s, a=%s, a not in r: %v", r, s, a, err)
			}
			if err := in(a, s); err != nil {
				t.Errorf("Intersect: r=%s, s=%s, a=%s, a not in s: %v", r, s, a, err)
			}
			if isZero, overlaps := a == (image_.Rectangle[float32]{}), r.Overlaps(s); isZero == overlaps {
				t.Errorf("Intersect: r=%s, s=%s, a=%s: isZero=%t same as overlaps=%t",
					r, s, a, isZero, overlaps)
			}
			largerThanA := [4]image_.Rectangle[float32]{a, a, a, a}
			largerThanA[0].Min.X--
			largerThanA[1].Min.Y--
			largerThanA[2].Max.X++
			largerThanA[3].Max.Y++
			for i, b := range largerThanA {
				if b.Empty() {
					// b isn't actually larger than a.
					continue
				}
				if in(b, r) == nil && in(b, s) == nil {
					t.Errorf("Intersect: r=%s, s=%s, a=%s, b=%s, i=%d: intersection could be larger",
						r, s, a, b, i)
				}
			}
		}
	}

	// The union should be the smallest rectangle a such that every point in r
	// is in a and every point in s is in a.
	for _, r := range rects {
		for _, s := range rects {
			a := r.Union(s)
			if err := in(r, a); err != nil {
				t.Errorf("Union: r=%s, s=%s, a=%s, r not in a: %v", r, s, a, err)
			}
			if err := in(s, a); err != nil {
				t.Errorf("Union: r=%s, s=%s, a=%s, s not in a: %v", r, s, a, err)
			}
			if a.Empty() {
				// You can't get any smaller than a.
				continue
			}
			smallerThanA := [4]image_.Rectangle[float32]{a, a, a, a}
			smallerThanA[0].Min.X++
			smallerThanA[1].Min.Y++
			smallerThanA[2].Max.X--
			smallerThanA[3].Max.Y--
			for i, b := range smallerThanA {
				if in(r, b) == nil && in(s, b) == nil {
					t.Errorf("Union: r=%s, s=%s, a=%s, b=%s, i=%d: union could be smaller",
						r, s, a, b, i)
				}
			}
		}
	}
}

func TestRectangle_Border(t *testing.T) {
	r := image_.Rect(100, 200, 400, 300)

	insets := []int{
		-100,
		-1,
		+0,
		+1,
		+20,
		+49,
		+50,
		+51,
		+149,
		+150,
		+151,
	}

	for _, inset := range insets {
		border := r.Border(inset)

		outer, inner := r, r.Inset(inset)
		if inset < 0 {
			outer, inner = inner, outer
		}

		got := 0
		for _, b := range border {
			got += b.Area()
		}
		want := outer.Area() - inner.Area()
		if got != want {
			t.Errorf("inset=%d: total area: got %d, want %d", inset, got, want)
		}

		for i, bi := range border {
			for j, bj := range border {
				if i <= j {
					continue
				}
				if !bi.Intersect(bj).Empty() {
					t.Errorf("inset=%d: %v and %v overlap", inset, bi, bj)
				}
			}
		}

		for _, b := range border {
			if got := outer.Intersect(b); got != b {
				t.Errorf("inset=%d: outer intersection: got %v, want %v", inset, got, b)
			}
			if got := inner.Intersect(b); !got.Empty() {
				t.Errorf("inset=%d: inner intersection: got %v, want empty", inset, got)
			}
		}
	}
}

func TestRectangle_BorderPoint(t *testing.T) {
	r := image_.Rect(100, 200, 400, 300)

	insets := []image_.Point[int]{
		image_.Pt(-100, -100),
		image_.Pt(-1, -1),
		image_.Pt(+0, +0),
		image_.Pt(+1, +1),
		image_.Pt(+20, +20),
		image_.Pt(+49, +49),
		image_.Pt(+50, +50),
		image_.Pt(+51, +51),
		image_.Pt(+149, +149),
		image_.Pt(+150, +150),
		image_.Pt(+151, +151),
	}

	for _, inset := range insets {
		border := r.BorderPoint(inset)

		outer, inner := r, r.InsetPoint(inset)
		if outer.Area() < inner.Area() {
			outer, inner = inner, outer
		}

		got := 0
		for _, b := range border {
			got += b.Area()
		}
		want := outer.Area() - inner.Area()
		if got != want {
			t.Errorf("inset=%d: total area: got %d, want %d", inset, got, want)
		}

		for i, bi := range border {
			for j, bj := range border {
				if i <= j {
					continue
				}
				if !bi.Intersect(bj).Empty() {
					t.Errorf("inset=%d: %v and %v overlap", inset, bi, bj)
				}
			}
		}

		for _, b := range border {
			if got := outer.Intersect(b); got != b {
				t.Errorf("inset=%d: outer intersection: got %v, want %v", inset, got, b)
			}
			if got := inner.Intersect(b); !got.Empty() {
				t.Errorf("inset=%d: inner intersection: got %v, want empty", inset, got)
			}
		}
	}
}

func TestRectangle_BorderRectangle(t *testing.T) {
	r := image_.Rect(100, 200, 400, 300)

	insets := []image_.Rectangle[int]{
		image_.Rect(-100, -100, -100, -100),
		image_.Rect(-1, -1, -1, -1),
		image_.Rect(+0, +0, +0, +0),
		image_.Rect(+1, +1, +1, +1),
		image_.Rect(+20, +20, +20, +20),
		image_.Rect(+49, +49, +49, +49),
		image_.Rect(+50, +50, +50, +50),
		image_.Rect(+51, +51, +51, +51),
		image_.Rect(+149, +149, +149, +149),
		image_.Rect(+150, +150, +150, +150),
		image_.Rect(+151, +151, +151, +151),
	}

	for _, inset := range insets {
		border := r.BorderRectangle(inset)

		outer, inner := r, r.InsetRectangle(inset)
		if outer.Area() < inner.Area() {
			outer, inner = inner, outer
		}

		got := 0
		for _, b := range border {
			got += b.Area()
		}
		want := outer.Area() - inner.Area()
		if got != want {
			t.Errorf("inset=%d: total area: got %d, want %d", inset, got, want)
		}

		for i, bi := range border {
			for j, bj := range border {
				if i <= j {
					continue
				}
				if !bi.Intersect(bj).Empty() {
					t.Errorf("inset=%d: %v and %v overlap", inset, bi, bj)
				}
			}
		}

		for _, b := range border {
			if got := outer.Intersect(b); got != b {
				t.Errorf("inset=%d: outer intersection: got %v, want %v", inset, got, b)
			}
			if got := inner.Intersect(b); !got.Empty() {
				t.Errorf("inset=%d: inner intersection: got %v, want empty", inset, got)
			}
		}
	}
}

func TestRectangle_FlexIn(t *testing.T) {
	tests := []struct {
		r    image_.Rectangle[int]
		box  image_.Rectangle[int]
		want image_.Rectangle[int]
	}{
		{image_.Rectangle[int]{}, image_.Rectangle[int]{}, image_.Rectangle[int]{}},
		{image_.Rect(0, 0, 10, 10), image_.Rectangle[int]{}, image_.Rectangle[int]{}},
		{image_.Rect(0, 0, 10, 10), image_.Rect(1, 2, 3, 4), image_.Rect(1, 2, 3, 4)},
		{image_.Rect(0, 0, 10, 10), image_.Rect(-1, -2, 13, 14), image_.Rect(0, 0, 10, 10)},
		{image_.Rect(0, 0, 10, 10), image_.Rect(1, -2, 3, 14), image_.Rect(1, 0, 3, 10)},
		{image_.Rect(0, 0, 10, 10), image_.Rect(1, 2, 13, 14), image_.Rect(1, 2, 11, 12)},
		{image_.Rect(0, 0, 10, 10), image_.Rect(1, 2, 13, 10), image_.Rect(1, 2, 11, 10)},
		{image_.Rect(0, 0, 10, 10), image_.Rect(20, 20, 20, 20), image_.Rect(20, 20, 20, 20)},
		{image_.Rect(0, 0, 10, 10), image_.Rect(1, 20, 3, 20), image_.Rect(1, 20, 3, 20)},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := tt.r.FlexIn(tt.box)
			if got != tt.want {
				t.Errorf("(%v).FlexIn(%v) got (%v), want (%v)", tt.r, tt.box, got, tt.want)
			}
		})
	}
}
