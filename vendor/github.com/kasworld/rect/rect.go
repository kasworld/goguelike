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

// rectangle in 2d space
package rect

type Rect struct {
	X, Y, W, H int
}

// func (r Rect) RandAxis(rndIntRange func(st, ed int) int, axis int) int {
// 	if axis == 0 {
// 		return rndIntRange(r.X, r.X+r.W)
// 	} else {
// 		return rndIntRange(r.Y, r.Y+r.H)
// 	}
// }
// func (r Rect) RandVector(rndIntRange func(st, ed int) int) [2]int {
// 	return [2]int{
// 		rndIntRange(r.X, r.X+r.W),
// 		rndIntRange(r.Y, r.Y+r.H),
// 	}
// }

// func (r Rect) RandXY(rndIntRange func(st, ed int) int) (int, int) {
// 	return rndIntRange(r.X, r.X+r.W), rndIntRange(r.Y, r.Y+r.H)
// }

func (r Rect) Center() [2]int {
	return [2]int{
		r.X + r.W/2, r.Y + r.H/2,
	}
}

func (r Rect) SizeVector() [2]int {
	return [2]int{r.W, r.H}
}

func (r Rect) Enlarge(size [2]int) Rect {
	return Rect{
		r.X - size[0], r.Y - size[1],
		r.W + size[0]*2, r.H + size[1]*2,
	}
}
func (r Rect) Shrink(size [2]int) Rect {
	return Rect{
		r.X + size[0], r.Y + size[1],
		r.W - size[0]*2, r.H - size[1]*2,
	}
}
func (r Rect) ShrinkSym(n int) Rect {
	return Rect{r.X + n, r.Y + n, r.W - n*2, r.H - n*2}
}

func (r Rect) X1() int {
	return r.X
}
func (r Rect) X2() int {
	return r.X + r.W
}
func (r Rect) Y1() int {
	return r.Y
}
func (r Rect) Y2() int {
	return r.Y + r.H
}

func (r Rect) MakeRectBy4Driect(c [2]int, direct4 int) Rect {
	w1, w2 := c[0]-r.X1(), r.X2()-c[0]
	h1, h2 := c[1]-r.Y1(), r.Y2()-c[1]
	switch direct4 {
	case 0:
		return Rect{r.X, r.Y, w1, h1}
	case 1:
		return Rect{c[0], r.Y, w2, h1}
	case 2:
		return Rect{c[0], c[1], w2, h2}
	case 3:
		return Rect{r.X, c[1], w1, h2}
	}
	return Rect{}
}

func (r1 Rect) IsOverlap(r2 Rect) bool {
	return !((r1.X1() >= r2.X2() || r1.X2() <= r2.X1()) ||
		(r1.Y1() >= r2.Y2() || r1.Y2() <= r2.Y1()))
}

func (r1 Rect) IsIn(r2 Rect) bool {
	if r1.X < r2.X || r1.X2() > r2.X2() {
		return false
	}
	if r1.Y < r2.Y || r1.Y2() > r2.Y2() {
		return false
	}
	return true
}

func min(i, j int) int {
	if i < j {
		return i
	}
	return j
}
func max(i, j int) int {
	if i > j {
		return i
	}
	return j
}

func (r1 Rect) Union(r2 Rect) Rect {
	r := Rect{
		min(r1.X1(), r2.X1()),
		min(r1.Y1(), r2.Y1()),
		0, 0,
	}
	r.W = max(r1.X2(), r2.X2()) - r.X1()
	r.H = max(r1.Y2(), r2.Y2()) - r.Y1()
	return r
}

func (r1 Rect) Intersection(r2 Rect) Rect {
	r := Rect{
		max(r1.X1(), r2.X1()),
		max(r1.Y1(), r2.Y1()),
		0, 0,
	}
	r.W = min(r1.X2(), r2.X2()) - r.X1()
	r.H = min(r1.Y2(), r2.Y2()) - r.Y1()
	return r
}

func (r Rect) Contain(p [2]int) bool {
	return r.X <= p[0] && p[0] < r.X2() &&
		r.Y <= p[1] && p[1] < r.Y2()
}

func (r Rect) RelPos(x, y int) (int, int) {
	return x - r.X, y - r.Y
}
