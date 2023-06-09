// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package image

import (
	"math"

	constraints_ "github.com/searKing/golang/go/exp/constraints"
	"golang.org/x/exp/constraints"
)

// UnionPoints returns the smallest rectangle that contains all points.
// an empty rectangle is a empty set, Not a point
func UnionPoints[E constraints_.Number](pts ...Point[E]) Rectangle[E] {
	if len(pts) == 0 {
		return Rectangle[E]{}
	}

	r := Rectangle[E]{
		Min: pts[0],
		Max: pts[0],
	}
	for _, p := range pts[1:] {
		if p.X < r.Min.X {
			r.Min.X = p.X
		}
		if p.Y < r.Min.Y {
			r.Min.Y = p.Y
		}
		if p.X > r.Max.X {
			r.Max.X = p.X
		}
		if p.Y > r.Max.Y {
			r.Max.Y = p.Y
		}
	}
	return r
}

// UnionRectangles returns the smallest rectangle that contains all rectangles, empty rectangles excluded.
func UnionRectangles[E constraints_.Number](rs ...Rectangle[E]) Rectangle[E] {
	var ur Rectangle[E]
	for _, r := range rs {
		ur = ur.Union(r)
	}
	return ur
}

// ScaleLineSegment segment's size to length flexible in limit
func ScaleLineSegment[E constraints_.Number](segment Point[E], length E, limit Point[E]) Point[E] {
	var swapped = segment.X > segment.Y
	if swapped { // swap (X,Y) -> (Y,X)
		segment.X, segment.Y = segment.Y, segment.X
		limit.X, limit.Y = limit.Y, limit.X
	}

	dx := length - (segment.Y - segment.X)
	segment.X -= E(math.Round(float64(dx) / 2.0))
	if segment.X < limit.X {
		segment.X = limit.X
	}
	segment.Y = segment.X + length
	if segment.Y > limit.Y {
		segment.Y = limit.Y
		segment.X = segment.Y - length
		if segment.X < limit.X {
			segment.X = limit.X
		}
	}

	if swapped {
		segment.X, segment.Y = segment.Y, segment.X
	}
	return segment
}

// ScaleRectangleBySize scale rect to size flexible in limit
func ScaleRectangleBySize[E constraints.Integer](rect Rectangle[E], size Point[E], limit Rectangle[E]) Rectangle[E] {
	// padding in x direction
	x := ScaleLineSegment(Pt(rect.Min.X, rect.Max.X), size.X, Pt(limit.Min.X, limit.Max.X))
	// padding in y direction
	y := ScaleLineSegment(Pt(rect.Min.Y, rect.Max.Y), size.Y, Pt(limit.Min.Y, limit.Max.Y))

	return limit.Intersect(Rect(x.X, y.X, x.Y, y.Y))
}
