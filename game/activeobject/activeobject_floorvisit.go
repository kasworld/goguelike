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

	"github.com/kasworld/goguelike/config/viewportdata"
	"github.com/kasworld/goguelike/game/gamei"
	"github.com/kasworld/goguelike/lib/g2id"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_idnoti"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
)

func (ao *ActiveObject) Noti_EnterFloor(f gamei.FloorI) {
	ao.currrentFloor = f
	ao.SetNeedTANoti()
	if _, exist := ao.uuid2VisitArea.GetByID(f.GetUUID()); !exist {
		ao.uuid2VisitArea.Add(f)
	}
	if aio := ao.ai; aio != nil {
		aio.ResetPlan()
	}
}

func (ao *ActiveObject) UpdateBySightMat2(
	f gamei.FloorI, vpCenterX, vpCenterY int,
	sightMat *viewportdata.ViewportSight2, sight float32) {
	va, exist := ao.uuid2VisitArea.GetByID(f.GetUUID())
	if !exist {
		ao.log.Fatal("floor not visited %v %v", ao, f.GetUUID())
		va = ao.uuid2VisitArea.Add(f)
	}
	va.UpdateBySightMat2(f.GetTerrain().GetTiles(), vpCenterX, vpCenterY, sightMat, sight)
}

func (ao *ActiveObject) forgetAnyFloor() error {
	for _, fuuid := range ao.uuid2VisitArea.GetIDList() {
		return ao.forgetFloorByUUID(fuuid)
	}
	return fmt.Errorf("no visit floor")
}

func (ao *ActiveObject) forgetFloorByUUID(fuuid string) error {
	if err := ao.uuid2VisitArea.Forget(fuuid); err != nil {
		return err
	}
	if aoconn := ao.clientConn; aoconn != nil {
		return ao.clientConn.SendNotiPacket(
			c2t_idnoti.ForgetFloor,
			&c2t_obj.NotiForgetFloor_data{
				FloorG2ID: g2id.NewFromString(fuuid),
			},
		)
	}
	return nil
}

func (ao *ActiveObject) makeFloorComplete(f gamei.FloorI) error {
	va, _ := ao.uuid2VisitArea.GetByID(f.GetUUID())
	va.MakeComplete()

	fi := f.ToPacket_FloorInfo()
	if aoconn := ao.clientConn; aoconn != nil {
		return ao.clientConn.SendNotiPacket(
			c2t_idnoti.FloorTiles,
			&c2t_obj.NotiFloorTiles_data{
				FI:    fi,
				Tiles: f.GetTerrain().GetTiles(),
			},
		)
	}
	return nil
}
