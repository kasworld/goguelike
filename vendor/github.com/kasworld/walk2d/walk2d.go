// Copyright 2015,2016,2017,2018,2019 SeukWon Kang (kasworld@gmail.com)
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package walk2d

import (
	"math"

	"github.com/kasworld/direction"
	"github.com/kasworld/go-abs"
)

type DoFn func(int, int) bool

// x2 not included
func HLine(x1, x2, y int, fn DoFn) bool {
	s := abs.Signi(x2 - x1)
	for x := x1; x != x2; x += s {
		if fn(x, y) {
			return true
		}
	}
	return false
}

// y2 not included
func VLine(y1, y2, x int, fn DoFn) bool {
	s := abs.Signi(y2 - y1)
	for y := y1; y != y2; y += s {
		if fn(x, y) {
			return true
		}
	}
	return false
}

// x2,y2 not included
func Line(x1, y1, x2, y2 int, fn DoFn) bool {
	diffx := x2 - x1
	diffy := y2 - y1
	xs, xa := abs.SignAbsi(x2 - x1)
	ys, ya := abs.SignAbsi(y2 - y1)
	if xa > ya { // horizontal
		for x := x1 + xs; x != x2; x += xs {
			y := y1 + (x-x1)*diffy/diffx
			if fn(x, y) {
				return true
			}
		}
	} else {
		for y := y1 + ys; y != y2; y += ys {
			x := x1 + (y-y1)*diffx/diffy
			if fn(x, y) {
				return true
			}
		}
	}
	return false
}

// x2,y2 not included
func Rect(x1, y1, x2, y2 int, fn DoFn) bool {
	x2--
	y2--
	if HLine(x1, x2, y1, fn) {
		return true
	}
	if VLine(y1, y2, x2, fn) {
		return true
	}
	if HLine(x2, x1, y2, fn) {
		return true
	}
	if VLine(y2, y1, x1, fn) {
		return true
	}
	return false
}

// x2,y2 not included
func FillHV(x1, y1, x2, y2 int, fn DoFn) bool {
	fn2 := func(x, y int) bool {
		return VLine(y1, y2, x, fn)
	}
	return HLine(x1, x2, y1, fn2)
}

// e not included
func Fill(x, y, s, e int, fn DoFn) bool {
	sign := abs.Signi(e - s)
	for i := s; i != e; i += sign {
		if Rect(x-i, y-i, x+i, y+i, fn) {
			return true
		}
	}
	return false
}

func Ellipses(x1, y1, x2, y2 int, fn DoFn) bool {
	cx, cy := (x1+x2)/2, (y1+y2)/2
	rx := (x2 - x1) / 2
	ry := (y2 - y1) / 2
	rate := float64(ry) / float64(rx)
	for x := -rx; x < rx; x++ {
		theta := math.Acos(float64(x) / float64(rx))
		y := int(float64(x) * math.Tan(theta) * rate)
		if x == 0 {
			y = ry
		}
		if VLine(cy-y, cy+y, cx+x, fn) {
			return true
		}
	}
	return false
}

func Near8Way(ox, oy int, fn DoFn) bool {
	for i := direction.Direction_Type(1); i <= 8; i++ {
		x, y := ox+i.Vector()[0], oy+i.Vector()[1]
		if fn(x, y) {
			return true
		}
	}
	return false
}
func Near4Way(ox, oy int, fn DoFn) bool {
	for i := direction.Direction_Type(1); i <= 8; i += 2 {
		x, y := ox+i.Vector()[0], oy+i.Vector()[1]
		if fn(x, y) {
			return true
		}
	}
	return false
}

// type PosList [][2]int
