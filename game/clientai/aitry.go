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
	"github.com/kasworld/goguelike/config/leveldata"
	"github.com/kasworld/goguelike/enum/condition"
	"github.com/kasworld/goguelike/enum/equipslottype"
	"github.com/kasworld/goguelike/enum/fieldobjacttype"
	"github.com/kasworld/goguelike/enum/potiontype"
	"github.com/kasworld/goguelike/enum/scrolltype"
	"github.com/kasworld/goguelike/enum/way9type"
	"github.com/kasworld/goguelike/game/bias"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_idcmd"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_packet"
)

var tryAutoActFn = []func(cai *ClientAI) bool{
	tryAutoBattle,
	tryAutoPickup,
	tryAutoEquip,
	tryAutoUsePotion,
	tryAutoRecyclePotionScroll,
	tryAutoRecycleEquip,
}

func (cai *ClientAI) actByControlMode() {
	if cai.OLNotiData == nil || cai.OLNotiData.ActiveObj.HP <= 0 {
		return
	}
	for _, v := range tryAutoActFn {
		if v(cai) {
			return
		}
	}
}

func tryAutoBattle(cai *ClientAI) bool {
	cf := cai.currentFloor()
	if cai.OLNotiData == nil {
		return false
	}
	playerX, playerY := cai.GetPlayerXY()
	if !cf.IsValidPos(playerX, playerY) {
		return false
	}
	if !cf.Tiles[playerX][playerY].CanBattle() {
		return false
	}
	w, h := cf.Tiles.GetXYLen()
	for _, ao := range cai.OLNotiData.ActiveObjList {
		if !ao.Alive {
			continue
		}
		if ao.G2ID == cai.AccountInfo.ActiveObjG2ID {
			continue
		}
		if !cf.IsValidPos(ao.X, ao.Y) {
			continue
		}
		if !cf.Tiles[ao.X][ao.Y].CanBattle() {
			continue
		}
		isContact, dir := way9type.CalcContactDirWrappedXY(
			playerX, playerY, ao.X, ao.Y, w, h)
		if isContact && dir != way9type.Center {
			cai.ReqWithRspFnWithAuth(c2t_idcmd.Attack,
				&c2t_obj.ReqAttack_data{Dir: dir},
				func(hd c2t_packet.Header, rsp interface{}) error {
					return nil
				})
			return true
		}
	}
	return false
}

func tryAutoPickup(cai *ClientAI) bool {
	cf := cai.currentFloor()
	if cai.OLNotiData == nil {
		return false
	}
	if cai.OLNotiData.ActiveObj.Conditions.TestByCondition(condition.Float) {
		return false
	}

	playerX, playerY := cai.GetPlayerXY()
	w, h := cf.Tiles.GetXYLen()
	for _, po := range cai.OLNotiData.CarryObjList {
		isContact, dir := way9type.CalcContactDirWrappedXY(
			playerX, playerY, po.X, po.Y, w, h)
		if !isContact {
			continue
		}
		if dir == way9type.Center {
			cai.ReqWithRspFnWithAuth(c2t_idcmd.Pickup,
				&c2t_obj.ReqPickup_data{G2ID: po.G2ID},
				func(hd c2t_packet.Header, rsp interface{}) error {
					return nil
				})
			return true
		} else {
			cai.ReqWithRspFnWithAuth(c2t_idcmd.Move,
				&c2t_obj.ReqMove_data{Dir: dir},
				func(hd c2t_packet.Header, rsp interface{}) error {
					return nil
				})
			return true
		}
	}
	return false
}

func tryAutoEquip(cai *ClientAI) bool {
	if cai.OLNotiData == nil {
		return false
	}
	for _, po := range cai.OLNotiData.ActiveObj.EquippedPo {
		if cai.needUnEquipCarryObj(po.GetBias()) {
			cai.ReqWithRspFnWithAuth(c2t_idcmd.UnEquip,
				&c2t_obj.ReqUnEquip_data{G2ID: po.G2ID},
				func(hd c2t_packet.Header, rsp interface{}) error {
					return nil
				})
			return true
		}
	}
	for _, po := range cai.OLNotiData.ActiveObj.EquipBag {
		if cai.isBetterCarryObj(po.EquipType, po.GetBias()) {
			cai.ReqWithRspFnWithAuth(c2t_idcmd.Equip,
				&c2t_obj.ReqEquip_data{G2ID: po.G2ID},
				func(hd c2t_packet.Header, rsp interface{}) error {
					return nil
				})
			return true
		}
	}
	return false
}

func tryAutoUsePotion(cai *ClientAI) bool {
	if cai.OLNotiData == nil {
		return false
	}
	for _, po := range cai.OLNotiData.ActiveObj.PotionBag {
		if cai.needUsePotion(po) {
			cai.ReqWithRspFnWithAuth(c2t_idcmd.DrinkPotion,
				&c2t_obj.ReqDrinkPotion_data{G2ID: po.G2ID},
				func(hd c2t_packet.Header, rsp interface{}) error {
					return nil
				})
			return true
		}
	}

	for _, po := range cai.OLNotiData.ActiveObj.ScrollBag {
		if cai.needUseScroll(po) {
			cai.ReqWithRspFnWithAuth(c2t_idcmd.ReadScroll,
				&c2t_obj.ReqReadScroll_data{G2ID: po.G2ID},
				func(hd c2t_packet.Header, rsp interface{}) error {
					return nil
				})
			return true
		}
	}

	return false
}

func tryAutoRecycleEquip(cai *ClientAI) bool {
	if cai.OLNotiData == nil {
		return false
	}
	if cai.OLNotiData.ActiveObj.Conditions.TestByCondition(condition.Float) {
		return false
	}
	if cai.onFieldObj == nil {
		return false
	}
	if cai.onFieldObj.ActType != fieldobjacttype.RecycleCarryObj {
		return false
	}
	return cai.recycleEqbag()
}

func tryAutoRecyclePotionScroll(cai *ClientAI) bool {
	if cai.OLNotiData == nil {
		return false
	}
	if cai.OLNotiData.ActiveObj.Conditions.TestByCondition(condition.Float) {
		return false
	}
	if cai.onFieldObj == nil {
		return false
	}
	if cai.onFieldObj.ActType != fieldobjacttype.RecycleCarryObj {
		return false
	}
	if cai.recycleUselessPotion() {
		return true
	}
	if cai.recycleUselessScroll() {
		return true
	}
	return false
}

/////////

func (cai *ClientAI) isBetterCarryObj(EquipType equipslottype.EquipSlotType, PoBias bias.Bias) bool {
	aoEnvBias := cai.TowerBias().Add(cai.currentFloor().GetBias()).Add(cai.OLNotiData.ActiveObj.Bias)
	newBiasAbs := aoEnvBias.Add(PoBias).AbsSum()
	for _, v := range cai.OLNotiData.ActiveObj.EquippedPo {
		if v.EquipType == EquipType {
			return newBiasAbs > aoEnvBias.Add(v.GetBias()).AbsSum()+1
		}
	}
	return newBiasAbs > aoEnvBias.AbsSum()+1
}

func (cai *ClientAI) needUnEquipCarryObj(PoBias bias.Bias) bool {
	aoEnvBias := cai.TowerBias().Add(cai.currentFloor().GetBias()).Add(cai.OLNotiData.ActiveObj.Bias)

	currentBias := aoEnvBias.Add(PoBias)
	newBias := aoEnvBias
	return newBias.AbsSum() > currentBias.AbsSum()+1
}

func (cai *ClientAI) needUseScroll(po *c2t_obj.ScrollClient) bool {
	cf := cai.currentFloor()
	switch po.ScrollType {
	case scrolltype.FloorMap:
		if cf.Visited.CalcCompleteRate() < 1.0 {
			return true
		}
	}
	return false
}

func (cai *ClientAI) needUsePotion(po *c2t_obj.PotionClient) bool {
	pao := cai.OLNotiData.ActiveObj
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
		return pao.Sight <= leveldata.Sight(cai.level)
	case potiontype.MajorSpanVision:
		return pao.Sight <= leveldata.Sight(cai.level)
	case potiontype.PerfectSpanVision:
		return pao.Sight <= leveldata.Sight(cai.level)
	}
	return false
}

func (cai *ClientAI) recycleEqbag() bool {
	var poList c2t_obj.CarryObjEqByLen
	poList = append(poList, cai.OLNotiData.ActiveObj.EquipBag...)
	if len(poList) == 0 {
		return false
	}
	poList.Sort()
	cai.ReqWithRspFnWithAuth(c2t_idcmd.Recycle,
		&c2t_obj.ReqRecycle_data{G2ID: poList[0].G2ID},
		func(hd c2t_packet.Header, rsp interface{}) error {
			return nil
		})
	return true
}

func (cai *ClientAI) recycleUselessPotion() bool {
	for _, po := range cai.OLNotiData.ActiveObj.PotionBag {
		if potiontype.AIRecycleMap[po.PotionType] {
			cai.ReqWithRspFnWithAuth(c2t_idcmd.Recycle,
				&c2t_obj.ReqRecycle_data{G2ID: po.G2ID},
				func(hd c2t_packet.Header, rsp interface{}) error {
					return nil
				})
			return true
		}
	}
	return false
}

func (cai *ClientAI) recycleUselessScroll() bool {
	for _, po := range cai.OLNotiData.ActiveObj.ScrollBag {
		if scrolltype.AIRecycleMap[po.ScrollType] {
			cai.ReqWithRspFnWithAuth(c2t_idcmd.Recycle,
				&c2t_obj.ReqRecycle_data{G2ID: po.G2ID},
				func(hd c2t_packet.Header, rsp interface{}) error {
					return nil
				})
			return true
		}
	}
	return false
}
