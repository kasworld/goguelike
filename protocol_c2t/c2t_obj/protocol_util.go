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

package c2t_obj

import (
	"fmt"
	"sort"

	"github.com/kasworld/goguelike/config/gameconst"
	"github.com/kasworld/goguelike/enum/carryingobjecttype"
	"github.com/kasworld/goguelike/enum/turnresulttype"
	"github.com/kasworld/goguelike/game/bias"
)

func (po CarryObjClientOnFloor) String() string {
	switch po.CarryingObjectType {
	default:
		return fmt.Sprintf("unknonw po %#v", po)
	case carryingobjecttype.Equip:
		return fmt.Sprintf("%v[%v]",
			po.EquipType.String(),
			po.Faction.String(),
		)
	case carryingobjecttype.Money:
		return fmt.Sprintf("$%v", po.Value)
	case carryingobjecttype.Potion:
		return fmt.Sprintf("Potion%v", po.PotionType.String())
	case carryingobjecttype.Scroll:
		return fmt.Sprintf("Scroll%v", po.ScrollType.String())
	}
}

func (pao PlayerActiveObjInfo) CalcWeight() float64 {
	weight := 0.0
	for _, v := range pao.EquippedPo {
		weight += v.Weight()
	}
	for _, v := range pao.EquipBag {
		weight += v.Weight()
	}
	weight += float64(len(pao.PotionBag)) * gameconst.PotionGram
	weight += float64(len(pao.ScrollBag)) * gameconst.ScrollGram
	weight += float64(pao.Wallet) * gameconst.MoneyGram
	return weight
}

func (pao PlayerActiveObjInfo) CalcDamageGive() float64 {
	var DamageGive float64
	for _, v := range pao.TurnResult {
		switch v.ResultType {
		case turnresulttype.AttackTo:
			DamageGive += v.Arg
		}
	}
	return DamageGive
}
func (pao PlayerActiveObjInfo) CalcDamageTake() float64 {
	var DamageTake float64
	for _, v := range pao.TurnResult {
		switch v.ResultType {
		case turnresulttype.AttackedFrom:
			DamageTake += v.Arg
		case turnresulttype.DamagedByTile:
			DamageTake += v.Arg
		}
	}
	return DamageTake
}

func (pao PlayerActiveObjInfo) Exist(id string) bool {
	for _, v := range pao.EquippedPo {
		if v == nil {
			continue
		}
		if v.UUID == id {
			return true
		}
	}
	for _, v := range pao.EquipBag {
		if v == nil {
			continue
		}
		if v.UUID == id {
			return true
		}
	}
	for _, v := range pao.PotionBag {
		if v == nil {
			continue
		}
		if v.UUID == id {
			return true
		}
	}
	for _, v := range pao.ScrollBag {
		if v == nil {
			continue
		}
		if v.UUID == id {
			return true
		}
	}
	return false
}

func (ifa TurnResultClient) String() string {
	return fmt.Sprintf("%v DstObj:%v Arg:%3.1f",
		ifa.ResultType.String(),
		ifa.DstUUID,
		ifa.Arg,
	)
}

func (po EquipClient) String() string {
	return fmt.Sprintf("%v %v%.0f",
		po.Name,
		po.Faction.Rune(),
		po.BiasLen,
	)
}

func (po EquipClient) GetBias() bias.Bias {
	return bias.NewByFaction(po.Faction, po.BiasLen)
}

func (po EquipClient) Weight() float64 {
	return po.BiasLen * gameconst.EquipABSGram
}

type EquipClientByUUID []*EquipClient

func (objList EquipClientByUUID) Len() int { return len(objList) }
func (objList EquipClientByUUID) Swap(i, j int) {
	objList[i], objList[j] = objList[j], objList[i]
}
func (objList EquipClientByUUID) Less(i, j int) bool {
	po1 := objList[i]
	po2 := objList[j]
	if po1.EquipType == po2.EquipType {
		return po1.UUID < po2.UUID
	} else {
		return po1.EquipType < po2.EquipType
	}
}
func (objList EquipClientByUUID) Sort() {
	sort.Sort(objList)
}

func (po PotionClient) String() string {
	return fmt.Sprintf("Potion%v", po.PotionType.String())
}

func (po PotionClient) Weight() int {
	return gameconst.PotionGram
}

type PotionClientByUUID []*PotionClient

func (objList PotionClientByUUID) Len() int { return len(objList) }
func (objList PotionClientByUUID) Swap(i, j int) {
	objList[i], objList[j] = objList[j], objList[i]
}
func (objList PotionClientByUUID) Less(i, j int) bool {
	po1 := objList[i]
	po2 := objList[j]
	return po1.UUID < po2.UUID
}
func (objList PotionClientByUUID) Sort() {
	sort.Sort(objList)
}

func (po ScrollClient) String() string {
	return fmt.Sprintf("Scroll%v", po.ScrollType.String())
}

func (po ScrollClient) Weight() int {
	return gameconst.ScrollGram
}

type ScrollClientByUUID []*ScrollClient

func (objList ScrollClientByUUID) Len() int { return len(objList) }
func (objList ScrollClientByUUID) Swap(i, j int) {
	objList[i], objList[j] = objList[j], objList[i]
}
func (objList ScrollClientByUUID) Less(i, j int) bool {
	po1 := objList[i]
	po2 := objList[j]
	return po1.UUID < po2.UUID
}
func (objList ScrollClientByUUID) Sort() {
	sort.Sort(objList)
}

type CarryObjEqByLen []*EquipClient

func (objList CarryObjEqByLen) Len() int { return len(objList) }
func (objList CarryObjEqByLen) Swap(i, j int) {
	objList[i], objList[j] = objList[j], objList[i]
}
func (objList CarryObjEqByLen) Less(i, j int) bool {
	po1 := objList[i]
	po2 := objList[j]
	return po1.BiasLen < po2.BiasLen
}
func (objList CarryObjEqByLen) Sort() {
	sort.Stable(objList)
}

type FieldObjByType []*FieldObjClient

func (objList FieldObjByType) Len() int { return len(objList) }
func (objList FieldObjByType) Swap(i, j int) {
	objList[i], objList[j] = objList[j], objList[i]
}
func (objList FieldObjByType) Less(i, j int) bool {
	po1 := objList[i]
	po2 := objList[j]
	if po1.ActType == po2.ActType {
		return po1.ID < po2.ID
	}
	return po1.ActType < po2.ActType
}
func (objList FieldObjByType) Sort() {
	sort.Stable(objList)
}

type FloorInfoByName []*FloorInfo

func (objList FloorInfoByName) Len() int { return len(objList) }
func (objList FloorInfoByName) Swap(i, j int) {
	objList[i], objList[j] = objList[j], objList[i]
}
func (objList FloorInfoByName) Less(i, j int) bool {
	po1 := objList[i]
	po2 := objList[j]
	return po1.Name > po2.Name
}
func (objList FloorInfoByName) Sort() {
	sort.Stable(objList)
}
