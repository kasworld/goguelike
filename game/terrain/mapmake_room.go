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

package terrain

import (
	"github.com/kasworld/goguelike/enum/tile"
	"github.com/kasworld/goguelike/game/terrain/room"
	"github.com/kasworld/rect"
)

func (tr *Terrain) randRoomRect2(RoomGrid, RoomMeanSize, Stddev, Min int) rect.Rect {

	rx, ry := tr.rnd.NormIntRange(RoomMeanSize, Stddev), tr.rnd.NormIntRange(RoomMeanSize, Stddev)
	if rx < Min {
		rx = Min
	}
	if ry < Min {
		ry = Min
	}
	rx = rx/RoomGrid*RoomGrid + 1
	ry = ry/RoomGrid*RoomGrid + 1
	stx, sty := tr.rnd.IntRange(0, tr.Xlen-rx), tr.rnd.IntRange(0, tr.Ylen-ry)
	stx = stx / RoomGrid * RoomGrid
	sty = sty / RoomGrid * RoomGrid
	stx, sty = stx%tr.Xlen, sty%tr.Ylen
	return rect.Rect{stx, sty, rx, ry}
}
func (tr *Terrain) addRoomManual(rt rect.Rect, bgtile, walltile tile.Tile, terrace bool) error {
	r := room.New(rt, bgtile)
	r.DrawRectWall(walltile, terrace)
	return tr.roomManager.AddRoom(r)
}

func (tr *Terrain) addMazeRoomManual(rt rect.Rect, bgtile, walltile tile.Tile, terrace bool,
	xn, yn int, connerFill bool,
) error {

	r := room.New(rt, bgtile)
	if err := r.DrawMaze(xn, yn, walltile, connerFill); err != nil {
		return err
	}
	r.DrawRectWall(walltile, false)
	return tr.roomManager.AddRoom(r)
}
