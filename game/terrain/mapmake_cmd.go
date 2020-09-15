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
	"github.com/kasworld/goguelike/lib/scriptparse"
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

func cmdFinalizeTerrain(tr *Terrain, ca *scriptparse.CmdArgs) error {
	tr.crpCache = nil
	tr.findList = nil
	tr.tileLayer.DrawRooms(tr.roomManager.GetRoomList())
	tr.tileLayer.DrawCorridors(tr.corridorList)
	tr.renderServiceTileArea()
	return nil
}
