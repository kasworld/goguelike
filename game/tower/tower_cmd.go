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

package tower

import (
	"github.com/kasworld/goguelike/enum/achievetype"
	"github.com/kasworld/goguelike/game/cmd2tower"
	"github.com/kasworld/goguelike/game/fieldobject"
	"github.com/kasworld/goguelike/game/gamei"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_error"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_idnoti"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
)

func (tw *Tower) processCmd2Tower(data interface{}) {
	switch pk := data.(type) {
	default:
		tw.log.Fatal("unknown tower cmd %#v", data)

	case *cmd2tower.ActiveObjEnterTower:
		pk.RspCh <- tw.Call_ActiveObjEnterTower(pk.ActiveObj)

	case *cmd2tower.ActiveObjResumeTower:
		pk.RspCh <- tw.Call_ActiveObjResumeTower(pk.ActiveObj)

	case *cmd2tower.ActiveObjSuspendFromTower:
		pk.RspCh <- tw.Call_ActiveObjSuspendFromTower(pk.ActiveObj)

	case *cmd2tower.ActiveObjLeaveTower:
		pk.RspCh <- tw.Call_ActiveObjLeaveTower(pk.ActiveObjUUID)

	case *cmd2tower.AdminFloorMove:
		pk.RspCh <- tw.Call_AdminFloorMove(pk.ActiveObj, pk.RecvPacket)

	case *cmd2tower.FloorMove:
		pk.RspCh <- tw.Call_FloorMove(pk.ActiveObj, pk.FloorUUID)

	case *cmd2tower.AdminTowerCmd:
		pk.RspCh <- tw.Call_AdminTowerCmd(pk.ActiveObj, pk.RecvPacket)

	case *cmd2tower.ActiveObjUsePortal:
		tw.Call_ActiveObjUsePortal(pk.ActiveObj, pk.SrcFloor, pk.P1, pk.P2)

	case *cmd2tower.ActiveObjTrapTeleport:
		tw.Call_ActiveObjTrapTeleport(pk.ActiveObj, pk.SrcFloor, pk.DstFloorName)

	case *cmd2tower.ActiveObjRebirth:
		tw.Call_ActiveObjRebirth(pk.ActiveObj)
	}
}

func (tw *Tower) Call_ActiveObjEnterTower(ao gamei.ActiveObjectI) error {
	if err := tw.ao2Floor.ActiveObjEnterTower(ao.GetHomeFloor(), ao); err != nil {
		tw.log.Error("%v", err)
		return err
	}
	if err := tw.id2ao.Add(ao); err != nil {
		tw.log.Fatal("%v", err)
	}
	aocon := ao.GetClientConn()
	if aocon == nil {
		tw.log.Fatal("ao connection nil %v", ao)
		return nil
	}
	if err := aocon.SendNotiPacket(
		c2t_idnoti.EnterTower,
		&c2t_obj.NotiEnterTower_data{
			TowerInfo: tw.towerInfo,
		},
	); err != nil {
		tw.log.Fatal("%v", err)
	}

	return nil // continue login
}

func (tw *Tower) Call_ActiveObjResumeTower(ao gamei.ActiveObjectI) error {
	f := ao.GetCurrentFloor()
	if f == nil {
		f = ao.GetHomeFloor()
	}
	if err := tw.ao2Floor.ActiveObjResumeToFloor(f, ao); err != nil {
		tw.log.Fatal("%v", err)
		return err
	}
	_, err := tw.id2aoSuspend.DelByUUID(ao.GetUUID())
	if err != nil {
		tw.log.Fatal("%v", err)
	}
	if err := tw.id2ao.Add(ao); err != nil {
		tw.log.Fatal("%v", err)
	}
	// send noti visitarea in ao
	aocon := ao.GetClientConn()
	if aocon == nil {
		tw.log.Fatal("ao connection nil %v", ao)
		return nil
	}
	if err := aocon.SendNotiPacket(
		c2t_idnoti.EnterTower,
		&c2t_obj.NotiEnterTower_data{
			TowerInfo: tw.towerInfo,
		},
	); err != nil {
		tw.log.Fatal("%v", err)
	}
	for _, f := range tw.floorMan.GetFloorList() {
		va := ao.GetVisitFloor(f.GetUUID())
		if va == nil {
			continue
		}
		fi := f.ToPacket_FloorInfo()
		if err := aocon.SendNotiPacket(
			c2t_idnoti.FloorTiles,
			&c2t_obj.NotiFloorTiles_data{
				FI:    fi,
				Tiles: f.GetTerrain().GetTiles().DupWithFilter(va.GetXYNolock),
			},
		); err != nil {
			tw.log.Error("%v", err)
		}
	}
	return nil // continue login
}

func (tw *Tower) Call_ActiveObjSuspendFromTower(ao gamei.ActiveObjectI) error {
	ao, err := tw.id2ao.DelByUUID(ao.GetUUID())
	if err != nil {
		tw.log.Fatal("%v", err)
		return nil
	}
	tw.ao2Floor.ActiveObjLeaveFloor(ao)
	if err := tw.id2aoSuspend.Add(ao); err != nil {
		tw.log.Fatal("%v", err)
	}
	return nil
}

func (tw *Tower) Call_ActiveObjLeaveTower(ActiveObjUUID string) error {
	ao, err := tw.id2ao.DelByUUID(ActiveObjUUID)
	if err != nil {
		tw.log.Fatal("%v", err)
		return nil
	}
	tw.ao2Floor.ActiveObjLeaveFloor(ao)
	return nil
}

func (tw *Tower) Call_AdminFloorMove(
	ActiveObj gamei.ActiveObjectI,
	RecvPacket *c2t_obj.ReqAdminFloorMove_data) c2t_error.ErrorCode {

	switch cmd := RecvPacket.Floor; cmd {
	case "Before":
		aoFloor := tw.ao2Floor.GetFloorByActiveObjID(ActiveObj.GetUUID())
		aoFloorIndex, err := tw.floorMan.GetFloorIndexByUUID(aoFloor.GetUUID())
		if err != nil {
			tw.log.Error("floor not found %v", aoFloor)
			return c2t_error.ObjectNotFound
		}
		dstFloor := tw.floorMan.GetFloorByIndexWrap(aoFloorIndex - 1)
		x, y, err := dstFloor.SearchRandomActiveObjPosInRoomOrRandPos()
		if err != nil {
			tw.log.Error("fail to find rand pos %v %v %v", dstFloor, ActiveObj, err)
		}
		if err := tw.ao2Floor.ActiveObjMoveToFloor(dstFloor, ActiveObj, x, y); err != nil {
			tw.log.Fatal("%v", err)
			return c2t_error.ActionProhibited
		}
		ActiveObj.GetAchieveStat().Inc(achievetype.Admin)
		return c2t_error.None
	case "Next":
		aoFloor := tw.ao2Floor.GetFloorByActiveObjID(ActiveObj.GetUUID())
		aoFloorIndex, err := tw.floorMan.GetFloorIndexByUUID(aoFloor.GetUUID())
		if err != nil {
			tw.log.Error("floor not found %v", aoFloor)
			return c2t_error.ObjectNotFound
		}
		dstFloor := tw.floorMan.GetFloorByIndexWrap(aoFloorIndex + 1)
		x, y, err := dstFloor.SearchRandomActiveObjPosInRoomOrRandPos()
		if err != nil {
			tw.log.Error("fail to find rand pos %v %v %v", dstFloor, ActiveObj, err)
		}
		if err := tw.ao2Floor.ActiveObjMoveToFloor(dstFloor, ActiveObj, x, y); err != nil {
			tw.log.Fatal("%v", err)
			return c2t_error.ActionProhibited
		}
		ActiveObj.GetAchieveStat().Inc(achievetype.Admin)
		return c2t_error.None
	default:
		dstFloor := tw.floorMan.GetFloorByUUID(cmd)
		if dstFloor == nil {
			tw.log.Error("floor not found %v", cmd)
			return c2t_error.ObjectNotFound
		}
		x, y, err := dstFloor.SearchRandomActiveObjPosInRoomOrRandPos()
		if err != nil {
			tw.log.Error("fail to find rand pos %v %v %v", dstFloor, ActiveObj, err)
		}
		if err := tw.ao2Floor.ActiveObjMoveToFloor(dstFloor, ActiveObj, x, y); err != nil {
			tw.log.Fatal("%v", err)
			return c2t_error.ActionProhibited
		}
		ActiveObj.GetAchieveStat().Inc(achievetype.Admin)
		return c2t_error.None
	}
}

func (tw *Tower) Call_FloorMove(
	ActiveObj gamei.ActiveObjectI, FloorUUID string) c2t_error.ErrorCode {

	dstFloor := tw.floorMan.GetFloorByUUID(FloorUUID)
	if dstFloor == nil {
		tw.log.Error("floor not found %v", FloorUUID)
		return c2t_error.ObjectNotFound
	}
	x, y, err := dstFloor.SearchRandomActiveObjPosInRoomOrRandPos()
	if err != nil {
		tw.log.Error("fail to find rand pos %v %v %v", dstFloor, ActiveObj, err)
	}
	if err := tw.ao2Floor.ActiveObjMoveToFloor(dstFloor, ActiveObj, x, y); err != nil {
		tw.log.Fatal("%v", err)
		return c2t_error.ActionProhibited
	}
	return c2t_error.None
}

func (tw *Tower) Call_AdminTowerCmd(ActiveObj gamei.ActiveObjectI,
	RecvPacket *c2t_obj.ReqAdminTowerCmd_data) c2t_error.ErrorCode {
	tw.log.Debug("%v %v", ActiveObj, RecvPacket)
	return c2t_error.None
}

func (tw *Tower) Call_ActiveObjUsePortal(
	ActiveObj gamei.ActiveObjectI,
	SrcFloor gamei.FloorI,
	P1, P2 *fieldobject.FieldObject,
) {
	srcPosMan := SrcFloor.GetActiveObjPosMan()
	if srcPosMan == nil {
		tw.log.Warn("pos man nil %v", SrcFloor)
		return
	}
	if srcPosMan.GetByUUID(ActiveObj.GetUUID()) == nil {
		tw.log.Warn("ActiveObj not in floor %v %v", ActiveObj, SrcFloor)
		return
	}
	dstFloor := tw.floorMan.GetFloorByName(P2.FloorName)
	if dstFloor == nil {
		tw.log.Fatal("dstFloor not found %v", P2.FloorName)
		return
	}
	tw.log.Debug("ActiveObjUsePortal %v %v to %v", ActiveObj, SrcFloor, dstFloor)
	if err := tw.ao2Floor.ActiveObjMoveToFloor(dstFloor, ActiveObj, P2.X, P2.Y); err != nil {
		tw.log.Fatal("%v", err)
	}
}

func (tw *Tower) Call_ActiveObjTrapTeleport(
	ActiveObj gamei.ActiveObjectI,
	SrcFloor gamei.FloorI,
	DstFloorName string,
) {
	dstFloor := tw.floorMan.GetFloorByName(DstFloorName)
	if dstFloor == nil {
		tw.log.Fatal("dstFloor not found %v", DstFloorName)
		return
	}
	tw.log.Debug("ActiveObjTrapTeleport %v to %v", ActiveObj, dstFloor)
	x, y, err := dstFloor.SearchRandomActiveObjPos()
	if err != nil {
		tw.log.Error("fail to find rand pos %v %v %v", dstFloor, ActiveObj, err)
	}
	if err := tw.ao2Floor.ActiveObjMoveToFloor(dstFloor, ActiveObj, x, y); err != nil {
		tw.log.Fatal("%v", err)
	}
}

func (tw *Tower) Call_ActiveObjRebirth(ao gamei.ActiveObjectI) {
	if ao.IsAlive() {
		tw.log.Fatal("ao is alive %v HP:%v/%v",
			ao, ao.GetHP(), ao.GetTurnData().HPMax)
	}
	var dstFloor gamei.FloorI
	// system ao respawn to home floor == not user ao
	if aoconn := ao.GetClientConn(); aoconn == nil {
		dstFloor = ao.GetHomeFloor()
	} else {
		dstFloor = ao.GetCurrentFloor()
	}
	if err := tw.ao2Floor.ActiveObjRebirthToFloor(dstFloor, ao); err != nil {
		tw.log.Fatal("%v", err)
		return
	}
	tw.log.Debug("ActiveObjRebirth %v to %v", ao, dstFloor)
}
