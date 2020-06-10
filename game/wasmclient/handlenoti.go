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

package wasmclient

import (
	"fmt"
	"sync/atomic"
	"syscall/js"
	"time"

	"github.com/kasworld/goguelike/enum/clientcontroltype"
	"github.com/kasworld/goguelike/enum/condition"
	"github.com/kasworld/goguelike/enum/fieldobjacttype"
	"github.com/kasworld/goguelike/enum/way9type"
	"github.com/kasworld/goguelike/lib/g2id"

	"github.com/kasworld/goguelike/config/gameconst"
	"github.com/kasworld/goguelike/config/leveldata"
	"github.com/kasworld/goguelike/enum/turnresulttype"
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
	c2t_idnoti.EnterTower:     objRecvNotiFn_EnterTower,
	c2t_idnoti.LeaveTower:     objRecvNotiFn_LeaveTower,
	c2t_idnoti.EnterFloor:     objRecvNotiFn_EnterFloor,
	c2t_idnoti.LeaveFloor:     objRecvNotiFn_LeaveFloor,
	c2t_idnoti.Ageing:         objRecvNotiFn_Ageing,
	c2t_idnoti.Death:          objRecvNotiFn_Death,
	c2t_idnoti.ReadyToRebirth: objRecvNotiFn_ReadyToRebirth,
	c2t_idnoti.Rebirthed:      objRecvNotiFn_Rebirthed,
	c2t_idnoti.Broadcast:      objRecvNotiFn_Broadcast,
	c2t_idnoti.VPTiles:        objRecvNotiFn_VPTiles,
	c2t_idnoti.ObjectList:     objRecvNotiFn_ObjectList,
	c2t_idnoti.FloorTiles:     objRecvNotiFn_FloorTiles,
	c2t_idnoti.FoundFieldObj:  objRecvNotiFn_FoundFieldObj,
	c2t_idnoti.ForgetFloor:    objRecvNotiFn_ForgetFloor,
	c2t_idnoti.ActivateTrap:   objRecvNotiFn_ActivateTrap,
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
	gVP2d.NotiMessage.AppendTf(tcsInfo,
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
	gVP2d.NotiMessage.AppendTf(tcsInfo,
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
	cf, exist := app.G2ID2ClientFloor[robj.FI.G2ID]
	if !exist {
		// new floor
		cf = clientfloor.New(robj.FI)
		app.G2ID2ClientFloor[robj.FI.G2ID] = cf
		app.systemMessage.Append(wrapspan.ColorTextf("yellow",
			"Found floor %v", cf.FloorInfo.Name))
		gVP2d.NotiMessage.AppendTf(tcsInfo,
			"Found floor %v", cf.FloorInfo.Name)
	}
	app.systemMessage.Appendf("Enter floor %v", cf.FloorInfo.Name)
	cf.EnterFloor()
	gVP2d.NotiMessage.AppendTf(tcsInfo,
		"Enter floor %v", cf.FloorInfo.Name)
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
	gVP2d.NotiMessage.AppendTf(tcsInfo,
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
	gVP2d.NotiMessage.AppendTf(tcsInfo,
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
	gVP2d.NotiMessage.AppendTf(tcsWarn,
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
	gVP2d.NotiMessage.AppendTf(tcsInfo,
		"Ready to rebirth")

	if autoActs.GetByIDBase("AutoRebirth").State == 0 {
		go app.sendPacket(c2t_idcmd.Rebirth,
			&c2t_obj.ReqRebirth_data{},
		)
	} else {
		commandButtons.GetByIDBase("Rebirth").Enable()
	}
	// clear carryobj cache
	app.CaObjG2ID2CaObjClient = make(map[g2id.G2ID]interface{})
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
	gVP2d.NotiMessage.AppendTf(tcsInfo,
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
	gVP2d.NotiMessage.AppendTf(tcsInfo,
		"Broadcast %v", robj.Msg)
	return nil
}

func objRecvNotiFn_ObjectList(recvobj interface{}, header c2t_packet.Header, obj interface{}) error {
	robj, ok := obj.(*c2t_obj.NotiObjectList_data)
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
			js.Global().Call("requestAnimationFrame", js.FuncOf(app.drawCanvas))
		}
	}()
	app.ServerClientTimeDiff = robj.Time.Sub(time.Now())
	app.olNotiHeader = header
	if app.olNotiData == nil {
		app.lastOLNotiData = robj
	} else {
		app.lastOLNotiData = app.olNotiData
	}
	app.olNotiData = robj
	newOLNotiData := robj
	app.onFieldObj = nil

	app.ClientJitter.Act()
	app.ServerJitter.ActByValue(robj.Time)

	c2t_obj.EquipClientByG2ID(robj.ActiveObj.EquipBag).Sort()
	c2t_obj.PotionClientByG2ID(robj.ActiveObj.PotionBag).Sort()
	c2t_obj.ScrollClientByG2ID(robj.ActiveObj.ScrollBag).Sort()

	oldLevel := 0
	exp := 0
	if app.lastOLNotiData != nil {
		oldLevel = int(leveldata.CalcLevelFromExp(float64(app.lastOLNotiData.ActiveObj.Exp)))
		app.HPdiff = newOLNotiData.ActiveObj.HP - app.lastOLNotiData.ActiveObj.HP
		app.SPdiff = newOLNotiData.ActiveObj.SP - app.lastOLNotiData.ActiveObj.SP
		exp = newOLNotiData.ActiveObj.Exp - app.lastOLNotiData.ActiveObj.Exp

		if app.lastOLNotiData.ActiveObj.Bias.NearFaction() != robj.ActiveObj.Bias.NearFaction() {
			app.systemMessage.Appendf(
				"Faction changed to %v", robj.ActiveObj.Bias.NearFaction().String())
			gVP2d.NotiMessage.AppendTf(tcsInfo,
				"Faction changed to %v", robj.ActiveObj.Bias.NearFaction().String())

		}
		if app.lastOLNotiData.ActiveObj.Bias != robj.ActiveObj.Bias {
			app.systemMessage.Append("Bias changed")
		}

	}
	newLevel := int(leveldata.CalcLevelFromExp(float64(newOLNotiData.ActiveObj.Exp)))
	if oldLevel < newLevel {
		soundmap.Play("levelupsound")
		app.systemMessage.Append(wrapspan.ColorTextf("yellow",
			"Level Up to %v", newLevel))
		gVP2d.NotiMessage.AppendTf(tcsInfo,
			"Level Up to %v", newLevel)
	} else if oldLevel > newLevel {
		soundmap.Play("leveldownsound")
		app.systemMessage.Append(wrapspan.ColorTextf("yellow",
			"Level Down to %v", newLevel))
		gVP2d.NotiMessage.AppendTf(tcsInfo,
			"Level Down to %v", newLevel)
	}
	if exp != 0 {
		gVP2d.ActiveObjMessage.Appendf("yellow", calcSizeRateByValue(exp), time.Second,
			"Exp %+d", exp)
	}
	if app.HPdiff > 0 {
		gVP2d.ActiveObjMessage.Appendf("LimeGreen", calcSizeRateByValue(app.HPdiff), time.Second,
			"HP %+d", app.HPdiff)
	} else if app.HPdiff < 0 {
		gVP2d.ActiveObjMessage.Appendf("Red", calcSizeRateByValue(-app.HPdiff), time.Second,
			"HP %+d", app.HPdiff)
	}
	if app.SPdiff > 0 {
		gVP2d.ActiveObjMessage.Appendf("LimeGreen", calcSizeRateByValue(app.SPdiff), time.Second,
			"SP %+d", app.SPdiff)
	} else if app.SPdiff < 0 {
		gVP2d.ActiveObjMessage.Appendf("Red", calcSizeRateByValue(-app.SPdiff), time.Second,
			"SP %+d", app.SPdiff)
	}

	for i := 0; i < condition.Condition_Count; i++ {
		ci := condition.Condition(i)
		if newOLNotiData.ActiveObj.Conditions.TestByCondition(ci) {
			gVP2d.ActiveObjMessage.Appendf("OrangeRed", 1.0, time.Second,
				"%s", ci)
		}
	}

	for _, v := range newOLNotiData.ActiveObj.EquipBag {
		app.CaObjG2ID2CaObjClient[v.G2ID] = v
	}
	for _, v := range newOLNotiData.ActiveObj.EquippedPo {
		app.CaObjG2ID2CaObjClient[v.G2ID] = v
	}
	for _, v := range newOLNotiData.ActiveObj.PotionBag {
		app.CaObjG2ID2CaObjClient[v.G2ID] = v
	}
	for _, v := range newOLNotiData.ActiveObj.ScrollBag {
		app.CaObjG2ID2CaObjClient[v.G2ID] = v
	}

	for _, v := range newOLNotiData.ActiveObjList {
		app.AOG2ID2AOClient[v.G2ID] = v
	}
	for _, v := range newOLNotiData.CarryObjList {
		if _, exist := app.CaObjG2ID2CaObjClient[v.G2ID]; !exist {
			app.CaObjG2ID2CaObjClient[v.G2ID] = v
		}
	}

	if actStr := app.ActionResult2String(newOLNotiData.ActiveObj.Act); actStr != "" {
		app.systemMessage.Append(actStr)
	}

	for _, v := range newOLNotiData.ActiveObj.TurnResult {
		switch v.ResultType {
		default:
			jslog.Error("unknown acttype %v", v)

		case turnresulttype.DropCarryObj:
			o, exist := app.CaObjG2ID2CaObjClient[v.DstG2ID]
			if exist {
				app.systemMessage.Appendf("Drop %v", obj2ColorStr(o))
			} else {
				jslog.Errorf("not found %v", v)
			}

		case turnresulttype.DropMoney:
			app.systemMessage.Appendf("Drop %v", makeMoneyColor(int(v.Arg)))

		case turnresulttype.DropMoneyInsteadOfCarryObj:
			o, exist := app.CaObjG2ID2CaObjClient[v.DstG2ID]
			if exist {
				app.systemMessage.Appendf("Drop %v for %v",
					makeMoneyColor(int(v.Arg)), obj2ColorStr(o),
				)
			} else {
				jslog.Errorf("not found %v", v)
			}

		case turnresulttype.AttackTo:
			dstao, exist := app.AOG2ID2AOClient[v.DstG2ID]
			aostr := ""
			if exist {
				aostr = obj2ColorStr(dstao)
			}
			app.systemMessage.Appendf("%s %v",
				wrapspan.ColorTextf("Lime", "Damage  %4.1f to", v.Arg),
				aostr)
		case turnresulttype.Kill:
			dstao, exist := app.AOG2ID2AOClient[v.DstG2ID]
			aostr := ""
			nickname := ""
			if exist {
				aostr = obj2ColorStr(dstao)
				nickname = dstao.NickName
			}
			app.systemMessage.Appendf("Kill %v", aostr)
			gVP2d.NotiMessage.AppendTf(tcsInfo, "You kill %v", nickname)
		case turnresulttype.AttackedFrom:
			dstao, exist := app.AOG2ID2AOClient[v.DstG2ID]
			aostr := ""
			if exist {
				aostr = obj2ColorStr(dstao)
			}
			app.systemMessage.Appendf("%s %v",
				wrapspan.ColorTextf("Red", "Damage %4.1f from", v.Arg),
				aostr)
		case turnresulttype.KilledBy:
			dstao, exist := app.AOG2ID2AOClient[v.DstG2ID]
			aostr := ""
			nickname := ""
			if exist {
				aostr = obj2ColorStr(dstao)
				nickname = dstao.NickName
			}
			app.systemMessage.Appendf("Killed by %v", aostr)
			gVP2d.NotiMessage.AppendTf(tcsInfo, "You are killed by %v", nickname)

		case turnresulttype.HPDamageFromTrap:
			app.systemMessage.Appendf("HP Damage from Trap %4.1f", v.Arg)

		case turnresulttype.SPDamageFromTrap:
			app.systemMessage.Appendf("SP Damage from Trap %4.1f", v.Arg)

		case turnresulttype.DamagedByTile:

		case turnresulttype.DeadByTile:
			app.systemMessage.Append("Die by Tile damage")
			gVP2d.NotiMessage.AppendTf(tcsInfo, "Die by Tile damage")
		}
	}

	SoundByActResult(newOLNotiData.ActiveObj.Act)
	wLimit := leveldata.WeightLimit(newLevel)
	if newOLNotiData.ActiveObj.Conditions.TestByCondition(condition.Burden) {
		wLimit /= 2
	}
	app.OverLoadRate = newOLNotiData.ActiveObj.CalcWeight() / wLimit
	if app.OverLoadRate > 1 {
		soundmap.Play("overloadsound")
		gVP2d.ActiveObjMessage.Appendf("red", app.OverLoadRate, time.Second,
			"Overloaded")
	}
	if pao := app.GetPlayerActiveObjClient(); pao != nil && pao.DamageTake > 1 {
		soundmap.Play("damagesound")
	}

	if newOLNotiData.ActiveObj.HP <= 0 {
		if app.remainTurn2Rebirth > 0 {
			gVP2d.NotiMessage.AppendTf(tcsRebirth, "Remain to rebirth %v", app.remainTurn2Rebirth)
		} else {
			gVP2d.NotiMessage.AppendTf(tcsRebirth, "Ready to rebirth %v", -app.remainTurn2Rebirth)
		}
		app.remainTurn2Rebirth--
	}

	if app.FloorInfo == nil {
		jslog.Error("app.FloorInfo not set")
		return nil
	}
	if app.FloorInfo.G2ID != newOLNotiData.FloorG2ID {
		jslog.Errorf("not current floor objlist data %v %v",
			app.currentFloor().FloorInfo.G2ID, newOLNotiData.FloorG2ID,
		)
		return nil
	}

	cf, exist := app.G2ID2ClientFloor[newOLNotiData.FloorG2ID]
	if !exist {
		jslog.Warnf("floor not added %v", newOLNotiData.FloorG2ID)
		return nil
	}
	for _, v := range newOLNotiData.FieldObjList {
		cf.FieldObjPosMan.AddOrUpdateToXY(v, v.X, v.Y)
	}
	playerX, playerY := app.GetPlayerXY()
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
	if newOLNotiData.ActiveObj.RemainTurn2Act <= 0 {
		app.actByControlMode()
	}

	return nil
}

// from noti obj list
func (app *WasmClient) actByControlMode() {
	if app.olNotiData == nil || app.olNotiData.ActiveObj.HP <= 0 {
		return
	}

	switch gameOptions.GetByIDBase("Viewport").State {
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
	if app.FloorInfo.G2ID != app.taNotiData.FloorG2ID {
		jslog.Warnf("not current floor vptile data %v %v",
			app.currentFloor().FloorInfo.G2ID, app.taNotiData.FloorG2ID,
		)
		return nil
	}
	cf, exist := app.G2ID2ClientFloor[app.taNotiData.FloorG2ID]
	if !exist {
		jslog.Warnf("floor not added %v", app.taNotiData.FloorG2ID)
		return nil
	}

	oldComplete := cf.Visited.IsComplete()
	if err := cf.UpdateFromViewportTile(app.taNotiData, gInitData.ViewportXYLenList); err != nil {
		jslog.Warn("%v", err)
		return nil
	}
	if !oldComplete && cf.Visited.IsComplete() { // just completed
		app.systemMessage.Append(wrapspan.ColorTextf("yellow",
			"Discover floor complete %v", cf.FloorInfo.Name))
		gVP2d.NotiMessage.AppendTf(tcsInfo,
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
	cf, exist := app.G2ID2ClientFloor[robj.FI.G2ID]
	if !exist {
		// new floor
		cf = clientfloor.New(robj.FI)
		app.G2ID2ClientFloor[robj.FI.G2ID] = cf
		app.systemMessage.Append(wrapspan.ColorTextf("yellow",
			"Found floor %v", cf.FloorInfo.Name))
		gVP2d.NotiMessage.AppendTf(tcsInfo,
			"Found floor %v", cf.FloorInfo.Name)
	}

	oldComplete := cf.Visited.IsComplete()
	cf.ReplaceFloorTiles(robj)
	if !oldComplete && cf.Visited.IsComplete() {
		app.systemMessage.Append(wrapspan.ColorTextf("yellow",
			"Discover floor %v complete", cf.FloorInfo.Name))
		gVP2d.NotiMessage.AppendTf(tcsInfo,
			"Discover floor %v complete", cf.FloorInfo.Name)
	}
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

	fromFloor, exist := app.G2ID2ClientFloor[robj.FloorG2ID]
	if !exist {
		jslog.Warnf("FoundFieldObj unknonw floor %v", robj)
	}
	if fromFloor.FieldObjPosMan.Get1stObjAt(robj.FieldObj.X, robj.FieldObj.Y) == nil {
		fromFloor.FieldObjPosMan.AddOrUpdateToXY(robj.FieldObj, robj.FieldObj.X, robj.FieldObj.Y)
		notiString := "Found Hidden FieldObj"
		app.systemMessage.Append(wrapspan.ColorText("yellow",
			notiString))
		gVP2d.NotiMessage.AppendTf(tcsInfo,
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

	forgetFloor, exist := app.G2ID2ClientFloor[robj.FloorG2ID]
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
		gVP2d.NotiMessage.AppendTf(tcsWarn,
			"Activate %v", robj.FieldObjAct.String())
	} else {
		app.systemMessage.Append(wrapspan.ColorTextf("Yellow",
			"Evade %v", robj.FieldObjAct.String()))
		gVP2d.NotiMessage.AppendTf(tcsInfo,
			"Evade %v", robj.FieldObjAct.String())
	}
	return nil
}
