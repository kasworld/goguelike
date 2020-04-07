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

package floor

import (
	"github.com/kasworld/goguelike/enum/achievetype"
	"github.com/kasworld/goguelike/game/cmd2floor"
	"github.com/kasworld/goguelike/game/gamei"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_error"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_idnoti"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
)

func (f *Floor) processCmd2Floor(data interface{}) {
	switch pk := data.(type) {
	default:
		f.log.Fatal("unknown pk recv %v %#v", f, data)

	case *cmd2floor.ReqLeaveFloor:
		if err := f.aoPosMan.Del(pk.ActiveObj); err != nil {
			f.log.Fatal("%v %v", f, err)
		}
		if conn := pk.ActiveObj.GetClientConn(); conn != nil {
			if err := conn.SendNotiPacket(
				c2t_idnoti.LeaveFloor,
				&c2t_obj.NotiLeaveFloor_data{
					FI: f.ToPacket_FloorInfo(),
				},
			); err != nil {
				f.log.Error("%v %v", f, err)
			}
		}

	case *cmd2floor.ReqEnterFloor:
		err := f.aoPosMan.AddOrUpdateToXY(pk.ActiveObj, pk.X, pk.Y)
		if err != nil {
			f.log.Fatal("%v %v", f, err)
		}
		pk.ActiveObj.Noti_EnterFloor(f)
		if conn := pk.ActiveObj.GetClientConn(); conn != nil {
			if err := conn.SendNotiPacket(
				c2t_idnoti.EnterFloor,
				&c2t_obj.NotiEnterFloor_data{
					FI: f.ToPacket_FloorInfo(),
				},
			); err != nil {
				f.log.Error("%v %v", f, err)
			}
		}

	case *cmd2floor.ReqRebirth2Floor:
		err := f.aoPosMan.AddOrUpdateToXY(pk.ActiveObj, pk.X, pk.Y)
		if err != nil {
			f.log.Fatal("%v %v", f, err)
		}
		pk.ActiveObj.Noti_EnterFloor(f)
		pk.ActiveObj.Noti_Rebirth()
		if conn := pk.ActiveObj.GetClientConn(); conn != nil {
			if err := conn.SendNotiPacket(
				c2t_idnoti.Rebirthed,
				&c2t_obj.NotiRebirthed_data{},
			); err != nil {
				f.log.Error("%v %v %v", f, pk.ActiveObj, err)
			}
		}

	case *cmd2floor.APIAdminTeleport2Floor:
		pk.RspCh <- f.Call_APIAdminTeleport2Floor(pk.ActiveObj, pk.ReqPk)

	case *cmd2floor.APIAdminCmd2Floor:
		pk.RspCh <- f.Call_APIAdminCmd2Floor(pk.ActiveObj, pk.ReqPk)

	}
}

func (f *Floor) Call_APIAdminTeleport2Floor(
	ActiveObj gamei.ActiveObjectI, ReqPk *c2t_obj.ReqAdminTeleport_data) c2t_error.ErrorCode {

	if f.aoPosMan.GetByUUID(ActiveObj.GetUUID()) == nil {
		f.log.Warn("ActiveObj not in floor %v %v", f, ActiveObj)
		return c2t_error.ActionProhibited
	}
	x, y, err := f.SearchRandomActiveObjPos()
	if err != nil {
		f.log.Error("fail to teleport %v %v %v", f, ActiveObj, err)
		return c2t_error.ActionCanceled
	}
	x, y = f.terrain.WrapXY(x, y)
	if err := f.aoPosMan.UpdateToXY(ActiveObj, x, y); err != nil {
		f.log.Fatal("move ao fail %v %v %v", f, ActiveObj, err)
		return c2t_error.ActionCanceled
	}
	ActiveObj.SetNeedTANoti()
	ActiveObj.GetAchieveStat().Inc(achievetype.Admin)
	return c2t_error.None
}

func (f *Floor) Call_APIAdminCmd2Floor(
	ActiveObj gamei.ActiveObjectI, ReqPk *c2t_obj.ReqAdminFloorCmd_data) c2t_error.ErrorCode {
	return c2t_error.None
}
