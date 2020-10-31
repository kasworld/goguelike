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

package clientai

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/kasworld/goguelike/config/leveldata"
	"github.com/kasworld/goguelike/game/clientfloor"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_gob"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_idcmd"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_idnoti"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_json"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_packet"
)

var DemuxNoti2ByteFnMap = [...]func(me interface{}, hd c2t_packet.Header, rbody []byte) error{
	c2t_idnoti.Invalid:        bytesRecvNotiFn_Invalid,        // Invalid make empty packet error
	c2t_idnoti.EnterTower:     bytesRecvNotiFn_EnterTower,     // EnterTower
	c2t_idnoti.EnterFloor:     bytesRecvNotiFn_EnterFloor,     // EnterFloor
	c2t_idnoti.LeaveFloor:     bytesRecvNotiFn_LeaveFloor,     // LeaveFloor
	c2t_idnoti.LeaveTower:     bytesRecvNotiFn_LeaveTower,     // LeaveTower
	c2t_idnoti.Ageing:         bytesRecvNotiFn_Ageing,         // Ageing          // floor
	c2t_idnoti.Death:          bytesRecvNotiFn_Death,          // Death
	c2t_idnoti.ReadyToRebirth: bytesRecvNotiFn_ReadyToRebirth, // ReadyToRebirth
	c2t_idnoti.Rebirthed:      bytesRecvNotiFn_Rebirthed,      // Rebirthed
	c2t_idnoti.Broadcast:      bytesRecvNotiFn_Broadcast,      // Broadcast       // global chat broadcast from web admin
	c2t_idnoti.VPObjList:      bytesRecvNotiFn_VPObjList,      // VPObjList       // in viewport, every turn
	c2t_idnoti.VPTiles:        bytesRecvNotiFn_VPTiles,        // VPTiles         // in viewport, when viewport changed only
	c2t_idnoti.FloorTiles:     bytesRecvNotiFn_FloorTiles,     // FloorTiles      // for rebuild known floor
	c2t_idnoti.FieldObjList:   bytesRecvNotiFn_FieldObjList,   // FieldObjList    // for rebuild known floor
	c2t_idnoti.FoundFieldObj:  bytesRecvNotiFn_FoundFieldObj,  // FoundFieldObj   // hidden field obj
	c2t_idnoti.ForgetFloor:    bytesRecvNotiFn_ForgetFloor,    // ForgetFloor
	c2t_idnoti.ActivateTrap:   bytesRecvNotiFn_ActivateTrap,   // ActivateTrap
}

func bytesRecvNotiFn_Invalid(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	return fmt.Errorf("Not implemented %v", hd)
}

func bytesRecvNotiFn_EnterTower(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_gob.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("Packet type miss match %v", rbody)
	}
	pkbody, ok := robj.(*c2t_obj.NotiEnterTower_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	cai, ok := me.(*ClientAI)
	if !ok {
		return fmt.Errorf("recvobj type mismatch %v", me)
	}
	cai.TowerInfo = pkbody.TowerInfo

	return nil
}
func bytesRecvNotiFn_LeaveTower(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_gob.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("Packet type miss match %v", rbody)
	}
	pkbody, ok := robj.(*c2t_obj.NotiLeaveTower_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	_ = pkbody
	cai, ok := me.(*ClientAI)
	if !ok {
		return fmt.Errorf("recvobj type mismatch %v", me)
	}
	_ = cai
	return nil
}

func bytesRecvNotiFn_EnterFloor(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_gob.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("Packet type miss match %v", rbody)
	}
	pkbody, ok := robj.(*c2t_obj.NotiEnterFloor_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	cai, ok := me.(*ClientAI)
	if !ok {
		return fmt.Errorf("recvobj type mismatch %v", me)
	}
	cai.FloorInfo = pkbody.FI
	cf, exist := cai.Name2ClientFloor[pkbody.FI.Name]
	if !exist {
		// new floor
		cf = clientfloor.New(pkbody.FI)
		cai.Name2ClientFloor[pkbody.FI.Name] = cf
	}
	cf.EnterFloor()

	return nil
}
func bytesRecvNotiFn_LeaveFloor(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_gob.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("Packet type miss match %v", rbody)
	}
	pkbody, ok := robj.(*c2t_obj.NotiLeaveFloor_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	_ = pkbody
	cai, ok := me.(*ClientAI)
	if !ok {
		return fmt.Errorf("recvobj type mismatch %v", me)
	}
	_ = cai
	return nil
}

func bytesRecvNotiFn_Ageing(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	return nil
}

func bytesRecvNotiFn_Death(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	cai, ok := me.(*ClientAI)
	if !ok {
		return fmt.Errorf("recvobj type mismatch %v", me)
	}
	_ = cai
	if cai.config.DisconnectOnDeath {
		cai.sendRecvStop()
	}
	return nil
}

func bytesRecvNotiFn_ReadyToRebirth(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	cai, ok := me.(*ClientAI)
	if !ok {
		return fmt.Errorf("recvobj type mismatch %v", me)
	}
	go cai.ReqWithRspFnWithAuth(c2t_idcmd.Rebirth,
		&c2t_obj.ReqRebirth_data{},
		func(hd c2t_packet.Header, rsp interface{}) error {
			return nil
		})
	return nil
}

func bytesRecvNotiFn_Rebirthed(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	return nil
}

func bytesRecvNotiFn_Broadcast(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	return nil
}

func bytesRecvNotiFn_VPObjList(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_gob.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("Packet type miss match %v", rbody)
	}
	pkbody, ok := robj.(*c2t_obj.NotiVPObjList_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	cai, ok := me.(*ClientAI)
	if !ok {
		return fmt.Errorf("recvobj type mismatch %v", me)
	}
	cai.OLNotiData = pkbody
	cai.ServerClientTimeDiff = pkbody.Time.Sub(time.Now())
	oldOLNotiData := cai.OLNotiData
	cai.OLNotiData = pkbody
	newOLNotiData := pkbody
	cai.onFieldObj = nil

	atomic.StoreInt32(&cai.movePacketPerTurn, 0)
	// cai.SentAction = nil
	cai.ServerJitter.ActByValue(pkbody.Time)

	c2t_obj.EquipClientByUUID(pkbody.ActiveObj.EquipBag).Sort()
	c2t_obj.PotionClientByUUID(pkbody.ActiveObj.PotionBag).Sort()
	c2t_obj.ScrollClientByUUID(pkbody.ActiveObj.ScrollBag).Sort()

	if oldOLNotiData != nil {
		cai.HPdiff = newOLNotiData.ActiveObj.HP - oldOLNotiData.ActiveObj.HP
		cai.SPdiff = newOLNotiData.ActiveObj.SP - oldOLNotiData.ActiveObj.SP
	}
	newLevel := int(leveldata.CalcLevelFromExp(float64(newOLNotiData.ActiveObj.Exp)))

	cai.playerActiveObjClient = nil
	if ainfo := cai.AccountInfo; ainfo != nil {
		for _, v := range pkbody.ActiveObjList {
			if v.UUID == cai.AccountInfo.ActiveObjUUID {
				cai.playerActiveObjClient = v
			}
		}
	}

	cai.IsOverLoad = newOLNotiData.ActiveObj.CalcWeight() >= leveldata.WeightLimit(newLevel)

	if cai.FloorInfo == nil {
		cai.log.Error("cai.FloorInfo not set")
		return nil
	}
	if cai.FloorInfo.Name != newOLNotiData.FloorName {
		cai.log.Error("not current floor objlist data %v %v",
			cai.currentFloor().FloorInfo.Name, newOLNotiData.FloorName,
		)
		return nil
	}

	cf, exist := cai.Name2ClientFloor[newOLNotiData.FloorName]
	if !exist {
		cai.log.Warn("floor not added %v", newOLNotiData.FloorName)
		return nil
	}
	for _, v := range newOLNotiData.FieldObjList {
		cf.FieldObjPosMan.AddOrUpdateToXY(v, v.X, v.Y)
	}

	playerX, playerY := cai.GetPlayerXY()
	if cai.playerActiveObjClient != nil && cf.IsValidPos(playerX, playerY) {
		cai.onFieldObj = cf.GetFieldObjAt(playerX, playerY)
		cai.actByControlMode()
	}
	return nil
}

func bytesRecvNotiFn_VPTiles(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_gob.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("Packet type miss match %v", rbody)
	}
	pkbody, ok := robj.(*c2t_obj.NotiVPTiles_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	cai, ok := me.(*ClientAI)
	if !ok {
		return fmt.Errorf("recvobj type mismatch %v", me)
	}

	if cai.FloorInfo == nil {
		cai.log.Warn("OrangeRed cai.FloorInfo not set")
		return nil
	}
	if cai.FloorInfo.Name != pkbody.FloorName {
		cai.log.Warn("not current floor vptile data %v %v",
			cai.currentFloor().FloorInfo.Name, pkbody.FloorName,
		)
		return nil
	}
	cf, exist := cai.Name2ClientFloor[pkbody.FloorName]
	if !exist {
		cai.log.Warn("floor not added %v", pkbody.FloorName)
		return nil
	}

	oldComplete := cf.Visited.IsComplete()
	if err := cf.UpdateFromViewportTile(pkbody, cai.ViewportXYLenList); err != nil {
		cai.log.Warn("%v", err)
		return nil
	}
	if !oldComplete && cf.Visited.IsComplete() {
		// just completed
	}

	return nil
}

func bytesRecvNotiFn_FloorTiles(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_gob.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("Packet type miss match %v", rbody)
	}
	pkbody, ok := robj.(*c2t_obj.NotiFloorTiles_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	cai, ok := me.(*ClientAI)
	if !ok {
		return fmt.Errorf("recvobj type mismatch %v", me)
	}
	cf, exist := cai.Name2ClientFloor[pkbody.FI.Name]
	if !exist {
		// new floor
		cf = clientfloor.New(pkbody.FI)
		cai.Name2ClientFloor[pkbody.FI.Name] = cf
	}

	oldComplete := cf.Visited.IsComplete()
	cf.ReplaceFloorTiles(pkbody)
	if !oldComplete && cf.Visited.IsComplete() {
		// floor complete
	}
	return nil
}

// FieldObjList    // for rebuild known floor
func bytesRecvNotiFn_FieldObjList(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_json.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("Packet type miss match %v", rbody)
	}
	recved, ok := robj.(*c2t_obj.NotiFieldObjList_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	return fmt.Errorf("Not implemented %v", recved)
}

func bytesRecvNotiFn_FoundFieldObj(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_gob.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("Packet type miss match %v", rbody)
	}
	pkbody, ok := robj.(*c2t_obj.NotiFoundFieldObj_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	cai, ok := me.(*ClientAI)
	if !ok {
		return fmt.Errorf("recvobj type mismatch %v", me)
	}
	fromFloor, exist := cai.Name2ClientFloor[pkbody.FloorName]
	if !exist {
		cai.log.Fatal("FoundFieldObj unknonw floor %v", pkbody)
		return fmt.Errorf("FoundFieldObj unknonw floor %v", pkbody)
	}
	if fromFloor.FieldObjPosMan.Get1stObjAt(pkbody.FieldObj.X, pkbody.FieldObj.Y) == nil {
		fromFloor.FieldObjPosMan.AddOrUpdateToXY(pkbody.FieldObj, pkbody.FieldObj.X, pkbody.FieldObj.Y)
	}
	return nil
}

func bytesRecvNotiFn_ForgetFloor(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_gob.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("Packet type miss match %v", rbody)
	}
	pkbody, ok := robj.(*c2t_obj.NotiForgetFloor_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	cai, ok := me.(*ClientAI)
	if !ok {
		return fmt.Errorf("recvobj type mismatch %v", me)
	}

	forgetFloor, exist := cai.Name2ClientFloor[pkbody.FloorName]
	if exist {
		forgetFloor.Forget()
	}
	return nil
}

func bytesRecvNotiFn_ActivateTrap(me interface{}, hd c2t_packet.Header, rbody []byte) error {
	robj, err := c2t_gob.UnmarshalPacket(hd, rbody)
	if err != nil {
		return fmt.Errorf("Packet type miss match %v", rbody)
	}
	pkbody, ok := robj.(*c2t_obj.NotiActivateTrap_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", robj)
	}
	cai, ok := me.(*ClientAI)
	if !ok {
		return fmt.Errorf("recvobj type mismatch %v", me)
	}
	_ = cai
	_ = pkbody
	// g2log.Debug("%v", robj)
	return nil
}
