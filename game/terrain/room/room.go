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
	"github.com/kasworld/goguelike/lib/maze2"
	"github.com/kasworld/prettystring"
	"github.com/kasworld/rect"
	"github.com/kasworld/uuidstr"
)

func (r Room) String() string {
	return fmt.Sprintf("Room[%v Area:%v]", r.UUID, r.Area)
}

type Room struct {
	UUID       string
	BgTile     tile.Tile
	Area       rect.Rect
	Tiles      [][]tile.Tile
	ConnectPos [][2]int // door outer pos , out of room area
	// for sort
	RecyclerCount   int
	PortalCount     int
	TrapCount       int
	AreaAttackCount int
}

func New(rt rect.Rect, bgTile tile.Tile) *Room {
	if rt.W < 4 || rt.H < 4 {
		panic(fmt.Sprintf("room to small %v", rt))
	}
	r := &Room{
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

func (r *Room) DrawMaze(rnd *g2rand.G2Rand, xn, yn int, walltile tile.Tile, connerFill bool) error {
	m := maze2.New(rnd, xn, yn)
	ma, err := m.ToBoolMatrix(r.Area.W-1, r.Area.H-1, connerFill)
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

func (r *Room) DrawRectWall(rnd *g2rand.G2Rand, walltile tile.Tile, terrace bool) error {
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
	r.AddWindowRand_N(rnd, wallrect, tile.Window)
	r.AddWindowRand_S(rnd, wallrect, tile.Window)
	r.AddWindowRand_W(rnd, wallrect, tile.Window)
	r.AddWindowRand_E(rnd, wallrect, tile.Window)
	r.AddDoorRand_N(rnd, wallrect, tile.Door, shiftConnetPos)
	r.AddDoorRand_S(rnd, wallrect, tile.Door, shiftConnetPos)
	r.AddDoorRand_W(rnd, wallrect, tile.Door, shiftConnetPos)
	r.AddDoorRand_E(rnd, wallrect, tile.Door, shiftConnetPos)
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

func (r *Room) AddWindowRand_N(rnd *g2rand.G2Rand, wallrect rect.Rect, wintile tile.Tile) {
	x := rnd.IntRange(1, wallrect.W-1)
	y := 0
	r.Tiles[x][y] = wintile
}
func (r *Room) AddWindowRand_S(rnd *g2rand.G2Rand, wallrect rect.Rect, wintile tile.Tile) {
	x := rnd.IntRange(1, wallrect.W-1)
	y := wallrect.H - 1
	r.Tiles[x][y] = wintile
}
func (r *Room) AddWindowRand_W(rnd *g2rand.G2Rand, wallrect rect.Rect, wintile tile.Tile) {
	x := 0
	y := rnd.IntRange(1, wallrect.H-1)
	r.Tiles[x][y] = wintile
}
func (r *Room) AddWindowRand_E(rnd *g2rand.G2Rand, wallrect rect.Rect, wintile tile.Tile) {
	x := wallrect.W - 1
	y := rnd.IntRange(1, wallrect.H-1)
	r.Tiles[x][y] = wintile
}

func (r *Room) AddDoorRand_N(rnd *g2rand.G2Rand, wallrect rect.Rect, doortile tile.Tile, shiftConnectPos int) {
	x := rnd.IntRange(1, wallrect.W-1)
	y := 0
	r.Tiles[x][y] = doortile
	r.ConnectPos = append(r.ConnectPos, [2]int{wallrect.X + x, wallrect.Y + y - shiftConnectPos})
}
func (r *Room) AddDoorRand_S(rnd *g2rand.G2Rand, wallrect rect.Rect, doortile tile.Tile, shiftConnectPos int) {
	x := rnd.IntRange(1, wallrect.W-1)
	y := wallrect.H - 1
	r.Tiles[x][y] = doortile
	r.ConnectPos = append(r.ConnectPos, [2]int{wallrect.X + x, wallrect.Y + y + shiftConnectPos})
}
func (r *Room) AddDoorRand_W(rnd *g2rand.G2Rand, wallrect rect.Rect, doortile tile.Tile, shiftConnectPos int) {
	x := 0
	y := rnd.IntRange(1, wallrect.H-1)
	r.Tiles[x][y] = doortile
	r.ConnectPos = append(r.ConnectPos, [2]int{wallrect.X + x - shiftConnectPos, wallrect.Y + y})
}
func (r *Room) AddDoorRand_E(rnd *g2rand.G2Rand, wallrect rect.Rect, doortile tile.Tile, shiftConnectPos int) {
	x := wallrect.W - 1
	y := rnd.IntRange(1, wallrect.H-1)
	r.Tiles[x][y] = doortile
	r.ConnectPos = append(r.ConnectPos, [2]int{wallrect.X + x + shiftConnectPos, wallrect.Y + y})
}

func (r *Room) GetUUID() string {
	return r.UUID
}

func (r *Room) RndDoorOuter(rnd *g2rand.G2Rand) [2]int {
	return r.ConnectPos[rnd.Intn(len(r.ConnectPos))]
}

func (r *Room) RndPos(rnd *g2rand.G2Rand) (int, int) {
	return r.Area.X + rnd.Intn(r.Area.W), r.Area.Y + rnd.Intn(r.Area.H)
}

func (info *Room) StringForm() string {
	return prettystring.PrettyString(info, 4)
}
