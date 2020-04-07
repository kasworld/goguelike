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

package resourcetilearea

import (
	"fmt"
	"image"
	"os"

	"github.com/kasworld/goguelike/game/terrain/resourcetile"
	"github.com/kasworld/wrapper"
)

type ResourceTileArea [][]resourcetile.ResourceTile

func New(x, y int) ResourceTileArea {
	rta := make(ResourceTileArea, x)
	for i, _ := range rta {
		rta[i] = make([]resourcetile.ResourceTile, y)
	}
	return rta
}
func (rta ResourceTileArea) String() string {
	return fmt.Sprintf("ResourceTileArea[%v %v]", len(rta), len(rta[0]))
}

func (rta ResourceTileArea) Dup() ResourceTileArea {
	w, h := len(rta), len(rta[0])
	rtn := make(ResourceTileArea, w)
	for x, _ := range rtn {
		rtn[x] = make([]resourcetile.ResourceTile, h)
		copy(rtn[x], rta[x])
	}
	return rtn
}

// for draw2d
func (rta ResourceTileArea) GetXYLen() (int, int) {
	return len(rta), len(rta[0])
}

// for draw2d
func (rta ResourceTileArea) OpXY(x, y int, v interface{}) {
	rv := v.(resourcetile.ResourceTypeValue)
	rta[x][y][rv.T] = rv.V
}

func (rta ResourceTileArea) GetXY(x, y int) *resourcetile.ResourceTile {
	return &rta[x][y]
}

////////

var list8Ways = [8][2]int{ // clock wise
	{0, -1},
	{1, -1},
	{1, 0},
	{1, 1},
	{0, 1},
	{-1, 1},
	{-1, 0},
	{-1, -1},
}

func (rta ResourceTileArea) Ageing(
	rndIntNFn func(n int) int,
	count int) ResourceTileArea {
	Xlen, Ylen := len(rta), len(rta[0])
	xWrapper := wrapper.New(Xlen).Wrap
	yWrapper := wrapper.New(Ylen).Wrap
	for i := 0; i < count; i++ {
		for x := 0; x < Xlen; x++ {
			for y, _ := range rta[x] {
				d := list8Ways[rndIntNFn(8)]
				x2 := xWrapper(x + d[0])
				y2 := yWrapper(y + d[1])

				t1 := &rta[x][y]
				t2 := &rta[x2][y2]
				t1.AgeingWith(t2, rndIntNFn)
				t1.ActTurnGeneration(rndIntNFn)
				t2.ActTurnGeneration(rndIntNFn)
			}
		}
	}
	return rta
}

func (rta ResourceTileArea) FromImage(filename string) error {
	Xlen, Ylen := len(rta), len(rta[0])
	reader, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer reader.Close()

	m, _, err := image.Decode(reader)
	if err != nil {
		return err
	}

	bounds := m.Bounds()
	if Xlen != bounds.Dx() || Ylen != bounds.Dy() {
		return fmt.Errorf(
			"image %v size %v %v not match %v %v",
			filename, bounds.Dx(), bounds.Dy(), Xlen, Ylen)
	}
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := m.At(x, y).RGBA()
			// A color's RGBA method returns values in the range [0, 65535].
			rta[x][y].FromRGBA(r, g, b, a)
		}
	}
	return nil
}
