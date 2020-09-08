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
	"github.com/kasworld/goguelike/lib/scriptparse"
	"github.com/kasworld/rect"
)

func cmdNewTerrain(tr *Terrain, ca *scriptparse.CmdArgs) error {
	var name string
	var w, h, aocount, pocount int
	var actTurnBoost float64
	if err := ca.GetArgs(&name, &w, &h, &aocount, &pocount, &actTurnBoost); err != nil {
		return err
	}
	return tr.execNewTerrain(name, w, h, aocount, pocount, actTurnBoost)
}

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

func cmdConnectRooms(tr *Terrain, ca *scriptparse.CmdArgs) error {
	var tl tile.Tile
	var connectCount int
	var allconnect bool
	var diagonal bool
	if err := ca.GetArgs(&tl, &connectCount, &allconnect, &diagonal); err != nil {
		return err
	}
	tr.connectRooms(connectCount, allconnect, tl, diagonal)
	return nil
}

func cmdFinalizeTerrain(tr *Terrain, ca *scriptparse.CmdArgs) error {
	tr.crpCache = nil
	tr.findList = nil
	tr.tileLayer.DrawRooms(tr.roomManager.GetRoomList())
	tr.tileLayer.DrawCorridors(tr.corridorList)
	tr.renderServiceTileArea()
	return nil
}
