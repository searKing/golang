// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package image

import (
	"image"
	"math"
)

// UnionPoints returns the smallest rectangle that contains all points.
func UnionPoints(pts ...image.Point) image.Rectangle {
	if len(pts) == 0 {
		return image.Rectangle{}
	}

	r := image.Rectangle{
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

// UnionRectangles returns the smallest rectangle that contains all rectangles.
func UnionRectangles(rs ...image.Rectangle) image.Rectangle {
	var ur image.Rectangle
	for _, r := range rs {
		ur = ur.Union(r)
	}
	return ur
}

func scale(segment image.Point, length int, limit image.Point) image.Point {
	var swapped = segment.X > segment.Y
	if swapped { // swap (X,Y) -> (Y,X)
		segment.X = segment.X ^ segment.Y
		segment.Y = segment.X ^ segment.Y
		segment.X = segment.X ^ segment.Y
		limit.X = limit.X ^ limit.Y
		limit.Y = limit.X ^ limit.Y
		limit.X = limit.X ^ limit.Y
	}

	dx := length - (segment.Y - segment.X)
	segment.X -= int(math.Round(float64(dx) / 2.0))
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
		segment.X = segment.X ^ segment.Y
		segment.Y = segment.X ^ segment.Y
		segment.X = segment.X ^ segment.Y
	}
	return segment
}

// ScaleRectangleBySize scale rect to size flexible in limit
func ScaleRectangleBySize(rect image.Rectangle, size image.Point, limit image.Rectangle) image.Rectangle {
	// padding in x direction
	x := scale(image.Pt(rect.Min.X, rect.Max.X), size.X, image.Pt(limit.Min.X, limit.Max.X))
	// padding in x direction
	y := scale(image.Pt(rect.Min.Y, rect.Max.Y), size.X, image.Pt(limit.Min.Y, limit.Max.Y))

	return limit.Intersect(image.Rect(x.X, y.X, x.Y, y.Y))
}
