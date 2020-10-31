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

package wasmclientgl

import (
	"fmt"
	"sync/atomic"
	"syscall/js"
	"time"

	"github.com/kasworld/goguelike/config/gameconst"
	"github.com/kasworld/goguelike/config/leveldata"
	"github.com/kasworld/goguelike/enum/clientcontroltype"
	"github.com/kasworld/goguelike/enum/condition"
	"github.com/kasworld/goguelike/enum/fieldobjacttype"
	"github.com/kasworld/goguelike/enum/turnresulttype"
	"github.com/kasworld/goguelike/enum/way9type"
	"github.com/kasworld/goguelike/game/aoactreqrsp"
	"github.com/kasworld/goguelike/game/clientfloor"
	"github.com/kasworld/goguelike/game/soundmap"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_idcmd"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_idnoti"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_packet"
	"github.com/kasworld/gowasmlib/jslog"
	"github.com/kasworld/gowasmlib/wrapspan"
)

var ProcessRecvObjNotiFnMap = [...]func(recvobj interface{}, header c2t_packet.Header, body interface{}) error{
	c2t_idnoti.Invalid:        objRecvNotiFn_Invalid,        // Invalid make empty packet error
	c2t_idnoti.EnterTower:     objRecvNotiFn_EnterTower,     // EnterTower
	c2t_idnoti.EnterFloor:     objRecvNotiFn_EnterFloor,     // EnterFloor
	c2t_idnoti.LeaveFloor:     objRecvNotiFn_LeaveFloor,     // LeaveFloor
	c2t_idnoti.LeaveTower:     objRecvNotiFn_LeaveTower,     // LeaveTower
	c2t_idnoti.Ageing:         objRecvNotiFn_Ageing,         // Ageing          // floor
	c2t_idnoti.Death:          objRecvNotiFn_Death,          // Death
	c2t_idnoti.ReadyToRebirth: objRecvNotiFn_ReadyToRebirth, // ReadyToRebirth
	c2t_idnoti.Rebirthed:      objRecvNotiFn_Rebirthed,      // Rebirthed
	c2t_idnoti.Broadcast:      objRecvNotiFn_Broadcast,      // Broadcast       // global chat broadcast from web admin
	c2t_idnoti.VPObjList:      objRecvNotiFn_VPObjList,      // VPObjList       // in viewport, every turn
	c2t_idnoti.VPTiles:        objRecvNotiFn_VPTiles,        // VPTiles         // in viewport, when viewport changed only
	c2t_idnoti.FloorTiles:     objRecvNotiFn_FloorTiles,     // FloorTiles      // for rebuild known floor
	c2t_idnoti.FieldObjList:   objRecvNotiFn_FieldObjList,   // FieldObjList    // for rebuild known floor
	c2t_idnoti.FoundFieldObj:  objRecvNotiFn_FoundFieldObj,  // FoundFieldObj   // hidden field obj
	c2t_idnoti.ForgetFloor:    objRecvNotiFn_ForgetFloor,    // ForgetFloor
	c2t_idnoti.ActivateTrap:   objRecvNotiFn_ActivateTrap,   // ActivateTrap
}

// Invalid make empty packet error
func objRecvNotiFn_Invalid(me interface{}, hd c2t_packet.Header, body interface{}) error {
	robj, ok := body.(*c2t_obj.NotiInvalid_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", body)
	}
	return fmt.Errorf("Not implemented %v", robj)
}

func objRecvNotiFn_EnterTower(recvobj interface{}, header c2t_packet.Header, obj interface{}) error {
	robj, ok := obj.(*c2t_obj.NotiEnterTower_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", obj)
	}
	app, ok := recvobj.(*WasmClient)
	if !ok {
		return fmt.Errorf("recvobj type mismatch %v", recvobj)
	}

	soundmap.Play("enterfloorsound")
	gInitData.TowerInfo = robj.TowerInfo
	app.systemMessage.Append(wrapspan.ColorTextf("yellow",
		"Enter tower %v", robj.TowerInfo.Name))
	app.NotiMessage.AppendTf(tcsInfo,
		"Enter tower %v", robj.TowerInfo.Name)
	return nil
}
func objRecvNotiFn_LeaveTower(recvobj interface{}, header c2t_packet.Header, obj interface{}) error {
	robj, ok := obj.(*c2t_obj.NotiLeaveTower_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", obj)
	}
	app, ok := recvobj.(*WasmClient)
	if !ok {
		return fmt.Errorf("recvobj type mismatch %v", recvobj)
	}

	soundmap.Play("enterfloorsound")
	app.systemMessage.Append(wrapspan.ColorTextf("yellow",
		"Leave tower %v", robj.TowerInfo.Name))
	app.NotiMessage.AppendTf(tcsInfo,
		"Leave tower %v", robj.TowerInfo.Name)
	return nil
}

func objRecvNotiFn_EnterFloor(recvobj interface{}, header c2t_packet.Header, obj interface{}) error {
	robj, ok := obj.(*c2t_obj.NotiEnterFloor_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", obj)
	}
	app, ok := recvobj.(*WasmClient)
	if !ok {
		return fmt.Errorf("recvobj type mismatch %v", recvobj)
	}

	soundmap.Play("enterfloorsound")
	app.FloorInfo = robj.FI
	cf, exist := app.Name2ClientFloor[robj.FI.Name]
	if !exist {
		// new floor
		cf = clientfloor.New(robj.FI)
		app.Name2ClientFloor[robj.FI.Name] = cf
		app.systemMessage.Append(wrapspan.ColorTextf("yellow",
			"Found floor %v", cf.FloorInfo.Name))
		app.NotiMessage.AppendTf(tcsInfo,
			"Found floor %v", cf.FloorInfo.Name)
	}
	app.systemMessage.Appendf("Enter floor %v", cf.FloorInfo.Name)
	app.NotiMessage.AppendTf(tcsInfo,
		"Enter floor %v", cf.FloorInfo.Name)
	cf.EnterFloor()
	return nil
}
func objRecvNotiFn_LeaveFloor(recvobj interface{}, header c2t_packet.Header, obj interface{}) error {
	robj, ok := obj.(*c2t_obj.NotiLeaveFloor_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", obj)
	}
	app, ok := recvobj.(*WasmClient)
	if !ok {
		return fmt.Errorf("recvobj type mismatch %v", recvobj)
	}

	soundmap.Play("enterfloorsound")
	app.systemMessage.Appendf("Leave floor %v", robj.FI.Name)
	app.NotiMessage.AppendTf(tcsInfo,
		"Leave floor %v", robj.FI.Name)
	return nil
}

func objRecvNotiFn_Ageing(recvobj interface{}, header c2t_packet.Header, obj interface{}) error {
	robj, ok := obj.(*c2t_obj.NotiAgeing_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", obj)
	}
	app, ok := recvobj.(*WasmClient)
	if !ok {
		return fmt.Errorf("recvobj type mismatch %v", recvobj)
	}
	_ = robj
	_ = app

	soundmap.Play("ageingsound")
	app.systemMessage.Append("Floor aged")
	app.NotiMessage.AppendTf(tcsInfo,
		"Floor aged")
	return nil
}

func objRecvNotiFn_Death(recvobj interface{}, header c2t_packet.Header, obj interface{}) error {
	robj, ok := obj.(*c2t_obj.NotiDeath_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", obj)
	}
	app, ok := recvobj.(*WasmClient)
	if !ok {
		return fmt.Errorf("recvobj type mismatch %v", recvobj)
	}
	_ = robj

	soundmap.Play("diesound")
	app.systemMessage.Append("You died.")
	app.NotiMessage.AppendTf(tcsWarn,
		"You died.")
	app.remainTurn2Rebirth = gameconst.ActiveObjRebirthWaitTurn

	return nil
}

func objRecvNotiFn_ReadyToRebirth(recvobj interface{}, header c2t_packet.Header, obj interface{}) error {
	robj, ok := obj.(*c2t_obj.NotiReadyToRebirth_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", obj)
	}
	app, ok := recvobj.(*WasmClient)
	if !ok {
		return fmt.Errorf("recvobj type mismatch %v", recvobj)
	}
	_ = app
	_ = robj

	app.systemMessage.Append("Ready to rebirth")
	app.NotiMessage.AppendTf(tcsInfo,
		"Ready to rebirth")
	if autoActs.GetByIDBase("AutoRebirth").State == 0 {
		go app.sendPacket(c2t_idcmd.Rebirth,
			&c2t_obj.ReqRebirth_data{},
		)
	} else {
		commandButtons.GetByIDBase("Rebirth").Enable()
	}
	// clear carryobj cache
	app.CaObjUUID2CaObjClient = make(map[string]interface{})
	app.remainTurn2Rebirth = 0

	return nil
}

func objRecvNotiFn_Rebirthed(recvobj interface{}, header c2t_packet.Header, obj interface{}) error {
	robj, ok := obj.(*c2t_obj.NotiRebirthed_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", obj)
	}
	app, ok := recvobj.(*WasmClient)
	if !ok {
		return fmt.Errorf("recvobj type mismatch %v", recvobj)
	}
	_ = app
	_ = robj

	soundmap.Play("rebirthsound")
	app.systemMessage.Append("Back to life!")
	app.NotiMessage.AppendTf(tcsInfo,
		"Back to life!")
	commandButtons.GetByIDBase("Rebirth").Disable()
	return nil
}

func objRecvNotiFn_Broadcast(recvobj interface{}, header c2t_packet.Header, obj interface{}) error {
	robj, ok := obj.(*c2t_obj.NotiBroadcast_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", obj)
	}
	app, ok := recvobj.(*WasmClient)
	if !ok {
		return fmt.Errorf("recvobj type mismatch %v", recvobj)
	}
	_ = app
	_ = robj

	soundmap.Play("broadcastsound")
	app.systemMessage.Appendf("Broadcast %v", robj.Msg)
	app.NotiMessage.AppendTf(tcsInfo,
		"Broadcast %v", robj.Msg)
	return nil
}

func objRecvNotiFn_VPObjList(recvobj interface{}, header c2t_packet.Header, obj interface{}) error {
	robj, ok := obj.(*c2t_obj.NotiVPObjList_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", obj)
	}
	app, ok := recvobj.(*WasmClient)
	if !ok {
		return fmt.Errorf("recvobj type mismatch %v", recvobj)
	}

	defer func() {
		app.waitObjList = false
		if app.needRefreshSet {
			app.needRefreshSet = false
			js.Global().Call("requestAnimationFrame", js.FuncOf(app.renderGLFrame))
		}
	}()

	app.ServerClientTimeDiff = robj.Time.Sub(time.Now())
	app.olNotiHeader = header

	app.lastOLNotiData = app.olNotiData
	if app.lastOLNotiData == nil {
		app.lastOLNotiData = robj
	}

	app.olNotiData = robj
	newOLNotiData := robj
	app.onFieldObj = nil

	app.ClientJitter.Act()
	app.ServerJitter.ActByValue(robj.Time)

	c2t_obj.EquipClientByUUID(robj.ActiveObj.EquipBag).Sort()
	c2t_obj.PotionClientByUUID(robj.ActiveObj.PotionBag).Sort()
	c2t_obj.ScrollClientByUUID(robj.ActiveObj.ScrollBag).Sort()

	oldLevel := 0
	exp := 0

	// check change
	oldLevel = int(leveldata.CalcLevelFromExp(float64(app.lastOLNotiData.ActiveObj.Exp)))
	app.HPdiff = newOLNotiData.ActiveObj.HP - app.lastOLNotiData.ActiveObj.HP
	app.SPdiff = newOLNotiData.ActiveObj.SP - app.lastOLNotiData.ActiveObj.SP
	exp = newOLNotiData.ActiveObj.Exp - app.lastOLNotiData.ActiveObj.Exp

	if app.lastOLNotiData.ActiveObj.Bias.NearFaction() != robj.ActiveObj.Bias.NearFaction() {
		app.systemMessage.Appendf(
			"Faction changed to %v", robj.ActiveObj.Bias.NearFaction().String())
		app.NotiMessage.AppendTf(tcsInfo,
			"Faction changed to %v", robj.ActiveObj.Bias.NearFaction().String())
	}
	if app.lastOLNotiData.ActiveObj.Bias != robj.ActiveObj.Bias {
		app.systemMessage.Append("Bias changed")
	}

	newLevel := int(leveldata.CalcLevelFromExp(float64(newOLNotiData.ActiveObj.Exp)))
	if oldLevel < newLevel {
		soundmap.Play("levelupsound")
		app.systemMessage.Append(wrapspan.ColorTextf("yellow",
			"Level Up to %v", newLevel))
		app.NotiMessage.AppendTf(tcsInfo,
			"Level Up to %v", newLevel)
	} else if oldLevel > newLevel {
		soundmap.Play("leveldownsound")
		app.systemMessage.Append(wrapspan.ColorTextf("yellow",
			"Level Down to %v", newLevel))
		app.NotiMessage.AppendTf(tcsInfo,
			"Level Down to %v", newLevel)
	}
	if exp != 0 {
		app.NotiMessage.Appendf("yellow", calcSizeRateByValue(exp/10), time.Second, "Exp %+d", exp)
	}
	if app.HPdiff > 0 {
		app.NotiMessage.Appendf("LimeGreen", calcSizeRateByValue(app.HPdiff), time.Second, "HP %+d", app.HPdiff)
	} else if app.HPdiff < 0 {
		app.NotiMessage.Appendf("Red", calcSizeRateByValue(-app.HPdiff), time.Second, "HP %+d", app.HPdiff)
	}
	if app.SPdiff > 0 {
		app.NotiMessage.Appendf("LimeGreen", calcSizeRateByValue(app.SPdiff), time.Second, "SP %+d", app.SPdiff)
	} else if app.SPdiff < 0 {
		app.NotiMessage.Appendf("Red", calcSizeRateByValue(-app.SPdiff), time.Second, "SP %+d", app.SPdiff)
	}

	for i := 0; i < condition.Condition_Count; i++ {
		ci := condition.Condition(i)
		if newOLNotiData.ActiveObj.Conditions.TestByCondition(ci) {
			app.NotiMessage.Appendf("OrangeRed", 1.0, time.Second, "%s", ci)
		}
	}

	for _, v := range newOLNotiData.ActiveObj.EquipBag {
		app.CaObjUUID2CaObjClient[v.UUID] = v
	}
	for _, v := range newOLNotiData.ActiveObj.EquippedPo {
		app.CaObjUUID2CaObjClient[v.UUID] = v
	}
	for _, v := range newOLNotiData.ActiveObj.PotionBag {
		app.CaObjUUID2CaObjClient[v.UUID] = v
	}
	for _, v := range newOLNotiData.ActiveObj.ScrollBag {
		app.CaObjUUID2CaObjClient[v.UUID] = v
	}

	for _, v := range newOLNotiData.ActiveObjList {
		app.AOUUID2AOClient[v.UUID] = v
	}
	for _, v := range newOLNotiData.CarryObjList {
		if _, exist := app.CaObjUUID2CaObjClient[v.UUID]; !exist {
			app.CaObjUUID2CaObjClient[v.UUID] = v
		}
	}

	if actStr := app.ActionResult2String(newOLNotiData.ActiveObj.Act); actStr != "" {
		app.systemMessage.Append(actStr)
	}

	for _, v := range newOLNotiData.ActiveObj.TurnResult {
		app.processTurnResult(v)
	}

	SoundByActResult(newOLNotiData.ActiveObj.Act)
	wLimit := leveldata.WeightLimit(newLevel)
	if newOLNotiData.ActiveObj.Conditions.TestByCondition(condition.Burden) {
		wLimit /= 2
	}
	app.OverLoadRate = newOLNotiData.ActiveObj.CalcWeight() / wLimit
	if app.OverLoadRate > 1 {
		soundmap.Play("overloadsound")
		app.NotiMessage.Appendf("red", app.OverLoadRate, time.Second, "Overloaded")
	}
	if pao := app.GetPlayerActiveObjClient(); pao != nil && pao.DamageTake > 1 {
		soundmap.Play("damagesound")
	}

	if newOLNotiData.ActiveObj.HP <= 0 {
		if app.remainTurn2Rebirth > 0 {
			app.NotiMessage.AppendTf(tcsRebirth, "Remain to rebirth %v", app.remainTurn2Rebirth)
		} else {
			app.NotiMessage.AppendTf(tcsRebirth, "Ready to rebirth %v", -app.remainTurn2Rebirth)
		}
		app.remainTurn2Rebirth--
	}

	if app.FloorInfo == nil {
		jslog.Error("app.FloorInfo not set")
		return nil
	}
	if app.FloorInfo.Name != newOLNotiData.FloorName {
		jslog.Errorf("not current floor objlist data %v %v",
			app.currentFloor().FloorInfo.Name, newOLNotiData.FloorName,
		)
		return nil
	}

	cf, exist := app.Name2ClientFloor[newOLNotiData.FloorName]
	if !exist {
		jslog.Warnf("floor not added %v", newOLNotiData.FloorName)
		return nil
	}
	for _, v := range newOLNotiData.FieldObjList {
		cf.FieldObjPosMan.AddToXY(v, v.X, v.Y)
	}

	playerX, playerY := app.GetPlayerXY()
	// cf.updateFieldObjInView(playerX, playerY)
	app.vp.processNotiVPObjList(cf, newOLNotiData, playerX, playerY)
	if cf.IsValidPos(playerX, playerY) {
		app.onFieldObj = cf.GetFieldObjAt(playerX, playerY)
	}

	app.DisplayTextInfo()
	atomic.StoreInt32(&app.movePacketPerTurn, 0)
	atomic.StoreInt32(&app.actPacketPerTurn, 0)

	if !cf.IsValidPos(playerX, playerY) {
		jslog.Warnf("ao pos out of floor %v [%v %v]", cf, playerX, playerY)
		return nil
	}
	if newOLNotiData.ActiveObj.AP > 0 {
		app.actByControlMode()
	}

	return nil
}

func (app *WasmClient) processTurnResult(v c2t_obj.TurnResultClient) {
	switch v.ResultType {
	default:
		jslog.Error("unknown acttype %v", v)

	case turnresulttype.DropCarryObj:
		o, exist := app.CaObjUUID2CaObjClient[v.DstUUID]
		if exist {
			app.systemMessage.Appendf("Drop %v", obj2ColorStr(o))
		} else {
			jslog.Errorf("not found %v", v)
		}

	case turnresulttype.DropMoney:
		app.systemMessage.Appendf("Drop %v", makeMoneyColor(int(v.Arg)))

	case turnresulttype.DropMoneyInsteadOfCarryObj:
		o, exist := app.CaObjUUID2CaObjClient[v.DstUUID]
		if exist {
			app.systemMessage.Appendf("Drop %v for %v",
				makeMoneyColor(int(v.Arg)), obj2ColorStr(o),
			)
		} else {
			jslog.Errorf("not found %v", v)
		}

	case turnresulttype.AttackTo:
		dstao, exist := app.AOUUID2AOClient[v.DstUUID]
		aostr := "??"
		if exist {
			aostr = obj2ColorStr(dstao)
		}
		app.systemMessage.Appendf("%s %v",
			wrapspan.ColorTextf("Lime", "Damage  %4.1f to", v.Arg),
			aostr)
	case turnresulttype.Kill:
		dstao, exist := app.AOUUID2AOClient[v.DstUUID]
		aostr := "??"
		nickname := "??"
		if exist {
			aostr = obj2ColorStr(dstao)
			nickname = dstao.NickName
		}
		app.systemMessage.Appendf("Kill %v", aostr)
		app.NotiMessage.AppendTf(tcsInfo, "You kill %v", nickname)
	case turnresulttype.AttackedFrom:
		if o := app.currentFloor().FieldObjPosMan.GetByUUID(v.DstUUID); o != nil {
			fo := o.(*c2t_obj.FieldObjClient)
			app.systemMessage.Appendf("Damage %4.1f from %v", v.Arg, fo.Message)
			app.NotiMessage.AppendTf(tcsInfo, "Damage %4.1f from %v", v.Arg, fo.Message)
			return
		}
		dstao, exist := app.AOUUID2AOClient[v.DstUUID]
		aostr := "??"
		if exist {
			aostr = obj2ColorStr(dstao)
		}
		app.systemMessage.Appendf("%s %v",
			wrapspan.ColorTextf("Red", "Damage %4.1f from", v.Arg),
			aostr)
	case turnresulttype.KilledBy:
		if o := app.currentFloor().FieldObjPosMan.GetByUUID(v.DstUUID); o != nil {
			fo := o.(*c2t_obj.FieldObjClient)
			app.systemMessage.Appendf("Killed by %v", fo.Message)
			app.NotiMessage.AppendTf(tcsInfo, "Killed by %v", fo.Message)
			return
		}
		dstao, exist := app.AOUUID2AOClient[v.DstUUID]
		aostr := "??"
		nickname := "??"
		if exist {
			aostr = obj2ColorStr(dstao)
			nickname = dstao.NickName
		}
		app.systemMessage.Appendf("Killed by %v", aostr)
		app.NotiMessage.AppendTf(tcsInfo, "Killed by %v", nickname)

	case turnresulttype.HPDamageFromTrap:
		app.systemMessage.Appendf("HP Damage from Trap %4.1f", v.Arg)

	case turnresulttype.SPDamageFromTrap:
		app.systemMessage.Appendf("SP Damage from Trap %4.1f", v.Arg)

	case turnresulttype.DamagedByTile:

	case turnresulttype.DeadByTile:
		app.systemMessage.Append("Die by Tile damage")
		app.NotiMessage.AppendTf(tcsInfo, "Die by Tile damage")

	case turnresulttype.ContagionFrom:
		dstao, exist := app.AOUUID2AOClient[v.DstUUID]
		aostr := "??"
		nickname := "??"
		if exist {
			aostr = obj2ColorStr(dstao)
			nickname = dstao.NickName
		}
		app.systemMessage.Appendf("Contagion from %v", aostr)
		app.NotiMessage.AppendTf(tcsInfo, "Contagion from %v", nickname)
	case turnresulttype.ContagionTo:
		dstao, exist := app.AOUUID2AOClient[v.DstUUID]
		aostr := "??"
		nickname := "??"
		if exist {
			aostr = obj2ColorStr(dstao)
			nickname = dstao.NickName
		}
		app.systemMessage.Appendf("Contagion to %v", aostr)
		app.NotiMessage.AppendTf(tcsInfo, "Contagion to %v", nickname)
	case turnresulttype.ContagionFromFail:
		dstao, exist := app.AOUUID2AOClient[v.DstUUID]
		aostr := "??"
		nickname := "??"
		if exist {
			aostr = obj2ColorStr(dstao)
			nickname = dstao.NickName
		}
		app.systemMessage.Appendf("Resist Contagion from %v", aostr)
		app.NotiMessage.AppendTf(tcsInfo, "Resist Contagion from %v", nickname)
	case turnresulttype.ContagionToFail:
		dstao, exist := app.AOUUID2AOClient[v.DstUUID]
		aostr := "??"
		nickname := "??"
		if exist {
			aostr = obj2ColorStr(dstao)
			nickname = dstao.NickName
		}
		app.systemMessage.Appendf("Fail Contagion to %v", aostr)
		app.NotiMessage.AppendTf(tcsInfo, "Fail Contagion to %v", nickname)
	}
}

// from noti obj list
func (app *WasmClient) actByControlMode() {
	if app.olNotiData == nil || app.olNotiData.ActiveObj.HP <= 0 {
		return
	}

	switch gameOptions.GetByIDBase("ViewMode").State {
	case 0: // play viewpot mode
		if app.moveByUserInput() {
			return
		}
	case 1: // floor viewport mode
	}

	if fo := app.onFieldObj; fo == nil ||
		(fo != nil && fieldobjacttype.ClientData[fo.ActType].ActOn) {
		for i, v := range autoActs.ButtonList {
			if v.State == 0 {
				if tryAutoActFn[i](app, v) {
					return
				}
			}
		}
	}
}

func (app *WasmClient) moveByUserInput() bool {
	cf := app.currentFloor()
	playerX, playerY := app.GetPlayerXY()
	if !cf.IsValidPos(playerX, playerY) {
		jslog.Errorf("ao out of floor %v %v", app.olNotiData.ActiveObj, cf)
		return false
	}
	w, h := cf.Tiles.GetXYLen()
	switch app.ClientColtrolMode {
	default:
		jslog.Errorf("invalid ClientColtrolMode %v", app.ClientColtrolMode)
	case clientcontroltype.Keyboard:
		if app.sendMovePacketByInput(app.KeyDir) {
			return true
		}
	case clientcontroltype.FollowMouse:
		if app.sendMovePacketByInput(app.MouseDir) {
			return true
		}
	case clientcontroltype.MoveToDest:
		playerPos := [2]int{playerX, playerY}
		if app.Path2dst == nil || len(app.Path2dst) == 0 {
			app.ClientColtrolMode = clientcontroltype.Keyboard
			return false
		}
		for i := len(app.Path2dst) - 1; i >= 0; i-- {
			nextPos := app.Path2dst[i]
			isContact, dir := way9type.CalcContactDirWrapped(playerPos, nextPos, w, h)
			if isContact {
				if dir == way9type.Center {
					// arrived
					app.Path2dst = nil
					app.vp.ClearMovePath()
					app.ClientColtrolMode = clientcontroltype.Keyboard
					return false
				} else {
					go app.sendPacket(c2t_idcmd.Move,
						&c2t_obj.ReqMove_data{Dir: dir},
					)
					return true
				}
			}
		}
		// fail to path2dest
		app.Path2dst = nil
		app.vp.ClearMovePath()
		app.ClientColtrolMode = clientcontroltype.Keyboard
	}
	return false
}

func SoundByActResult(ar *aoactreqrsp.ActReqRsp) {
	if ar != nil && ar.IsSuccess() {
		soundmap.PlayByAct(ar.Done.Act)
	}
}

func objRecvNotiFn_VPTiles(recvobj interface{}, header c2t_packet.Header, obj interface{}) error {
	robj, ok := obj.(*c2t_obj.NotiVPTiles_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", obj)
	}
	app, ok := recvobj.(*WasmClient)
	if !ok {
		return fmt.Errorf("recvobj type mismatch %v", recvobj)
	}

	if !app.waitObjList {
		app.waitObjList = true
	}

	app.taNotiHeader = header
	app.taNotiData = robj

	if app.FloorInfo == nil {
		jslog.Warn("OrangeRed", "app.FloorInfo not set")
		return nil
	}
	if app.FloorInfo.Name != app.taNotiData.FloorName {
		jslog.Warnf("not current floor vptile data %v %v",
			app.currentFloor().FloorInfo.Name, app.taNotiData.FloorName,
		)
		return nil
	}
	cf, exist := app.Name2ClientFloor[app.taNotiData.FloorName]
	if !exist {
		jslog.Warnf("floor not added %v", app.taNotiData.FloorName)
		return nil
	}

	oldComplete := cf.Visited.IsComplete()
	if err := cf.UpdateFromViewportTile(
		app.taNotiData, gInitData.ViewportXYLenList); err != nil {
		jslog.Warn("%v", err)
		return nil
	}
	app.vp.ProcessNotiVPTiles(
		cf, app.taNotiData, app.olNotiData, app.Path2dst)

	if !oldComplete && cf.Visited.IsComplete() { // just completed
		app.systemMessage.Append(wrapspan.ColorTextf("yellow",
			"Discover floor complete %v", cf.FloorInfo.Name))
		app.NotiMessage.AppendTf(tcsInfo,
			"Discover floor complete %v", cf.FloorInfo.Name)
	}
	return nil
}

func objRecvNotiFn_FloorTiles(recvobj interface{}, header c2t_packet.Header, obj interface{}) error {
	robj, ok := obj.(*c2t_obj.NotiFloorTiles_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", obj)
	}
	app, ok := recvobj.(*WasmClient)
	if !ok {
		return fmt.Errorf("recvobj type mismatch %v", recvobj)
	}
	cf, exist := app.Name2ClientFloor[robj.FI.Name]
	if !exist {
		// new floor
		cf = clientfloor.New(robj.FI)
		app.Name2ClientFloor[robj.FI.Name] = cf
		app.systemMessage.Append(wrapspan.ColorTextf("yellow",
			"Found floor %v", cf.FloorInfo.Name))
		app.NotiMessage.AppendTf(tcsInfo,
			"Found floor %v", cf.FloorInfo.Name)
	}

	oldComplete := cf.Visited.IsComplete()
	cf.ReplaceFloorTiles(robj)
	if !oldComplete && cf.Visited.IsComplete() {
		app.systemMessage.Append(wrapspan.ColorTextf("yellow",
			"Discover floor %v complete", cf.FloorInfo.Name))
		app.NotiMessage.AppendTf(tcsInfo,
			"Discover floor %v complete", cf.FloorInfo.Name)
	}
	return nil
}

// FieldObjList    // for rebuild known floor
func objRecvNotiFn_FieldObjList(recvobj interface{}, hd c2t_packet.Header, body interface{}) error {
	robj, ok := body.(*c2t_obj.NotiFieldObjList_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", body)
	}
	app, ok := recvobj.(*WasmClient)
	if !ok {
		return fmt.Errorf("recvobj type mismatch %v", recvobj)
	}
	cf, exist := app.Name2ClientFloor[robj.FI.Name]
	if !exist {
		// new floor
		cf = clientfloor.New(robj.FI)
		app.Name2ClientFloor[robj.FI.Name] = cf
		app.systemMessage.Append(wrapspan.ColorTextf("yellow",
			"Found floor %v", cf.FloorInfo.Name))
		app.NotiMessage.AppendTf(tcsInfo,
			"Found floor %v", cf.FloorInfo.Name)
	}
	cf.UpdateFieldObjList(robj.FOList)
	return nil
}

func objRecvNotiFn_FoundFieldObj(recvobj interface{}, header c2t_packet.Header, obj interface{}) error {
	robj, ok := obj.(*c2t_obj.NotiFoundFieldObj_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", obj)
	}
	app, ok := recvobj.(*WasmClient)
	if !ok {
		return fmt.Errorf("recvobj type mismatch %v", recvobj)
	}
	_ = app
	_ = robj

	fromFloor, exist := app.Name2ClientFloor[robj.FloorName]
	if !exist {
		jslog.Warnf("FoundFieldObj unknonw floor %v", robj)
	}
	if fromFloor.FieldObjPosMan.Get1stObjAt(robj.FieldObj.X, robj.FieldObj.Y) == nil {
		fromFloor.FieldObjPosMan.AddToXY(robj.FieldObj, robj.FieldObj.X, robj.FieldObj.Y)
		// fromFloor.addFieldObj(robj.FieldObj)
		notiString := "Found Hidden FieldObj"
		app.systemMessage.Append(wrapspan.ColorText("yellow",
			notiString))
		app.NotiMessage.AppendTf(tcsInfo,
			notiString)
	}
	return nil
}

func objRecvNotiFn_ForgetFloor(recvobj interface{}, header c2t_packet.Header, obj interface{}) error {
	robj, ok := obj.(*c2t_obj.NotiForgetFloor_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", obj)
	}
	app, ok := recvobj.(*WasmClient)
	if !ok {
		return fmt.Errorf("recvobj type mismatch %v", recvobj)
	}
	_ = app
	_ = robj

	forgetFloor, exist := app.Name2ClientFloor[robj.FloorName]
	if !exist {
		jslog.Warnf("ForgetFloor unknonw floor %v", robj)
	} else {
		forgetFloor.Forget()
	}
	return nil
}

func objRecvNotiFn_ActivateTrap(recvobj interface{}, header c2t_packet.Header, obj interface{}) error {
	robj, ok := obj.(*c2t_obj.NotiActivateTrap_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", obj)
	}
	app, ok := recvobj.(*WasmClient)
	if !ok {
		return fmt.Errorf("recvobj type mismatch %v", recvobj)
	}
	_ = app
	_ = robj

	if robj.Triggered {
		app.systemMessage.Append(wrapspan.ColorTextf("OrangeRed",
			"Activate %v", robj.FieldObjAct.String()))
		app.NotiMessage.AppendTf(tcsWarn,
			"Activate %v", robj.FieldObjAct.String())
	} else {
		app.systemMessage.Append(wrapspan.ColorTextf("Yellow",
			"Evade %v", robj.FieldObjAct.String()))
		app.NotiMessage.AppendTf(tcsInfo,
			"Evade %v", robj.FieldObjAct.String())
	}
	return nil
}
