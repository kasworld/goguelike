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
	"github.com/kasworld/goguelike/config/leveldata"
	"github.com/kasworld/goguelike/enum/condition"
	"github.com/kasworld/goguelike/enum/equipslottype"
	"github.com/kasworld/goguelike/enum/fieldobjacttype"
	"github.com/kasworld/goguelike/enum/potiontype"
	"github.com/kasworld/goguelike/enum/scrolltype"
	"github.com/kasworld/goguelike/enum/way9type"
	"github.com/kasworld/goguelike/game/bias"
	"github.com/kasworld/goguelike/lib/htmlbutton"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_idcmd"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
	"github.com/kasworld/gowasmlib/jslog"
)

var autoActs = htmlbutton.NewButtonGroup("AutoActs",
	[]*htmlbutton.HTMLButton{
		// not client ai
		htmlbutton.New("z", "AutoPlay", []string{"AutoPlay", "NoAutoPlay"},
			"ServerAI on/off", cmdToggleServerAI, 0),
		htmlbutton.New("x", "AutoRebirth", []string{"AutoRebirth", "NoAutoRebirth"},
			"auto rebirth on/off", cmdToggleAutoRebirth, 0),
		// client ai
		htmlbutton.New("c", "AutoBattle", []string{"AutoBattle", "NoAutoBattle"},
			"Auto Battle on/off", cmdToggleAutoBattle, 0),
		htmlbutton.New("v", "AutoPickup", []string{"AutoPickup", "NoAutoPickup"},
			"Auto Pickup on/off", cmdToggleAutoPickup, 0),
		htmlbutton.New("b", "AutoEquip", []string{"AutoEquip", "NoAutoEquip"},
			"Auto Equip/Unequip on/off", cmdToggleAutoEquip, 0),
		htmlbutton.New("n", "AutoUsePotionScroll", []string{"AutoUsePotionScroll", "NoAutoUsePotionScroll"},
			"Auto Use Potion and Scroll on/off", cmdToggleAutoUsePotionScroll, 0),
		htmlbutton.New("m", "AutoRecyclePotionScroll", []string{"AutoRecyclePotionScroll", "NoAutoRecyclePotionScroll"},
			"Auto Recycle Potion and Scroll on/off", cmdToggleAutoRecyclePotionScroll, 0),
		htmlbutton.New(",", "AutoRecycleEquip", []string{"AutoRecycleEquip", "NoAutoRecycleEquip"},
			"Auto Recycle CarryObj on/off", cmdToggleAutoRecycleEquip, 0),
	})

func cmdToggleAutoRebirth(obj interface{}, v *htmlbutton.HTMLButton) {
	v.Blur()
}

func cmdToggleAutoBattle(obj interface{}, v *htmlbutton.HTMLButton) {
	v.Blur()
}

func cmdToggleAutoPickup(obj interface{}, v *htmlbutton.HTMLButton) {
	v.Blur()
}

func cmdToggleAutoEquip(obj interface{}, v *htmlbutton.HTMLButton) {
	v.Blur()
}
func cmdToggleAutoUsePotionScroll(obj interface{}, v *htmlbutton.HTMLButton) {
	v.Blur()
}
func cmdToggleAutoRecyclePotionScroll(obj interface{}, v *htmlbutton.HTMLButton) {
	v.Blur()
}
func cmdToggleAutoRecycleEquip(obj interface{}, v *htmlbutton.HTMLButton) {
	v.Blur()
}

var tryAutoActFn = []func(app *WasmClient, v *htmlbutton.HTMLButton) bool{
	tryAutoPlay,
	tryAutoRebirth,
	tryAutoBattle,
	tryAutoPickup,
	tryAutoEquip,
	tryAutoUsePotionScroll,
	tryAutoRecyclePotionScroll,
	tryAutoRecycleEquip,
}

func tryAutoPlay(app *WasmClient, v *htmlbutton.HTMLButton) bool {
	// not ai fn
	return false
}
func cmdToggleServerAI(obj interface{}, v *htmlbutton.HTMLButton) {
	app, ok := obj.(*WasmClient)
	if !ok {
		jslog.Errorf("obj not app %v", obj)
		return
	}
	// jslog.Infof("autoplay %v", v)
	switch v.State {
	case 0: // Auto
		go app.reqAIPlay(true)
		app.systemMessage.Append("AutoPlayOn")
		// app.vp.NotiMessage.AppendTf(tcsInfo, "AutoPlayOn")
	case 1: // NoAuto
		go app.reqAIPlay(false)
		app.systemMessage.Append("AutoPlayOff")
		// app.vp.NotiMessage.AppendTf(tcsInfo, "AutoPlayOff")
	}
	v.Blur()
}

func tryAutoRebirth(app *WasmClient, v *htmlbutton.HTMLButton) bool {
	// not ai fn
	return false
}

func tryAutoBattle(app *WasmClient, v *htmlbutton.HTMLButton) bool {
	if v.State != 0 {
		return false
	}
	if app.olNotiData == nil {
		return false
	}
	cf := app.currentFloor()
	if app.olNotiData.FloorUUID != cf.FloorInfo.UUID {
		return false
	}
	playerX, playerY := app.GetPlayerXY()
	if !cf.IsValidPos(playerX, playerY) {
		return false
	}
	if !cf.Tiles[playerX][playerY].CanBattle() {
		return false
	}
	w, h := cf.Tiles.GetXYLen()
	for _, ao := range app.olNotiData.ActiveObjList {
		if !ao.Alive {
			continue
		}
		if ao.UUID == gInitData.AccountInfo.ActiveObjUUID {
			continue
		}
		if !cf.Tiles[ao.X][ao.Y].CanBattle() {
			continue
		}
		isContact, dir := way9type.CalcContactDirWrappedXY(
			playerX, playerY, ao.X, ao.Y, w, h)
		if isContact && dir != way9type.Center {
			go app.sendPacket(c2t_idcmd.Attack,
				&c2t_obj.ReqAttack_data{Dir: dir},
			)
			return true
		}
	}
	return false
}

func tryAutoPickup(app *WasmClient, v *htmlbutton.HTMLButton) bool {
	if v.State != 0 {
		return false
	}
	if app.olNotiData.ActiveObj.Conditions.TestByCondition(condition.Float) {
		return false
	}
	cf := app.currentFloor()
	if app.olNotiData == nil {
		return false
	}
	playerX, playerY := app.GetPlayerXY()
	w, h := cf.Tiles.GetXYLen()
	for _, po := range app.olNotiData.CarryObjList {
		isContact, dir := way9type.CalcContactDirWrappedXY(
			playerX, playerY, po.X, po.Y, w, h)
		if !isContact {
			continue
		}
		o := cf.FieldObjPosMan.Get1stObjAt(po.X, po.Y)
		if o != nil {
			fo := o.(*c2t_obj.FieldObjClient)
			if !fieldobjacttype.ClientData[fo.ActType].ActOn {
				continue
			}
		}

		if dir == way9type.Center {
			go app.sendPacket(c2t_idcmd.Pickup,
				&c2t_obj.ReqPickup_data{UUID: po.UUID},
			)
			return true
		} else {
			go app.sendPacket(c2t_idcmd.Move,
				&c2t_obj.ReqMove_data{Dir: dir},
			)
			return true
		}
	}
	return false
}

func tryAutoEquip(app *WasmClient, v *htmlbutton.HTMLButton) bool {
	if v.State != 0 {
		return false
	}
	if app.olNotiData == nil {
		return false
	}
	for _, po := range app.olNotiData.ActiveObj.EquippedPo {
		if app.needUnEquipCarryObj(po.GetBias()) {
			go app.sendPacket(c2t_idcmd.UnEquip,
				&c2t_obj.ReqUnEquip_data{UUID: po.UUID},
			)
			return true
		}
	}
	for _, po := range app.olNotiData.ActiveObj.EquipBag {
		if app.isBetterCarryObj(po.EquipType, po.GetBias()) {
			go app.sendPacket(c2t_idcmd.Equip,
				&c2t_obj.ReqEquip_data{UUID: po.UUID},
			)
			return true
		}
	}
	return false
}

func tryAutoUsePotionScroll(app *WasmClient, v *htmlbutton.HTMLButton) bool {
	if v.State != 0 {
		return false
	}
	if app.olNotiData == nil {
		return false
	}
	for _, po := range app.olNotiData.ActiveObj.PotionBag {
		if app.needUsePotion(po) {
			go app.sendPacket(c2t_idcmd.DrinkPotion,
				&c2t_obj.ReqDrinkPotion_data{UUID: po.UUID},
			)
			return true
		}
	}

	for _, po := range app.olNotiData.ActiveObj.ScrollBag {
		if app.needUseScroll(po) {
			go app.sendPacket(c2t_idcmd.ReadScroll,
				&c2t_obj.ReqReadScroll_data{UUID: po.UUID},
			)
			return true
		}
	}

	return false
}

func tryAutoRecycleEquip(app *WasmClient, v *htmlbutton.HTMLButton) bool {
	if v.State != 0 {
		return false
	}
	if app.olNotiData == nil {
		return false
	}
	if app.olNotiData.ActiveObj.Conditions.TestByCondition(condition.Float) {
		return false
	}
	if app.onFieldObj == nil {
		return false
	}
	if app.onFieldObj.ActType != fieldobjacttype.RecycleCarryObj {
		return false
	}
	return app.recycleEqbag()
}

func tryAutoRecyclePotionScroll(app *WasmClient, v *htmlbutton.HTMLButton) bool {
	if v.State != 0 {
		return false
	}
	if app.olNotiData == nil {
		return false
	}
	if app.olNotiData.ActiveObj.Conditions.TestByCondition(condition.Float) {
		return false
	}
	if app.onFieldObj == nil {
		return false
	}
	if app.onFieldObj.ActType != fieldobjacttype.RecycleCarryObj {
		return false
	}
	if app.recycleUselessPotion() {
		return true
	}
	if app.recycleUselessScroll() {
		return true
	}
	return false
}

/////////

func (app *WasmClient) isBetterCarryObj(EquipType equipslottype.EquipSlotType, PoBias bias.Bias) bool {
	aoEnvBias := app.TowerBias().Add(app.currentFloor().GetBias()).Add(app.olNotiData.ActiveObj.Bias)
	newBiasAbs := aoEnvBias.Add(PoBias).AbsSum()
	for _, v := range app.olNotiData.ActiveObj.EquippedPo {
		if v.EquipType == EquipType {
			return newBiasAbs > aoEnvBias.Add(v.GetBias()).AbsSum()+1
		}
	}
	return newBiasAbs > aoEnvBias.AbsSum()+1
}

func (app *WasmClient) needUnEquipCarryObj(PoBias bias.Bias) bool {
	aoEnvBias := app.TowerBias().Add(app.currentFloor().GetBias()).Add(app.olNotiData.ActiveObj.Bias)

	currentBias := aoEnvBias.Add(PoBias)
	newBias := aoEnvBias
	return newBias.AbsSum() > currentBias.AbsSum()+1
}

func (app *WasmClient) needUseScroll(po *c2t_obj.ScrollClient) bool {
	cf := app.currentFloor()
	switch po.ScrollType {
	case scrolltype.FloorMap:
		if cf.Visited.CalcCompleteRate() < 1.0 {
			return true
		}
	}
	return false
}

func (app *WasmClient) needUsePotion(po *c2t_obj.PotionClient) bool {
	pao := app.olNotiData.ActiveObj
	switch po.PotionType {
	case potiontype.MinorHealing:
		return pao.HPMax-pao.HP > 10
	case potiontype.MajorHealing:
		return pao.HPMax-pao.HP > 50
	case potiontype.GreatHealing:
		return pao.HPMax-pao.HP > 100

	case potiontype.MinorActing:
		return pao.SPMax-pao.SP > 10
	case potiontype.MajorActing:
		return pao.SPMax-pao.SP > 50
	case potiontype.GreatActing:
		return pao.SPMax-pao.SP > 100

	case potiontype.MinorHeal:
		return pao.HPMax-pao.HP > pao.HPMax/10
	case potiontype.MajorHeal:
		return pao.HPMax/2 > pao.HP
	case potiontype.CompleteHeal:
		return pao.HPMax/10 > pao.HP

	case potiontype.MinorAct:
		return pao.SPMax-pao.SP > pao.SPMax/10
	case potiontype.MajorAct:
		return pao.SPMax/2 > pao.SP
	case potiontype.CompleteAct:
		return pao.SPMax/10 > pao.SP

	case potiontype.MinorSpanHealing:
		return pao.HPMax/2 > pao.HP
	case potiontype.MinorSpanActing:
		return pao.SPMax/2 > pao.SP

	case potiontype.MinorSpanVision:
		return pao.Sight <= leveldata.Sight(app.level)
	case potiontype.MajorSpanVision:
		return pao.Sight <= leveldata.Sight(app.level)
	case potiontype.PerfectSpanVision:
		return pao.Sight <= leveldata.Sight(app.level)
	}
	return false
}

func (app *WasmClient) recycleEqbag() bool {
	var poList c2t_obj.CarryObjEqByLen
	poList = append(poList, app.olNotiData.ActiveObj.EquipBag...)
	if len(poList) == 0 {
		return false
	}
	poList.Sort()
	go app.sendPacket(c2t_idcmd.Recycle,
		&c2t_obj.ReqRecycle_data{UUID: poList[0].UUID},
	)
	return true
}

func (app *WasmClient) recycleUselessPotion() bool {
	for _, po := range app.olNotiData.ActiveObj.PotionBag {
		if potiontype.AIRecycleMap[po.PotionType] {
			go app.sendPacket(c2t_idcmd.Recycle,
				&c2t_obj.ReqRecycle_data{UUID: po.UUID},
			)
			return true
		}
	}
	return false
}

func (app *WasmClient) recycleUselessScroll() bool {
	for _, po := range app.olNotiData.ActiveObj.ScrollBag {
		if scrolltype.AIRecycleMap[po.ScrollType] {
			go app.sendPacket(c2t_idcmd.Recycle,
				&c2t_obj.ReqRecycle_data{UUID: po.UUID},
			)
			return true
		}
	}
	return false
}
