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

// 2d tile space direction
package direction

import (
	"math"

	"github.com/kasworld/go-abs"
)

type Direction_Type uint8

//go:generate stringer -type=Direction_Type
const (
	Dir_stop Direction_Type = iota
	Dir_n
	Dir_ne
	Dir_e
	Dir_se
	Dir_s
	Dir_sw
	Dir_w
	Dir_nw
	Direction_Count int = iota
)

func (i Direction_Type) String() string {
	return "Dir_" + attrib[i].Name
}

var attrib = []struct {
	Name string
	Vt   [2]int
	Len  float64
}{
	Dir_stop: {"Stop", [2]int{0, 0}, 0.0},
	Dir_n:    {"N", [2]int{0, -1}, 1.0},
	Dir_ne:   {"NE", [2]int{1, -1}, math.Sqrt2},
	Dir_e:    {"E", [2]int{1, 0}, 1.0},
	Dir_se:   {"SE", [2]int{1, 1}, math.Sqrt2},
	Dir_s:    {"S", [2]int{0, 1}, 1.0},
	Dir_sw:   {"SW", [2]int{-1, 1}, math.Sqrt2},
	Dir_w:    {"W", [2]int{-1, 0}, 1.0},
	Dir_nw:   {"NW", [2]int{-1, -1}, math.Sqrt2},
}

var vt2Dir = [3][3]Direction_Type{}

func init() {
	// println("init direction")
	for i, v := range attrib {
		vt2Dir[v.Vt[0]+1][v.Vt[1]+1] = Direction_Type(i)
	}
}

func (dt Direction_Type) Name() string {
	return attrib[dt].Name
}

func (dt Direction_Type) Vector() [2]int {
	return attrib[dt].Vt
}

func (dt Direction_Type) DxDy() (int, int) {
	return attrib[dt].Vt[0], attrib[dt].Vt[1]
}

func (dt Direction_Type) Dx() int {
	return attrib[dt].Vt[0]
}
func (dt Direction_Type) Dy() int {
	return attrib[dt].Vt[1]
}

func (dt Direction_Type) Len() float64 {
	return attrib[dt].Len
}

func (dt Direction_Type) TurnDir(turn int8) Direction_Type {
	if dt == Dir_stop {
		return dt
	}
	turn %= 8
	return Direction_Type((int8(dt)-1+turn+8)%8 + 1)
}

func (dt Direction_Type) ReverseDir() Direction_Type {
	vt := attrib[dt].Vt
	return Vt2Dir(-vt[0], -vt[1])
}
func (dt Direction_Type) InverseX() Direction_Type {
	vt := attrib[dt].Vt
	return Vt2Dir(-vt[0], vt[1])
}
func (dt Direction_Type) InverseY() Direction_Type {
	vt := attrib[dt].Vt
	return Vt2Dir(vt[0], -vt[1])
}

func (dt Direction_Type) IsValid() bool {
	return dt >= Dir_stop && dt <= Dir_nw
}

func (dt Direction_Type) Add(dt2 Direction_Type) Direction_Type {
	dx := dt.Dx() + dt2.Dx()
	dy := dt.Dy() + dt2.Dy()
	dx = abs.Signi(dx)
	dy = abs.Signi(dy)
	return Vt2Dir(dx, dy)
}

func (dt Direction_Type) MulXY(x, y int) (int, int) {
	return attrib[dt].Vt[0] * x, attrib[dt].Vt[1] * y
}

////

func Vt2Dir(x, y int) Direction_Type { // -1 ~ 1
	return vt2Dir[x+1][y+1]
}
func VtValidate(x, y int) bool {
	return x >= -1 && x <= 1 && y >= -1 && y <= 1
}

// func RandDir(rndfn func(st, ed int) int) Direction_Type {
// 	return Direction_Type(rndfn(1, 9))
// }

// find remote pos direction 8way
func DxDy2Dir8(dx, dy int) Direction_Type {
	dxs, dxv := abs.SignAbsi(dx)
	dys, dyv := abs.SignAbsi(dy)
	if dxv > dyv*2 {
		dys = 0
	}
	if dyv > dxv*2 {
		dxs = 0
	}
	return Vt2Dir(dxs, dys)
}

// find remote pos direction 4way
func DxDy2Dir4(dx, dy int) Direction_Type {
	dxs, dxv := abs.SignAbsi(dx)
	dys, dyv := abs.SignAbsi(dy)
	if dxv >= dyv {
		dys = 0
	} else {
		dxs = 0
	}
	return Vt2Dir(dxs, dys)
}

// contact only
func DxDy2Dir(dx, dy int) Direction_Type {
	dxs, _ := abs.SignAbsi(dx)
	dys, _ := abs.SignAbsi(dy)
	return Vt2Dir(dxs, dys)
}

// wrap dxdy
func CalcDxDyWrapped(dx, dy, xLen, yLen int) (int, int) {
	xs, xa := abs.SignAbsi(dx)
	if xa > xLen/2 {
		dx = xs * (xa - xLen)
	}
	ys, ya := abs.SignAbsi(dy)
	if ya > yLen/2 {
		dy = ys * (ya - yLen)
	}
	return dx, dy
}

// iscontact and dir with wrap
func CalcContactDirWrapped(from, to [2]int, xLen, yLen int) (bool, Direction_Type) {
	dx, dy := CalcDxDyWrapped(to[0]-from[0], to[1]-from[1], xLen, yLen)
	contact := VtValidate(dx, dy)
	if !contact {
		return false, 0
	}
	return true, Vt2Dir(dx, dy)
}

func CalcContactDirWrappedXY(x1, y1, x2, y2, xLen, yLen int) (bool, Direction_Type) {
	dx, dy := CalcDxDyWrapped(x2-x1, y2-y1, xLen, yLen)
	contact := VtValidate(dx, dy)
	if !contact {
		return false, 0
	}
	return true, Vt2Dir(dx, dy)
}

func Call8WayTile(ox, oy int, fn func(int, int) bool) []Direction_Type {
	TileDirs := []Direction_Type{}
	for i := Direction_Type(1); i <= 8; i++ {
		x, y := ox+i.Vector()[0], oy+i.Vector()[1]
		if fn(x, y) {
			TileDirs = append(TileDirs, i)
		}
	}
	return TileDirs
}
func Call4WayTile(ox, oy int, fn func(int, int) bool) []Direction_Type {
	TileDirs := []Direction_Type{}
	for i := Direction_Type(1); i <= 8; i += 2 {
		x, y := ox+i.Vector()[0], oy+i.Vector()[1]
		if fn(x, y) {
			TileDirs = append(TileDirs, i)
		}
	}
	return TileDirs
}
