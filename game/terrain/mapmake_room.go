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
	"fmt"

	"github.com/kasworld/goguelike/enum/tile"
	"github.com/kasworld/goguelike/game/terrain/room"
	"github.com/kasworld/goguelike/lib/scriptparse"
	"github.com/kasworld/rect"
)

func cmdAddRoom(tr *Terrain, ca *scriptparse.CmdArgs) error {
	var bgtile, walltile tile.Tile
	var terrace bool
	var x, w, y, h int
	if err := ca.GetArgs(&bgtile, &walltile, &terrace, &x, &y, &w, &h); err != nil {
		return err
	}
	if err := tr.addRoomManual(rect.Rect{x, y, w, h}, bgtile, walltile, terrace); err != nil {
		return fmt.Errorf("Room add fail %v", err)
	}
	return nil
}

func cmdAddMazeRoom(tr *Terrain, ca *scriptparse.CmdArgs) error {
	var bgtile, walltile tile.Tile
	var terrace bool
	var x, w, y, h int
	var xn, yn int
	var conerFill bool
	if err := ca.GetArgs(&bgtile, &walltile, &terrace, &x, &y, &w, &h, &xn, &yn, &conerFill); err != nil {
		return err
	}
	if err := tr.addMazeRoomManual(
		rect.Rect{x, y, w, h}, bgtile, walltile, terrace, xn, yn, conerFill); err != nil {
		return fmt.Errorf("Room add fail %v", err)
	}
	return nil
}

func cmdAddRandRooms(tr *Terrain, ca *scriptparse.CmdArgs) error {
	var bgtile, walltile tile.Tile
	var terrace bool
	var align int
	var count, mean, stddev, min int
	if err := ca.GetArgs(&bgtile, &walltile, &terrace, &align, &count, &mean, &stddev, &min); err != nil {
		return err
	}
	try := count
	for count > 0 && try > 0 {
		rt := tr.randRoomRect2(align, mean, stddev, min)
		if err := tr.addRoomManual(rt, bgtile, walltile, terrace); err == nil {
			count--
		} else {
			try--
		}
	}
	if try == 0 {
		tr.log.Warn("Room add insufficient")
		// return false
	}
	return nil
}

func (tr *Terrain) randRoomRect2(align, RoomMeanSize, Stddev, Min int) rect.Rect {

	roomW, roomH := tr.rnd.NormIntRange(RoomMeanSize, Stddev), tr.rnd.NormIntRange(RoomMeanSize, Stddev)
	roomW = roomW/align*align + 1
	roomH = roomH/align*align + 1
	if roomW < Min {
		roomW = Min
	}
	if roomH < Min {
		roomH = Min
	}
	roomX := tr.rnd.Intn((tr.Xlen-roomW)/align) * align
	roomY := tr.rnd.Intn((tr.Ylen-roomH)/align) * align
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
