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

package way9type

import (
	"math"

	"github.com/kasworld/go-abs"
)

func (dt Way9Type) Dx() int {
	return attrib[dt].Dx
}
func (dt Way9Type) Dy() int {
	return attrib[dt].Dy
}

func (dt Way9Type) Len() float64 {
	return attrib[dt].Len
}

func (dt Way9Type) DxDy() (int, int) {
	return dt.Dx(), dt.Dy()
}

func (dt Way9Type) TurnDir(turn int8) Way9Type {
	if dt == Center {
		return dt
	}
	turn %= 8
	return Way9Type((int8(dt)-1+turn+8)%8 + 1)
}

func (dt Way9Type) ReverseDir() Way9Type {
	return Vt2Dir(-dt.Dx(), -dt.Dy())
}
func (dt Way9Type) InverseX() Way9Type {
	return Vt2Dir(-dt.Dx(), dt.Dy())
}
func (dt Way9Type) InverseY() Way9Type {
	return Vt2Dir(dt.Dx(), -dt.Dy())
}

// is this need?
func (dt Way9Type) IsValid() bool {
	return dt >= Center && dt <= NorthWest
}

func (dt Way9Type) Add(dt2 Way9Type) Way9Type {
	dx := dt.Dx() + dt2.Dx()
	dy := dt.Dy() + dt2.Dy()
	dx = abs.Signi(dx)
	dy = abs.Signi(dy)
	return Vt2Dir(dx, dy)
}

func (dt Way9Type) MulXY(x, y int) (int, int) {
	return dt.Dx() * x, dt.Dy() * y
}

var attrib = [Way9Type_Count]struct {
	Dx  int
	Dy  int
	Len float64
}{
	Center:    {0, 0, 0.0},
	North:     {0, -1, 1.0},
	NorthEast: {1, -1, math.Sqrt2},
	East:      {1, 0, 1.0},
	SouthEast: {1, 1, math.Sqrt2},
	South:     {0, 1, 1.0},
	SouthWest: {-1, 1, math.Sqrt2},
	West:      {-1, 0, 1.0},
	NorthWest: {-1, -1, math.Sqrt2},
}

var vt2Dir = [3][3]Way9Type{}

func init() {
	for i, v := range attrib {
		vt2Dir[v.Dx+1][v.Dy+1] = Way9Type(i)
	}
}

////

func Vt2Dir(x, y int) Way9Type { // -1 ~ 1
	return vt2Dir[x+1][y+1]
}
func VtValidate(x, y int) bool {
	return x >= -1 && x <= 1 && y >= -1 && y <= 1
}

// find remote pos to 8way
func RemoteDxDy2Way9(dx, dy int) Way9Type {
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

// find remote pos to 4way
func RemoteDxDy2Way5(dx, dy int) Way9Type {
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
func ContactDxDy2Way9(dx, dy int) Way9Type {
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
func CalcContactDirWrapped(from, to [2]int, xLen, yLen int) (bool, Way9Type) {
	dx, dy := CalcDxDyWrapped(to[0]-from[0], to[1]-from[1], xLen, yLen)
	contact := VtValidate(dx, dy)
	if !contact {
		return false, 0
	}
	return true, Vt2Dir(dx, dy)
}

func CalcContactDirWrappedXY(x1, y1, x2, y2, xLen, yLen int) (bool, Way9Type) {
	dx, dy := CalcDxDyWrapped(x2-x1, y2-y1, xLen, yLen)
	contact := VtValidate(dx, dy)
	if !contact {
		return false, 0
	}
	return true, Vt2Dir(dx, dy)
}

func Call8WayTile(ox, oy int, fn func(int, int) bool) []Way9Type {
	TileDirs := []Way9Type{}
	for i := 1; i <= Way9Type_Count; i++ {
		w9 := Way9Type(i)
		x, y := ox+w9.Dx(), oy+w9.Dy()
		if fn(x, y) {
			TileDirs = append(TileDirs, w9)
		}
	}
	return TileDirs
}
func Call4WayTile(ox, oy int, fn func(int, int) bool) []Way9Type {
	TileDirs := []Way9Type{}
	for i := 1; i <= Way9Type_Count; i++ {
		w9 := Way9Type(i)
		x, y := ox+w9.Dx(), oy+w9.Dy()
		if fn(x, y) {
			TileDirs = append(TileDirs, w9)
		}
	}
	return TileDirs
}
