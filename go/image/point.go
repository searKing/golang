// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package image

import (
	"fmt"
	"image"
)

// A Point2f is an X, Y coordinate pair. The axes increase right and down.
type Point2f struct {
	X, Y float32
}

// String returns a string representation of p like "(3,4)".
func (p Point2f) String() string {
	return fmt.Sprintf("(%.2f,%.2f)", p.X, p.Y)
}

// Add returns the vector p+q.
func (p Point2f) Add(q Point2f) Point2f {
	return Point2f{p.X + q.X, p.Y + q.Y}
}

// Sub returns the vector p-q.
func (p Point2f) Sub(q Point2f) Point2f {
	return Point2f{p.X - q.X, p.Y - q.Y}
}

// Mul returns the vector p*k.
func (p Point2f) Mul(k float32) Point2f {
	return Point2f{p.X * k, p.Y * k}
}

// Div returns the vector p/k.
func (p Point2f) Div(k float32) Point2f {
	return Point2f{p.X / k, p.Y / k}
}

// In reports whether p is in r.
func (p Point2f) In(r Rectangle2f) bool {
	return r.Min.X <= p.X && p.X < r.Max.X &&
		r.Min.Y <= p.Y && p.Y < r.Max.Y
}

// Mod returns the point q in r such that p.X-q.X is a multiple of r's width
// and p.Y-q.Y is a multiple of r's height.
func (p Point2f) Mod(r image.Rectangle) image.Point {
	p2 := p.RoundPoint()
	return p2.Mod(r)
}

// Eq reports whether p and q are equal.
func (p Point2f) Eq(q Point2f) bool {
	return p == q
}

func (p Point2f) RoundPoint() image.Point {
	return image.Pt(round(p.X), round(p.Y))
}

// ZP2f is the zero Point2f.
//
// Deprecated: Use a literal image.Point2f{} instead.
var ZP2f Point2f

// Pt2f is shorthand for Point2f{X, Y}.
func Pt2f(X, Y float32) Point2f {
	return Point2f{X, Y}
}
