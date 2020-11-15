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

package activeobject

import (
	"fmt"

	"github.com/kasworld/goguelike/config/gameconst"
	"github.com/kasworld/goguelike/config/viewportdata"
	"github.com/kasworld/goguelike/game/fieldobject"
	"github.com/kasworld/goguelike/game/gamei"
	"github.com/kasworld/goguelike/lib/uuidposmani"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_idnoti"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
)

func (ao *ActiveObject) EnterFloor(f gamei.FloorI) {
	ao.currentFloor = f
	ao.SetNeedTANoti()
	if _, exist := ao.floor4ClientMan.GetByName(f.GetName()); !exist {
		ao.floor4ClientMan.Add(f)
	}
	if aio := ao.ai; aio != nil {
		ao.ResetPlan(aio)
	}
}

func (ao *ActiveObject) UpdateVisitAreaBySightMat2(
	f gamei.FloorI, vpCenterX, vpCenterY int,
	sightMat *viewportdata.ViewportSight2, sight float32) {
	f4c, exist := ao.floor4ClientMan.GetByName(f.GetName())
	if !exist {
		ao.log.Fatal("floor not visited %v %v", ao, f.GetName())
		f4c = ao.floor4ClientMan.Add(f)
	}
	f4c.Visit.UpdateBySightMat2(
		f.GetTerrain().GetTiles(),
		vpCenterX, vpCenterY,
		sightMat,
		sight)
}

func (ao *ActiveObject) forgetAnyFloor() error {
	for _, floorName := range ao.floor4ClientMan.GetNameList() {
		return ao.ForgetFloorByName(floorName)
	}
	return fmt.Errorf("no visit floor")
}

func (ao *ActiveObject) ForgetFloorByName(floorName string) error {
	if err := ao.floor4ClientMan.Forget(floorName); err != nil {
		return err
	}
	if aoconn := ao.clientConn; aoconn != nil {
		return ao.clientConn.SendNotiPacket(c2t_idnoti.ForgetFloor,
			&c2t_obj.NotiForgetFloor_data{
				FloorName: floorName,
			},
		)
	}
	return nil
}

func (ao *ActiveObject) MakeFloorComplete(f gamei.FloorI) error {
	f4c, exist := ao.floor4ClientMan.GetByName(f.GetName())
	if !exist {
		f4c = ao.floor4ClientMan.Add(f)
	}

	f4c.Visit.MakeComplete()
	f.GetFieldObjPosMan().IterAll(func(o uuidposmani.UUIDPosI, foX, foY int) bool {
		fo := o.(*fieldobject.FieldObject)
		f4c.FOPosMan.AddOrUpdateToXY(fo.ToPacket_FieldObjClient(foX, foY), foX, foY)
		return false
	})

	fi := f.ToPacket_FloorInfo()
	if aoconn := ao.clientConn; aoconn != nil {
		// send tile area
		posList, taList := f.GetTerrain().GetTiles().Split(gameconst.TileAreaSplitSize)
		for i := range posList {
			if err := ao.clientConn.SendNotiPacket(c2t_idnoti.FloorTiles,
				&c2t_obj.NotiFloorTiles_data{
					FI:    fi,
					X:     posList[i][0],
					Y:     posList[i][1],
					Tiles: taList[i],
				},
			); err != nil {
				return err
			}
		}
		// send fieldobj list
		fol := make([]*c2t_obj.FieldObjClient, 0)
		f4c.FOPosMan.IterAll(func(o uuidposmani.UUIDPosI, foX, foY int) bool {
			fo := o.(*c2t_obj.FieldObjClient)
			fol = append(fol, fo)
			return false
		})
		if err := ao.clientConn.SendNotiPacket(c2t_idnoti.FieldObjList,
			&c2t_obj.NotiFieldObjList_data{
				FI:     fi,
				FOList: fol,
			},
		); err != nil {
			return err
		}
	}
	return nil
}
