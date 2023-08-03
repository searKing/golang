// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package image

import (
	"image"
	"image/color"
	"math"
	"reflect"

	"github.com/searKing/golang/go/exp/constraints"
)

type Rectangle[E constraints.Number] struct {
	Min, Max Point[E]
}

// String returns a string representation of r like "(3,4)-(6,5)".
func (r Rectangle[E]) String() string {
	return r.Min.String() + "-" + r.Max.String()
}

// Dx returns r's width.
func (r Rectangle[E]) Dx() E {
	return r.Max.X - r.Min.X
}

// Dy returns r's height.
func (r Rectangle[E]) Dy() E {
	return r.Max.Y - r.Min.Y
}

// Size returns r's width and height.
func (r Rectangle[E]) Size() Point[E] {
	return Point[E]{
		r.Max.X - r.Min.X,
		r.Max.Y - r.Min.Y,
	}
}

// Area returns r's area, as width * height.
func (r Rectangle[E]) Area() E {
	dx, dy := r.Dx(), r.Dy()
	if dx <= 0 || dy <= 0 {
		return 0
	}
	return dx * dy
}

// Add returns the rectangle r translated by p.
func (r Rectangle[E]) Add(p Point[E]) Rectangle[E] {
	return Rectangle[E]{
		Point[E]{r.Min.X + p.X, r.Min.Y + p.Y},
		Point[E]{r.Max.X + p.X, r.Max.Y + p.Y},
	}
}

// Sub returns the rectangle r translated by -p.
func (r Rectangle[E]) Sub(p Point[E]) Rectangle[E] {
	return Rectangle[E]{
		Point[E]{r.Min.X - p.X, r.Min.Y - p.Y},
		Point[E]{r.Max.X - p.X, r.Max.Y - p.Y},
	}
}

// Mul returns the rectangle r translated by r*k.
func (r Rectangle[E]) Mul(k E) Rectangle[E] {
	return Rectangle[E]{r.Min.Mul(k), r.Max.Mul(k)}
}

// MulPoint returns the rectangle r translated by r.*p.
func (r Rectangle[E]) MulPoint(p Point[E]) Rectangle[E] {
	return Rectangle[E]{r.Min.MulPoint(p), r.Max.MulPoint(p)}
}

// MulRectangle returns the rectangle r translated by r.*p.
func (r Rectangle[E]) MulRectangle(p Rectangle[E]) Rectangle[E] {
	return Rectangle[E]{r.Min.MulPoint(p.Min), r.Max.MulPoint(p.Max)}
}

// Div returns the rectangle r translated by r/k.
func (r Rectangle[E]) Div(k E) Rectangle[E] {
	return Rectangle[E]{r.Min.Div(k), r.Max.Div(k)}
}

// DivPoint returns the rectangle r translated by r./p.
func (r Rectangle[E]) DivPoint(p Point[E]) Rectangle[E] {
	return Rectangle[E]{r.Min.DivPoint(p), r.Max.DivPoint(p)}
}

// DivRectangle returns the rectangle r translated by r./p.
func (r Rectangle[E]) DivRectangle(p Rectangle[E]) Rectangle[E] {
	return Rectangle[E]{r.Min.DivPoint(p.Min), r.Max.DivPoint(p.Max)}
}

// Inset returns the rectangle r inset by n, which may be negative. If either
// of r's dimensions is less than 2*n then an empty rectangle near the center
// of r will be returned.
func (r Rectangle[E]) Inset(n E) Rectangle[E] {
	return r.InsetPoint(Point[E]{X: n, Y: n})
}

// InsetPoint returns the rectangle r inset by n, which may be negative. If either
// of r's dimensions is less than n.X+n.Y then an empty rectangle near the center
// of r will be returned.
func (r Rectangle[E]) InsetPoint(n Point[E]) Rectangle[E] {
	return r.InsetRectangle(Rectangle[E]{Min: n, Max: n})
}

// InsetRectangle returns the rectangle r inset by n, which may be negative. If either
// of r's dimensions is less than (n.Min.X+n.Max.X, n.Min.Y+n.Max.Y), then an empty rectangle near the center
// of r will be returned.
func (r Rectangle[E]) InsetRectangle(n Rectangle[E]) Rectangle[E] {
	if r.Dx() < n.Min.X+n.Max.X {
		r.Min.X = (r.Min.X + r.Max.X) / 2
		r.Max.X = r.Min.X
	} else {
		r.Min.X += n.Min.X
		r.Max.X -= n.Max.X
	}
	if r.Dy() < n.Min.Y+n.Max.Y {
		r.Min.Y = (r.Min.Y + r.Max.Y) / 2
		r.Max.Y = r.Min.Y
	} else {
		r.Min.Y += n.Min.Y
		r.Max.Y -= n.Max.Y
	}
	return r
}

// Intersect returns the largest rectangle contained by both r and s. If the
// two rectangles do not overlap then the zero rectangle will be returned.
func (r Rectangle[E]) Intersect(s Rectangle[E]) Rectangle[E] {
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
		var zero Rectangle[E]
		return zero
	}
	return r
}

// Border returns four rectangles that together contain those points between r
// and r.Inset(inset). Visually:
//
//	00000000
//	00000000
//	11....22
//	11....22
//	11....22
//	33333333
//	33333333
//
// The inset may be negative, in which case the points will be outside r.
//
// Some of the returned rectangles may be empty. None of the returned
// rectangles will overlap.
func (r Rectangle[E]) Border(inset E) [4]Rectangle[E] {
	return r.BorderPoint(Point[E]{
		X: inset,
		Y: inset,
	})
}

// BorderPoint returns four rectangles that together contain those points between r
// and r.Inset(inset). Visually:
//
//	00000000
//	00000000
//	11....22
//	11....22
//	11....22
//	33333333
//	33333333
//
// The inset may be negative, in which case the points will be outside r.
//
// Some of the returned rectangles may be empty. None of the returned
// rectangles will overlap.
func (r Rectangle[E]) BorderPoint(inset Point[E]) [4]Rectangle[E] {
	return r.BorderRectangle(Rectangle[E]{
		Min: inset,
		Max: inset,
	})
}

// BorderRectangle returns four rectangles that together contain those points between r
// and r.Inset(inset). Visually:
//
//	00000000
//	00000000
//	11....22
//	11....22
//	11....22
//	33333333
//	33333333
//
// The inset may be negative, in which case the points will be outside r.
//
// Some of the returned rectangles may be empty. None of the returned
// rectangles will overlap.
func (r Rectangle[E]) BorderRectangle(inset Rectangle[E]) [4]Rectangle[E] {
	if inset.Min.X == 0 && inset.Min.Y == 0 && inset.Max.X == 0 && inset.Max.Y == 0 {
		return [4]Rectangle[E]{}
	}
	if r.Dx() <= inset.Min.X+inset.Max.X || r.Dy() <= inset.Min.Y+inset.Max.Y {
		return [4]Rectangle[E]{r}
	}

	x := [4]E{
		r.Min.X,
		r.Min.X + inset.Min.X,
		r.Max.X - inset.Max.X,
		r.Max.X,
	}
	y := [4]E{
		r.Min.Y,
		r.Min.Y + inset.Min.Y,
		r.Max.Y - inset.Max.Y,
		r.Max.Y,
	}
	if inset.Min.X < 0 {
		x[0], x[1] = x[1], x[0]
	}
	if inset.Max.X < 0 {
		x[2], x[3] = x[3], x[2]
	}
	if inset.Min.Y < 0 {
		y[0], y[1] = y[1], y[0]
	}
	if inset.Max.Y < 0 {
		y[2], y[3] = y[3], y[2]
	}

	// The top and bottom sections are responsible for filling the corners.
	// The top and bottom sections go from x[0] to x[3], across the y's.
	// The left and right sections go from y[1] to y[2], across the x's.

	return [4]Rectangle[E]{{
		// Top section.
		Min: Point[E]{
			X: x[0],
			Y: y[0],
		},
		Max: Point[E]{
			X: x[3],
			Y: y[1],
		},
	}, {
		// Left section.
		Min: Point[E]{
			X: x[0],
			Y: y[1],
		},
		Max: Point[E]{
			X: x[1],
			Y: y[2],
		},
	}, {
		// Right section.
		Min: Point[E]{
			X: x[2],
			Y: y[1],
		},
		Max: Point[E]{
			X: x[3],
			Y: y[2],
		},
	}, {
		// Bottom section.
		Min: Point[E]{
			X: x[0],
			Y: y[2],
		},
		Max: Point[E]{
			X: x[3],
			Y: y[3],
		},
	}}
}

// Union returns the smallest rectangle that contains both r and s.
func (r Rectangle[E]) Union(s Rectangle[E]) Rectangle[E] {
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
func (r Rectangle[E]) Empty() bool {
	return r.Min.X >= r.Max.X || r.Min.Y >= r.Max.Y
}

// Eq reports whether r and s contain the same set of points. All empty
// rectangles are considered equal.
func (r Rectangle[E]) Eq(s Rectangle[E]) bool {
	return r == s || r.Empty() && s.Empty()
}

// Overlaps reports whether r and s have a non-empty intersection.
func (r Rectangle[E]) Overlaps(s Rectangle[E]) bool {
	return !r.Empty() && !s.Empty() &&
		r.Min.X < s.Max.X && s.Min.X < r.Max.X &&
		r.Min.Y < s.Max.Y && s.Min.Y < r.Max.Y
}

// In reports whether every point in r is in s.
func (r Rectangle[E]) In(s Rectangle[E]) bool {
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
func (r Rectangle[E]) Canon() Rectangle[E] {
	if r.Max.X < r.Min.X {
		r.Min.X, r.Max.X = r.Max.X, r.Min.X
	}
	if r.Max.Y < r.Min.Y {
		r.Min.Y, r.Max.Y = r.Max.Y, r.Min.Y
	}
	return r
}

// At implements the Image interface.
func (r Rectangle[E]) At(x, y E) color.Color {
	if (Point[E]{x, y}).In(r) {
		return color.Opaque
	}
	return color.Transparent
}

// RGBA64At implements the RGBA64Image interface.
func (r Rectangle[E]) RGBA64At(x, y E) color.RGBA64 {
	if (Point[E]{x, y}).In(r) {
		return color.RGBA64{R: 0xffff, G: 0xffff, B: 0xffff, A: 0xffff}
	}
	return color.RGBA64{}
}

// Bounds implements the Image interface.
func (r Rectangle[E]) Bounds() Rectangle[E] {
	return r
}

// ColorModel implements the Image interface.
func (r Rectangle[E]) ColorModel() color.Model {
	return color.Alpha16Model
}

func (r Rectangle[E]) RoundRectangle() image.Rectangle {
	return image.Rect(round(r.Min.X), round(r.Min.Y), round(r.Max.X), round(r.Max.Y))
}

// UnionPoints returns the smallest rectangle that contains all points.
func (r Rectangle[E]) UnionPoints(pts ...Point[E]) Rectangle[E] {
	if len(pts) == 0 {
		return r
	}
	var pos int
	if r.Empty() { // an empty rectangle is an empty set, Not a point
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

// Scale scale rectangle's size to size, by expand from the mid-point of rect.
func (r Rectangle[E]) Scale(size Point[E]) Rectangle[E] {
	size = size.Sub(r.Size())
	minOffset := size.Div(2)
	maxOffset := size.Sub(minOffset)

	return Rectangle[E]{
		Min: Point[E]{X: r.Min.X - minOffset.X, Y: r.Min.Y - minOffset.Y},
		Max: Point[E]{X: r.Max.X + maxOffset.X, Y: r.Max.Y + maxOffset.Y},
	}
}

// ScaleByFactor scale rectangle's size to factor * size, by expand from the mid-point of rect.
func (r Rectangle[E]) ScaleByFactor(factor Point[E]) Rectangle[E] {
	if r.Empty() {
		return r
	}
	return r.Scale(r.Size().MulPoint(factor))
}

// FlexIn flex rect into box, shrink but not grow to fit the space available in its flex container
func (r Rectangle[E]) FlexIn(container Rectangle[E]) Rectangle[E] {
	r2 := r
	r2.Min = container.Min
	r2.Max = container.Min.Add(r.Size())
	if r.Min.X > r2.Min.X {
		r2.Max.X += r.Min.X - r2.Min.X
		r2.Min.X = r.Min.X
	}
	if r.Min.Y > r2.Min.Y {
		r2.Max.Y += r.Min.Y - r2.Min.Y
		r2.Min.Y = r.Min.Y
	}
	if r2.Max.X > container.Max.X {
		r2.Min.X -= r2.Max.X - container.Max.X
		r2.Max.X = container.Max.X
	}
	if r2.Max.Y > container.Max.Y {
		r2.Min.Y -= r2.Max.Y - container.Max.Y
		r2.Max.Y = container.Max.Y
	}
	if r2.Min.X < container.Min.X {
		r2.Min.X = container.Min.X
	}
	if r2.Min.Y < container.Min.Y {
		r2.Min.Y = container.Min.Y
	}
	return r2
}

// Rect is shorthand for Rectangle[E]{Pt(x0, y0), Pt(x1, y1)}. The returned
// rectangle has minimum and maximum coordinates swapped if necessary so that
// it is well-formed.
func Rect[E constraints.Number](x0, y0, x1, y1 E) Rectangle[E] {
	if x0 > x1 {
		x0, x1 = x1, x0
	}
	if y0 > y1 {
		y0, y1 = y1, y0
	}
	return Rectangle[E]{Point[E]{x0, y0}, Point[E]{x1, y1}}
}

func FromRectInt[E constraints.Number](rect image.Rectangle) Rectangle[E] {
	return Rect(E(rect.Min.X), E(rect.Min.Y), E(rect.Max.X), E(rect.Max.Y))
}

func round[E constraints.Number](x E) int {
	kind := reflect.TypeOf(x).Kind()
	switch kind {
	case reflect.Float32, reflect.Float64:
		return int(math.Round(float64(x)))
	}
	return int(x)
}
