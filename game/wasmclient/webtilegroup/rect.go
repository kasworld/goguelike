// Copyright 2014,2015,2016,2017,2018,2019,2020 SeukWon Kang (kasworld@gmail.com)
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package webtilegroup tile group from json
package webtilegroup

import "github.com/kasworld/goguelike/enum/way9type"

type Rect struct {
	X int
	Y int
	W int
	H int
}

func (rt Rect) Shrink(n int) Rect {
	return Rect{
		X: rt.X + n,
		Y: rt.Y + n,
		W: rt.W - 2*n,
		H: rt.H - 2*n,
	}
}

func (rt Rect) CenterX() int {
	return rt.X + rt.W/2
}

func (rt Rect) CenterY() int {
	return rt.Y + rt.H/2
}

func (rt Rect) X2() int {
	return rt.X + rt.W
}

func (rt Rect) Y2() int {
	return rt.Y + rt.H
}

func (rt Rect) CenterXY() (int, int) {
	return rt.X + rt.W/2, rt.Y + rt.H/2
}

func (rt Rect) ShiftXY(shiftX, shiftY int) Rect {
	return Rect{
		rt.X + shiftX,
		rt.Y + shiftY,
		rt.W, rt.H,
	}
}

func (rt Rect) ShiftX(shiftX int) Rect {
	return Rect{
		rt.X + shiftX,
		rt.Y,
		rt.W, rt.H,
	}
}

func (rt Rect) ShiftY(shiftY int) Rect {
	return Rect{
		rt.X,
		rt.Y + shiftY,
		rt.W, rt.H,
	}
}

func (rt Rect) MoveToXY(X, Y int) Rect {
	return Rect{
		X, Y,
		rt.W, rt.H,
	}
}

func (rt Rect) MoveToX(X int) Rect {
	return Rect{
		X, rt.Y,
		rt.W, rt.H,
	}
}

func (rt Rect) MoveToY(Y int) Rect {
	return Rect{
		rt.X, Y,
		rt.W, rt.H,
	}
}

func (rt Rect) AddDir(dir way9type.Way9Type, amount int) Rect {
	shiftX := dir.Dx() * amount
	shiftY := dir.Dy() * amount
	return rt.ShiftXY(shiftX, shiftY)
}

func (rt Rect) AddDirXY(dir way9type.Way9Type, x, y int) Rect {
	shiftX := dir.Dx() * x
	shiftY := dir.Dy() * y
	return rt.ShiftXY(shiftX, shiftY)
}

func (rt Rect) ScaleCenterDxDy(dx, dy int) Rect {
	cx := rt.X + rt.W/2
	cy := rt.Y + rt.H/2
	w := rt.W - dx*2
	h := rt.H - dy*2
	return Rect{
		cx - w/2, cy - h/2,
		w, h,
	}
}

func (rt Rect) ScaleCenterRate(xrate float64, yrate float64) Rect {
	cx := rt.X + rt.W/2
	cy := rt.Y + rt.H/2
	w := int(float64(rt.W) * xrate)
	h := int(float64(rt.H) * yrate)
	return Rect{
		cx - w/2, cy - h/2,
		w, h,
	}
}

func ContactInner(outRect Rect, inRect Rect, dir way9type.Way9Type) Rect {
	rtn := inRect
	switch dir.Dx() {
	case -1:
		rtn.X = outRect.X
	case 0:
		rtn.X = outRect.X + outRect.W/2 - inRect.W/2
	case 1:
		rtn.X = outRect.X + outRect.W - inRect.W
	}
	switch dir.Dy() {
	case -1:
		rtn.Y = outRect.Y
	case 0:
		rtn.Y = outRect.Y + outRect.H/2 - inRect.H/2
	case 1:
		rtn.Y = outRect.Y + outRect.H - inRect.H
	}
	return rtn
}

// point base align
func (rt Rect) AroundToPoint(x, y int, dir way9type.Way9Type) Rect {
	rtn := Rect{x, y, rt.W, rt.H}
	rtn.X = x + rt.W/2*dir.Dx() - 1
	rtn.Y = y + rt.H/2*dir.Dy() - 1
	return rtn
}

func isInInt32(x, r1, r2 int) bool {
	return x >= r1 && x < r2
}

func (rt Rect) IsXYInRect(x, y int) bool {
	return isInInt32(x, rt.X, rt.X+rt.W) && isInInt32(y, rt.Y, rt.Y+rt.H)
}
