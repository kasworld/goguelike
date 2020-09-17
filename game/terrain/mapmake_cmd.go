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
	"github.com/kasworld/findnear"
	"github.com/kasworld/goguelike/game/terrain/corridor"
	"github.com/kasworld/goguelike/game/terrain/resourcetilearea"
	"github.com/kasworld/goguelike/game/terrain/roommanager"
	"github.com/kasworld/goguelike/game/tilearea"
	"github.com/kasworld/goguelike/lib/scriptparse"
	"github.com/kasworld/goguelike/lib/uuidposman"
	"github.com/kasworld/wrapper"
)

func cmdNewTerrain(tr *Terrain, ca *scriptparse.CmdArgs) error {
	var name string
	var w, h int
	var actTurnBoost float64
	if err := ca.GetArgs(&name, &w, &h, &actTurnBoost); err != nil {
		return err
	}
	return tr.execNewTerrain(name, w, h, actTurnBoost)
}

func cmdActiveObjectsRand(tr *Terrain, ca *scriptparse.CmdArgs) error {
	var aocount int
	if err := ca.GetArgs(&aocount); err != nil {
		return err
	}
	tr.ActiveObjCount = aocount
	return nil
}

func cmdCarryObjectsRand(tr *Terrain, ca *scriptparse.CmdArgs) error {
	var cocount int
	if err := ca.GetArgs(&cocount); err != nil {
		return err
	}
	tr.CarryObjCount = cocount
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

func (tr *Terrain) execNewTerrain(
	name string, w, h int, actturnboost float64) error {
	tr.Xlen, tr.Ylen = w, h

	tr.XWrapper = wrapper.New(tr.Xlen)
	tr.YWrapper = wrapper.New(tr.Ylen)
	tr.XWrap = tr.XWrapper.GetWrapFn()
	tr.YWrap = tr.YWrapper.GetWrapFn()

	tr.Name = name
	tr.ActTurnBoost = actturnboost
	tr.serviceTileArea = tilearea.New(w, h)
	tr.tileLayer = tilearea.New(w, h)
	tr.resourceTileArea = resourcetilearea.New(w, h)
	tr.roomManager = roommanager.New(w, h)
	tr.corridorList = make([]*corridor.Corridor, 0)
	tr.foPosMan = uuidposman.New(tr.Xlen, tr.Ylen)

	tr.initCrpCache()
	tr.findList = findnear.NewXYLenList(w, h)
	return nil
}
