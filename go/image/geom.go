// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package image

import "image"

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
