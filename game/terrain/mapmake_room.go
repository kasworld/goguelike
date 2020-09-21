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

	"github.com/kasworld/goguelike/enum/tile_flag"
	"github.com/kasworld/goguelike/game/terrain/room"
	"github.com/kasworld/goguelike/lib/scriptparse"
	"github.com/kasworld/rect"
)

func cmdAddRoom(tr *Terrain, ca *scriptparse.CmdArgs) error {
	var bgtiles, walltiles tile_flag.TileFlag
	var terrace bool
	var x, w, y, h int
	if err := ca.GetArgs(&bgtiles, &walltiles, &terrace, &x, &y, &w, &h); err != nil {
		return err
	}
	if err := tr.addRoomManual(rect.Rect{x, y, w, h}, bgtiles, walltiles, terrace); err != nil {
		return fmt.Errorf("Room add fail %v", err)
	}
	return nil
}

func cmdAddMazeRoom(tr *Terrain, ca *scriptparse.CmdArgs) error {
	var bgtiles, walltiles tile_flag.TileFlag
	var terrace bool
	var x, w, y, h int
	var xn, yn int
	var conerFill bool
	if err := ca.GetArgs(&bgtiles, &walltiles, &terrace, &x, &y, &w, &h, &xn, &yn, &conerFill); err != nil {
		return err
	}
	if err := tr.addMazeRoomManual(
		rect.Rect{x, y, w, h}, bgtiles, walltiles, terrace, xn, yn, conerFill); err != nil {
		return fmt.Errorf("Room add fail %v", err)
	}
	return nil
}

func cmdAddRandRooms(tr *Terrain, ca *scriptparse.CmdArgs) error {
	var bgtiles, walltiles tile_flag.TileFlag
	var terrace bool
	var align int
	var count, mean, stddev int
	if err := ca.GetArgs(&bgtiles, &walltiles, &terrace, &align, &count, &mean, &stddev); err != nil {
		return err
	}
	try := count
	for count > 0 && try > 0 {
		rt := tr.randRoomRect2(align, mean, stddev, terrace)
		if err := tr.addRoomManual(rt, bgtiles, walltiles, terrace); err == nil {
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

func (tr *Terrain) roomRandPosSize(align, mean, stddev, trlen int, terrace bool) (int, int) {
	size := tr.rnd.NormIntRange(mean, stddev)
	size = size/align*align + 1
	if size < mean/2 {
		size = mean / 2
	}
	if size < 4 {
		size = 4
	}
	if terrace && size < 6 {
		size = 6
	}
	if size > mean*2 {
		size = mean * 2
	}
	pos := tr.rnd.Intn((trlen-size)/align) * align
	return pos, size
}

func (tr *Terrain) randRoomRect2(align, mean, stddev int, terrace bool) rect.Rect {
	x, w := tr.roomRandPosSize(align, mean, stddev, tr.Xlen, terrace)
	y, h := tr.roomRandPosSize(align, mean, stddev, tr.Ylen, terrace)
	return rect.Rect{x, y, w, h}
}

func (tr *Terrain) addRoomManual(rt rect.Rect, bgtiles, walltiles tile_flag.TileFlag, terrace bool) error {
	r := room.New(rt, bgtiles)
	r.DrawRectWall(tr.rnd, walltiles, terrace)
	return tr.roomManager.AddRoom(r)
}

func (tr *Terrain) addMazeRoomManual(rt rect.Rect, bgtiles, walltiles tile_flag.TileFlag, terrace bool,
	xn, yn int, connerFill bool,
) error {

	r := room.New(rt, bgtiles)
	if err := r.DrawMaze(tr.rnd, xn, yn, walltiles, connerFill); err != nil {
		return err
	}
	r.DrawRectWall(tr.rnd, walltiles, false)
	return tr.roomManager.AddRoom(r)
}
