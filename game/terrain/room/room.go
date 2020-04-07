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

package room

import (
	"fmt"

	"github.com/kasworld/g2rand"
	"github.com/kasworld/goguelike/enum/tile"
	"github.com/kasworld/goguelike/game/terrain/maze2"
	"github.com/kasworld/prettystring"
	"github.com/kasworld/rect"
	"github.com/kasworld/uuidstr"
)

func (r Room) String() string {
	return fmt.Sprintf("Room[%v Area:%v]", r.UUID, r.Area)
}

func (r *Room) IntRange(n1, n2 int) int {
	return r.rnd.Intn(n2-n1) + n1
}

type Room struct {
	rnd        *g2rand.G2Rand `prettystring:"hide"`
	UUID       string
	BgTile     tile.Tile
	Area       rect.Rect
	Tiles      [][]tile.Tile
	ConnectPos [][2]int // door outer pos , out of room area
	// for sort
	RecyclerCount int
	PortalCount   int
	TrapCount     int
}

func New(rt rect.Rect, bgTile tile.Tile) *Room {
	r := &Room{
		rnd:        g2rand.New(),
		UUID:       uuidstr.New(),
		BgTile:     bgTile,
		Area:       rt,
		Tiles:      make([][]tile.Tile, rt.W),
		ConnectPos: make([][2]int, 0, 4),
	}
	for x, _ := range r.Tiles {
		r.Tiles[x] = make([]tile.Tile, rt.H)
		for y := range r.Tiles[x] {
			r.Tiles[x][y] = bgTile
		}
	}
	return r
}

func (r *Room) DrawMaze(xn, yn int, walltile tile.Tile, connerFill bool) error {
	m := maze2.New(xn, yn)
	ma, err := m.ToMazeArea(r.Area.W-1, r.Area.H-1, connerFill)
	if err != nil {
		return fmt.Errorf("room %v %v", r, err)
	}
	for x, xv := range ma {
		for y, yv := range xv {
			if yv {
				r.Tiles[x][y] = walltile
			}
		}
	}
	return nil
}

func (r *Room) DrawRectWall(walltile tile.Tile, terrace bool) error {
	wallrect := r.Area
	shiftConnetPos := 1
	if terrace {
		wallrect = wallrect.Shrink([2]int{1, 1})
		shiftConnetPos = 2
	}
	r.DrawWall_N(wallrect, walltile)
	r.DrawWall_S(wallrect, walltile)
	r.DrawWall_W(wallrect, walltile)
	r.DrawWall_E(wallrect, walltile)
	r.AddWindowRand_N(wallrect, tile.Window)
	r.AddWindowRand_S(wallrect, tile.Window)
	r.AddWindowRand_W(wallrect, tile.Window)
	r.AddWindowRand_E(wallrect, tile.Window)
	r.AddDoorRand_N(wallrect, tile.Door, shiftConnetPos)
	r.AddDoorRand_S(wallrect, tile.Door, shiftConnetPos)
	r.AddDoorRand_W(wallrect, tile.Door, shiftConnetPos)
	r.AddDoorRand_E(wallrect, tile.Door, shiftConnetPos)
	return nil
}

func (r *Room) DrawWall_N(wallrect rect.Rect, walltile tile.Tile) {
	y := 0
	for x := 0; x < wallrect.W; x++ {
		r.Tiles[x][y] = walltile
	}
}
func (r *Room) DrawWall_S(wallrect rect.Rect, walltile tile.Tile) {
	y := wallrect.H - 1
	for x := 0; x < wallrect.W; x++ {
		r.Tiles[x][y] = walltile
	}
}
func (r *Room) DrawWall_W(wallrect rect.Rect, walltile tile.Tile) {
	x := 0
	for y := 0; y < wallrect.H; y++ {
		r.Tiles[x][y] = walltile
	}
}
func (r *Room) DrawWall_E(wallrect rect.Rect, walltile tile.Tile) {
	x := wallrect.W - 1
	for y := 0; y < wallrect.H; y++ {
		r.Tiles[x][y] = walltile
	}
}

func (r *Room) AddWindowRand_N(wallrect rect.Rect, wintile tile.Tile) {
	x := r.IntRange(1, wallrect.W-1)
	y := 0
	r.Tiles[x][y] = wintile
}
func (r *Room) AddWindowRand_S(wallrect rect.Rect, wintile tile.Tile) {
	x := r.IntRange(1, wallrect.W-1)
	y := wallrect.H - 1
	r.Tiles[x][y] = wintile
}
func (r *Room) AddWindowRand_W(wallrect rect.Rect, wintile tile.Tile) {
	x := 0
	y := r.IntRange(1, wallrect.H-1)
	r.Tiles[x][y] = wintile
}
func (r *Room) AddWindowRand_E(wallrect rect.Rect, wintile tile.Tile) {
	x := wallrect.W - 1
	y := r.IntRange(1, wallrect.H-1)
	r.Tiles[x][y] = wintile
}

func (r *Room) AddDoorRand_N(wallrect rect.Rect, doortile tile.Tile, shiftConnectPos int) {
	x := r.IntRange(1, wallrect.W-1)
	y := 0
	r.Tiles[x][y] = doortile
	r.ConnectPos = append(r.ConnectPos, [2]int{wallrect.X + x, wallrect.Y + y - shiftConnectPos})
}
func (r *Room) AddDoorRand_S(wallrect rect.Rect, doortile tile.Tile, shiftConnectPos int) {
	x := r.IntRange(1, wallrect.W-1)
	y := wallrect.H - 1
	r.Tiles[x][y] = doortile
	r.ConnectPos = append(r.ConnectPos, [2]int{wallrect.X + x, wallrect.Y + y + shiftConnectPos})
}
func (r *Room) AddDoorRand_W(wallrect rect.Rect, doortile tile.Tile, shiftConnectPos int) {
	x := 0
	y := r.IntRange(1, wallrect.H-1)
	r.Tiles[x][y] = doortile
	r.ConnectPos = append(r.ConnectPos, [2]int{wallrect.X + x - shiftConnectPos, wallrect.Y + y})
}
func (r *Room) AddDoorRand_E(wallrect rect.Rect, doortile tile.Tile, shiftConnectPos int) {
	x := wallrect.W - 1
	y := r.IntRange(1, wallrect.H-1)
	r.Tiles[x][y] = doortile
	r.ConnectPos = append(r.ConnectPos, [2]int{wallrect.X + x + shiftConnectPos, wallrect.Y + y})
}

// quadtree interface fn
// func (r *Room) GetRect() rect.Rect {
// 	return r.Area
// }

func (r *Room) GetUUID() string {
	return r.UUID
}

// // for draw2d
// func (r *Room) GetXYLen() (int, int) {
// 	return r.Area.W, r.Area.H
// }

// // for draw2d
// func (r *Room) OpXY(x, y int, v interface{}) {
// 	rv := v.(tile.Tile)
// 	r.Tiles[x-r.Area.X][y-r.Area.Y] = rv
// }

func (r *Room) RndDoorOuter() [2]int {
	return r.ConnectPos[r.rnd.Intn(len(r.ConnectPos))]
}

func (r *Room) RndPos() (int, int) {
	return r.Area.X + r.rnd.Intn(r.Area.W), r.Area.Y + r.rnd.Intn(r.Area.H)
}

func (info *Room) StringForm() string {
	return prettystring.PrettyString(info, 4)
}
