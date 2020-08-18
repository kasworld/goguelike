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

	roomW, roomH := tr.rnd.NormIntRange(RoomMeanSize, Stddev), tr.rnd.NormIntRange(RoomMeanSize, Stddev)
	roomW = roomW/RoomGrid*RoomGrid + 1
	roomH = roomH/RoomGrid*RoomGrid + 1
	if roomW < Min {
		roomW = Min
	}
	if roomH < Min {
		roomH = Min
	}
	roomX, roomY := tr.rnd.IntRange(0, tr.Xlen-roomW), tr.rnd.IntRange(0, tr.Ylen-roomH)
	roomX = roomX / RoomGrid * RoomGrid
	roomY = roomY / RoomGrid * RoomGrid
	roomX, roomY = roomX%tr.Xlen, roomY%tr.Ylen
	return rect.Rect{roomX, roomY, roomW, roomH}
}
func (tr *Terrain) addRoomManual(rt rect.Rect, bgtile, walltile tile.Tile, terrace bool) error {
	r := room.New(rt, bgtile)
	r.DrawRectWall(tr.rnd, walltile, terrace)
	return tr.roomManager.AddRoom(r)
}

func (tr *Terrain) addMazeRoomManual(rt rect.Rect, bgtile, walltile tile.Tile, terrace bool,
	xn, yn int, connerFill bool,
) error {

	r := room.New(rt, bgtile)
	if err := r.DrawMaze(tr.rnd, xn, yn, walltile, connerFill); err != nil {
		return err
	}
	r.DrawRectWall(tr.rnd, walltile, false)
	return tr.roomManager.AddRoom(r)
}
