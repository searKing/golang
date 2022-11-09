// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package image

import (
	"image"
	"image/color"
	"math"
)

type Rectangle2f struct {
	Min, Max Point2f
}

// String returns a string representation of r like "(3,4)-(6,5)".
func (r Rectangle2f) String() string {
	return r.Min.String() + "-" + r.Max.String()
}

// Dx returns r's width.
func (r Rectangle2f) Dx() float32 {
	return r.Max.X - r.Min.X
}

// Dy returns r's height.
func (r Rectangle2f) Dy() float32 {
	return r.Max.Y - r.Min.Y
}

// Size returns r's width and height.
func (r Rectangle2f) Size() Point2f {
	return Point2f{
		r.Max.X - r.Min.X,
		r.Max.Y - r.Min.Y,
	}
}

// Add returns the rectangle r translated by p.
func (r Rectangle2f) Add(p Point2f) Rectangle2f {
	return Rectangle2f{
		Point2f{r.Min.X + p.X, r.Min.Y + p.Y},
		Point2f{r.Max.X + p.X, r.Max.Y + p.Y},
	}
}

// Sub returns the rectangle r translated by -p.
func (r Rectangle2f) Sub(p Point2f) Rectangle2f {
	return Rectangle2f{
		Point2f{r.Min.X - p.X, r.Min.Y - p.Y},
		Point2f{r.Max.X - p.X, r.Max.Y - p.Y},
	}
}

// Inset returns the rectangle r inset by n, which may be negative. If either
// of r's dimensions is less than 2*n then an empty rectangle near the center
// of r will be returned.
func (r Rectangle2f) Inset(n float32) Rectangle2f {
	if r.Dx() < 2*n {
		r.Min.X = (r.Min.X + r.Max.X) / 2
		r.Max.X = r.Min.X
	} else {
		r.Min.X += n
		r.Max.X -= n
	}
	if r.Dy() < 2*n {
		r.Min.Y = (r.Min.Y + r.Max.Y) / 2
		r.Max.Y = r.Min.Y
	} else {
		r.Min.Y += n
		r.Max.Y -= n
	}
	return r
}

// Intersect returns the largest rectangle contained by both r and s. If the
// two rectangles do not overlap then the zero rectangle will be returned.
func (r Rectangle2f) Intersect(s Rectangle2f) Rectangle2f {
	if r.Min.X < s.Min.X {
		r.Min.X = s.Min.X
	}
	if r.Min.Y < s.Min.Y {
		r.Min.Y = s.Min.Y
	}
	if r.Max.X > s.Max.X {
		r.Max.X = s.Max.X
	}
	if r.Max.Y > s.Max.Y {
		r.Max.Y = s.Max.Y
	}
	// Letting r0 and s0 be the values of r and s at the time that the method
	// is called, this next line is equivalent to:
	//
	// if max(r0.Min.X, s0.Min.X) >= min(r0.Max.X, s0.Max.X) || likewiseForY { etc }
	if r.Empty() {
		return ZR2f
	}
	return r
}

// Union returns the smallest rectangle that contains both r and s.
func (r Rectangle2f) Union(s Rectangle2f) Rectangle2f {
	if r.Empty() {
		return s
	}
	if s.Empty() {
		return r
	}
	if r.Min.X > s.Min.X {
		r.Min.X = s.Min.X
	}
	if r.Min.Y > s.Min.Y {
		r.Min.Y = s.Min.Y
	}
	if r.Max.X < s.Max.X {
		r.Max.X = s.Max.X
	}
	if r.Max.Y < s.Max.Y {
		r.Max.Y = s.Max.Y
	}
	return r
}

// Empty reports whether the rectangle contains no points.
func (r Rectangle2f) Empty() bool {
	return r.Min.X >= r.Max.X || r.Min.Y >= r.Max.Y
}

// Eq reports whether r and s contain the same set of points. All empty
// rectangles are considered equal.
func (r Rectangle2f) Eq(s Rectangle2f) bool {
	return r == s || r.Empty() && s.Empty()
}

// Overlaps reports whether r and s have a non-empty intersection.
func (r Rectangle2f) Overlaps(s Rectangle2f) bool {
	return !r.Empty() && !s.Empty() &&
		r.Min.X < s.Max.X && s.Min.X < r.Max.X &&
		r.Min.Y < s.Max.Y && s.Min.Y < r.Max.Y
}

// In reports whether every point in r is in s.
func (r Rectangle2f) In(s Rectangle2f) bool {
	if r.Empty() {
		return true
	}
	// Note that r.Max is an exclusive bound for r, so that r.In(s)
	// does not require that r.Max.In(s).
	return s.Min.X <= r.Min.X && r.Max.X <= s.Max.X &&
		s.Min.Y <= r.Min.Y && r.Max.Y <= s.Max.Y
}

// Canon returns the canonical version of r. The returned rectangle has minimum
// and maximum coordinates swapped if necessary so that it is well-formed.
func (r Rectangle2f) Canon() Rectangle2f {
	if r.Max.X < r.Min.X {
		r.Min.X, r.Max.X = r.Max.X, r.Min.X
	}
	if r.Max.Y < r.Min.Y {
		r.Min.Y, r.Max.Y = r.Max.Y, r.Min.Y
	}
	return r
}

// At implements the Image interface.
func (r Rectangle2f) At(x, y float32) color.Color {
	if (Point2f{x, y}).In(r) {
		return color.Opaque
	}
	return color.Transparent
}

// RGBA64At implements the RGBA64Image interface.
func (r Rectangle2f) RGBA64At(x, y float32) color.RGBA64 {
	if (Point2f{x, y}).In(r) {
		return color.RGBA64{R: 0xffff, G: 0xffff, B: 0xffff, A: 0xffff}
	}
	return color.RGBA64{}
}

// Bounds implements the Image interface.
func (r Rectangle2f) Bounds() Rectangle2f {
	return r
}

// ColorModel implements the Image interface.
func (r Rectangle2f) ColorModel() color.Model {
	return color.Alpha16Model
}

func (r Rectangle2f) RoundRectangle() image.Rectangle {
	return image.Rect(round(r.Min.X), round(r.Min.Y), round(r.Max.X), round(r.Max.Y))
}

// UnionPoints returns the smallest rectangle that contains all points.
func (r Rectangle2f) UnionPoints(pts ...Point2f) Rectangle2f {
	if len(pts) == 0 {
		return r
	}
	var pos int
	if r.Empty() { // an empty rectangle is a empty set, Not a point
		r.Min = pts[0]
		r.Max = pts[0]
		pos = 1
	}
	for _, p := range pts[pos:] {
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

// ScaleByFactor scale rect to factor*size
func (r Rectangle2f) ScaleByFactor(factor Point2f) Rectangle2f {
	if r.Empty() {
		return r
	}
	factor = factor.Sub(Pt2f(1, 1))
	minOffset := Point2f{
		X: r.Dx() * factor.X / 2,
		Y: r.Dy() * factor.Y / 2,
	}
	maxOffset := Point2f{
		X: r.Dx() * factor.X,
		Y: r.Dy() * factor.Y,
	}.Sub(minOffset)

	return Rectangle2f{
		Min: Point2f{X: r.Min.X - minOffset.X, Y: r.Min.Y - minOffset.Y},
		Max: Point2f{X: r.Max.X + maxOffset.X, Y: r.Max.Y + maxOffset.Y},
	}
}

// ZR2f is the zero Rectangle2f.
//
// Deprecated: Use a literal image.Rectangle2f{} instead.
var ZR2f Rectangle2f

// Rect2f is shorthand for Rectangle2f{Pt(x0, y0), Pt(x1, y1)}. The returned
// rectangle has minimum and maximum coordinates swapped if necessary so that
// it is well-formed.
func Rect2f(x0, y0, x1, y1 float32) Rectangle2f {
	if x0 > x1 {
		x0, x1 = x1, x0
	}
	if y0 > y1 {
		y0, y1 = y1, y0
	}
	return Rectangle2f{Point2f{x0, y0}, Point2f{x1, y1}}
}

func Rect2fFromRect(rect image.Rectangle) Rectangle2f {
	return Rect2f(float32(rect.Min.X), float32(rect.Min.Y), float32(rect.Max.X), float32(rect.Max.Y))
}

func round(x float32) int {
	return int(math.Round(float64(x)))
}
