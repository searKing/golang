// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package image

import (
	"fmt"
	"image"

	"github.com/searKing/golang/go/exp/constraints"
)

// A Point is an X, Y coordinate pair. The axes increase right and down.
type Point[E constraints.Number] struct {
	X, Y E
}

// String returns a string representation of p like "(3,4)".
func (p Point[E]) String() string {
	return fmt.Sprintf("(%v,%v)", p.X, p.Y)
}

// Add returns the vector p+q.
func (p Point[E]) Add(q Point[E]) Point[E] {
	return Point[E]{p.X + q.X, p.Y + q.Y}
}

// Sub returns the vector p-q.
func (p Point[E]) Sub(q Point[E]) Point[E] {
	return Point[E]{p.X - q.X, p.Y - q.Y}
}

// Mul returns the vector p*k.
func (p Point[E]) Mul(k E) Point[E] {
	return Point[E]{p.X * k, p.Y * k}
}

// MulPoint returns the vector p.*k.
func (p Point[E]) MulPoint(k Point[E]) Point[E] {
	return Point[E]{p.X * k.X, p.Y * k.Y}
}

// Div returns the vector p/k.
func (p Point[E]) Div(k E) Point[E] {
	return Point[E]{p.X / k, p.Y / k}
}

// DivPoint returns the vector p./k.
func (p Point[E]) DivPoint(k Point[E]) Point[E] {
	return Point[E]{p.X / k.X, p.Y / k.Y}
}

// In reports whether p is in r.
func (p Point[E]) In(r Rectangle[E]) bool {
	return r.Min.X <= p.X && p.X < r.Max.X &&
		r.Min.Y <= p.Y && p.Y < r.Max.Y
}

// Mod returns the point q in r such that p.X-q.X is a multiple of r's width
// and p.Y-q.Y is a multiple of r's height.
func (p Point[E]) Mod(r image.Rectangle) image.Point {
	p2 := p.RoundPoint()
	return p2.Mod(r)
}

// Eq reports whether p and q are equal.
func (p Point[E]) Eq(q Point[E]) bool {
	return p == q
}

func (p Point[E]) RoundPoint() image.Point {
	return image.Pt(round(p.X), round(p.Y))
}

// Pt is shorthand for Point[E]{X, Y}.
func Pt[E constraints.Number](X, Y E) Point[E] {
	return Point[E]{X, Y}
}

func FromPtInt[E constraints.Number](q image.Point) Point[E] {
	return Point[E]{E(q.X), E(q.Y)}
}
