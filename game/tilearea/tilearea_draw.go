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

package tilearea

import (
	"github.com/kasworld/goguelike/enum/tile"
	"github.com/kasworld/goguelike/enum/tile_flag"
	"github.com/kasworld/goguelike/enum/tileoptype"
	"github.com/kasworld/goguelike/game/terrain/corridor"
	"github.com/kasworld/goguelike/game/terrain/room"
	"github.com/kasworld/goguelike/lib/boolmatrix"
	"github.com/kasworld/walk2d"
)

func (ta TileArea) DrawRooms(rs []*room.Room) {
	xWrap, yWrap := ta.GetXYWrapper()
	for _, r := range rs {
		roomRect := r.Area
		for x, xv := range r.Tiles {
			for y, yv := range xv {
				tax, tay := xWrap(roomRect.X+x), yWrap(roomRect.Y+y)
				ta[tax][tay].Op(tile_flag.TileTypeValue{Op: tileoptype.OverrideBits, Arg: r.BgTile})
				ta[tax][tay].Op(tile_flag.TileTypeValue{Op: tileoptype.OverrideBits, Arg: yv})
			}
		}
	}
}

func (ta TileArea) DrawCorridors(corridorList []*corridor.Corridor) {
	for _, v := range corridorList {
		for _, w := range v.P {
			ta[w[0]][w[1]].OverrideBits(v.Tile)
		}
	}
}

func (ta TileArea) HLine(x, w, y int, tl tile.Tile) {
	rv := tile_flag.TileTypeValue{Op: tileoptype.OverrideBits, Arg: tl}
	fn := func(ax, ay int) bool {
		ta.OpXY(ax, ay, rv)
		return false
	}
	walk2d.HLine(x, x+w, y, fn)
}

func (ta TileArea) VLine(x, y, h int, tl tile.Tile) {
	rv := tile_flag.TileTypeValue{Op: tileoptype.OverrideBits, Arg: tl}
	fn := func(ax, ay int) bool {
		ta.OpXY(ax, ay, rv)
		return false
	}
	walk2d.VLine(y, y+h, x, fn)
}

func (ta TileArea) Line(x1, y1, x2, y2 int, tl tile.Tile) {
	rv := tile_flag.TileTypeValue{Op: tileoptype.OverrideBits, Arg: tl}
	fn := func(ax, ay int) bool {
		ta.OpXY(ax, ay, rv)
		return false
	}
	walk2d.Line(x1, y1, x2, y2, fn)
}

func (ta TileArea) Rect(x, w, y, h int, tl tile.Tile) {
	rv := tile_flag.TileTypeValue{Op: tileoptype.OverrideBits, Arg: tl}
	fn := func(ax, ay int) bool {
		ta.OpXY(ax, ay, rv)
		return false
	}
	walk2d.Rect(x, y, x+w, y+h, fn)
}

func (ta TileArea) FillRect(x, w, y, h int, tl tile.Tile) {
	rv := tile_flag.TileTypeValue{Op: tileoptype.OverrideBits, Arg: tl}
	fn := func(ax, ay int) bool {
		ta.OpXY(ax, ay, rv)
		return false
	}
	walk2d.FillHV(x, y, x+w, y+h, fn)
}

func (ta TileArea) Ellipses(x, w, y, h int, tl tile.Tile) {
	rv := tile_flag.TileTypeValue{Op: tileoptype.OverrideBits, Arg: tl}
	fn := func(ax, ay int) bool {
		ta.OpXY(ax, ay, rv)
		return false
	}
	walk2d.Ellipses(x, y, x+w, y+h, fn)
}

func (ta TileArea) DrawBoolMapTrue(xWrap, yWrap func(i int) int,
	maX, maY int, ma boolmatrix.BoolMatrix, tl tile.Tile) {
	rv := tile_flag.TileTypeValue{Op: tileoptype.OverrideBits, Arg: tl}
	for x, xv := range ma {
		for y, yv := range xv {
			if yv {
				ta.OpXY(xWrap(maX+x), yWrap(maY+y), rv)
			}
		}
	}
}

func (ta TileArea) DrawBoolMapFalse(xWrap, yWrap func(i int) int,
	maX, maY int, ma boolmatrix.BoolMatrix, tl tile.Tile) {
	rv := tile_flag.TileTypeValue{Op: tileoptype.OverrideBits, Arg: tl}
	for x, xv := range ma {
		for y, yv := range xv {
			if !yv {
				ta.OpXY(xWrap(maX+x), yWrap(maY+y), rv)
			}
		}
	}
}
