// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package image_test

import (
	"image"
	"testing"

	image_ "github.com/searKing/golang/go/image"
)

func ScaleRectangleBySize(t *testing.T) {
	tests := []struct {
		rect     image.Rectangle
		size     image.Point
		limit    image.Rectangle
		wantRect image.Rectangle
	}{
		{
			rect:     image.Rect(0, 0, 10, 10),
			size:     image.Pt(5, 5),
			limit:    image.Rect(0, 0, 10, 10),
			wantRect: image.Rect(3, 3, 8, 8),
		},
		{
			rect:     image.Rect(0, 0, 10, 10),
			size:     image.Pt(5, 5),
			limit:    image.Rect(0, 4, 7, 10),
			wantRect: image.Rect(2, 4, 7, 9),
		},
		{
			rect:     image.Rect(0, 0, 10, 10),
			size:     image.Pt(5, 5),
			limit:    image.Rect(4, -8, 8, -4),
			wantRect: image.Rect(4, -8, 8, -4),
		},
		{
			rect:     image.Rect(0, 0, 10, 10),
			size:     image.Pt(50, 50),
			limit:    image.Rect(-100, -100, 100, 100),
			wantRect: image.Rect(-20, -20, 30, 30),
		},
		{
			rect:     image.Rect(0, 0, 10, 10),
			size:     image.Pt(5, 5),
			limit:    image.Rect(5, 5, 10, 10),
			wantRect: image.Rect(5, 5, 10, 10),
		},
	}

	for i, tt := range tests {
		gotRect := image_.ScaleRectangleBySize(tt.rect, tt.size, tt.limit)
		if !tt.wantRect.Eq(gotRect) {
			t.Errorf("#%d: expected %s got %s", i, tt.wantRect, gotRect)
		}
	}
}

func TestRectangle2f_ScaleByFactor(t *testing.T) {
	tests := []struct {
		rect     image_.Rectangle2f
		factor   image_.Point2f
		wantRect image_.Rectangle2f
	}{
		{
			rect:     image_.Rect2f(0, 0, 10, 10),
			factor:   image_.Pt2f(1, 1),
			wantRect: image_.Rect2f(0, 0, 10, 10),
		},
		{
			rect:     image_.Rect2f(0, 0, 10, 10),
			factor:   image_.Pt2f(-1, -1),
			wantRect: image_.Rect2f(10, 10, 0, 0),
		},
		{
			rect:     image_.Rect2f(0, 0, 10, 10),
			factor:   image_.Pt2f(2, 0.5),
			wantRect: image_.Rect2f(-5, 2.5, 15, 7.5),
		},
		{
			rect:     image_.Rect2f(0, 0, 10, 10),
			factor:   image_.Pt2f(-2, -0.5),
			wantRect: image_.Rect2f(15, 7.5, -5, 2.5),
		},
		{
			rect:     image_.Rect2f(-10, -10, 20, 20),
			factor:   image_.Pt2f(0, 0),
			wantRect: image_.Rect2f(5, 5, 5, 5),
		},
	}

	for i, tt := range tests {
		gotRect := tt.rect.ScaleByFactor(tt.factor)
		if !tt.wantRect.RoundRectangle().Eq(gotRect.RoundRectangle()) {
			t.Errorf("#%d: expected %s got %s", i, tt.wantRect, gotRect)
		}
	}
}

func TestRectangle2f_UnionPoints(t *testing.T) {
	tests := []struct {
		rect     image_.Rectangle2f
		pts      []image_.Point2f
		wantRect image_.Rectangle2f
	}{
		{
			rect:     image_.Rect2f(0, 0, 10, 10),
			pts:      []image_.Point2f{image_.Pt2f(1, 1), image_.Pt2f(-1, -1), image_.Pt2f(20, 20)},
			wantRect: image_.Rect2f(-1, -1, 20, 20),
		},
		{
			rect:     image_.Rect2f(0, 0, 10, 10),
			pts:      []image_.Point2f{},
			wantRect: image_.Rect2f(0, 0, 10, 10),
		},
	}

	for i, tt := range tests {
		gotRect := tt.rect.UnionPoints(tt.pts...)
		if !tt.wantRect.RoundRectangle().Eq(gotRect.RoundRectangle()) {
			t.Errorf("#%d: expected %s got %s", i, tt.wantRect, gotRect)
		}
	}
}
