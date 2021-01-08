// Copyright 2014,2015,2016,2017,2018,2019,2020,2021 SeukWon Kang (kasworld@gmail.com)
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package inventory

import (
	"fmt"

	"github.com/kasworld/goguelike/config/gameconst"
	"github.com/kasworld/goguelike/enum/equipslottype"
	"github.com/kasworld/goguelike/enum/towerachieve"
	"github.com/kasworld/goguelike/game/bias"
	"github.com/kasworld/goguelike/game/carryingobject"
	"github.com/kasworld/goguelike/game/gamei"
)

var _ gamei.InventoryI = &Inventory{}

func (inv *Inventory) GetEquipSlot() [equipslottype.EquipSlotType_Count]gamei.EquipObjI {
	return inv.equipSlot
}

func (inv *Inventory) SumEquipAttackBias() bias.Bias {
	rtn := bias.Bias{}
	for i, po := range inv.GetEquipSlot() {
		if po == nil {
			continue
		}
		if equipslottype.EquipSlotType(i).Attack() {
			rtn = rtn.Add(po.GetBias())
		}
	}
	return rtn
}
func (inv *Inventory) SumEquipDefenceBias() bias.Bias {
	rtn := bias.Bias{}
	for i, po := range inv.GetEquipSlot() {
		if po == nil {
			continue
		}
		if equipslottype.EquipSlotType(i).Defence() {
			rtn = rtn.Add(po.GetBias())
		}
	}
	return rtn
}

func (inv *Inventory) GetEquipedCount() int {
	eqCount := 0
	for _, v := range inv.equipSlot {
		if v != nil {
			eqCount++
		}
	}
	return eqCount
}

func (inv *Inventory) GetBagCount() int {
	return len(inv.bag)
}

func (inv *Inventory) GetTypeCount() (int, int, int) {
	var eqCount, potionCount, scrollCount int
	inv.mutexBag.RLock()
	for _, v := range inv.bag {
		switch v.(type) {
		default:
		case gamei.EquipObjI:
			eqCount++
		case gamei.PotionI:
			potionCount++
		case gamei.ScrollI:
			scrollCount++
		}
	}
	inv.mutexBag.RUnlock()
	return eqCount, potionCount, scrollCount
}

func (inv *Inventory) GetTypeList() (
	[]gamei.EquipObjI, []gamei.PotionI, []gamei.ScrollI) {
	rtne := make([]gamei.EquipObjI, 0)
	rtnp := make([]gamei.PotionI, 0)
	rtns := make([]gamei.ScrollI, 0)
	inv.mutexBag.RLock()
	for _, v := range inv.bag {
		switch e := v.(type) {
		default:
		case gamei.EquipObjI:
			rtne = append(rtne, e)
		case gamei.PotionI:
			rtnp = append(rtnp, e)
		case gamei.ScrollI:
			rtns = append(rtns, e)
		}
	}
	inv.mutexBag.RUnlock()
	return rtne, rtnp, rtns
}

func (inv *Inventory) GetEquipList() []gamei.EquipObjI {
	rtn := make([]gamei.EquipObjI, 0)
	inv.mutexBag.RLock()
	for _, v := range inv.bag {
		if e, ok := v.(gamei.EquipObjI); ok {
			rtn = append(rtn, e)
		}
	}
	inv.mutexBag.RUnlock()
	return rtn
}

func (inv *Inventory) GetPotionList() []gamei.PotionI {
	rtn := make([]gamei.PotionI, 0)
	inv.mutexBag.RLock()
	for _, v := range inv.bag {
		if e, ok := v.(gamei.PotionI); ok {
			rtn = append(rtn, e)
		}
	}
	inv.mutexBag.RUnlock()
	return rtn
}

func (inv *Inventory) GetScrollList() []gamei.ScrollI {
	rtn := make([]gamei.ScrollI, 0)
	inv.mutexBag.RLock()
	for _, v := range inv.bag {
		if e, ok := v.(gamei.ScrollI); ok {
			rtn = append(rtn, e)
		}
	}
	inv.mutexBag.RUnlock()
	return rtn
}

func (inv *Inventory) GetTotalWeight() float64 {
	rtn := float64(inv.wallet)*gameconst.MoneyGram +
		inv.poTotalWeight
	return rtn
}

func (inv *Inventory) GetTotalValue() float64 {
	rtn := float64(inv.wallet) + inv.poTotalValue
	return rtn
}

func (inv *Inventory) AddToBag(po gamei.CarryingObjectI) error {
	inv.mutexBag.Lock()
	if _, exist := inv.bag[po.GetUUID()]; exist {
		return fmt.Errorf("already owned %v", po)
	}
	inv.bag[po.GetUUID()] = po
	inv.mutexBag.Unlock()
	inv.poTotalWeight += po.GetWeight()
	inv.poTotalValue += po.GetValue()
	switch po.(type) {
	default:
		return fmt.Errorf("unknown obj")
	case gamei.EquipObjI:
		inv.towerAchieveStat.Inc(towerachieve.EquipIn)
	case gamei.PotionI:
		inv.towerAchieveStat.Inc(towerachieve.PotionIn)
	case gamei.ScrollI:
		inv.towerAchieveStat.Inc(towerachieve.ScrollIn)
	}
	return nil
}

func (inv *Inventory) GetByUUID(poid string) gamei.CarryingObjectI {
	if po, err := inv.getFromEquipByUUID(poid); err == nil {
		return po
	}
	inv.mutexBag.RLock()
	defer inv.mutexBag.RUnlock()
	return inv.bag[poid]
}

func (inv *Inventory) RecycleCarryObjByID(poid string) (float64, error) {
	po := inv.RemoveByUUID(poid)
	if po == nil {
		return 0, fmt.Errorf("not in inventory %v", poid)
	}
	recycleValue := po.GetValue() * gameconst.CarryObjRecycleRate
	inv.AddToWallet(carryingobject.NewMoney(recycleValue))
	return recycleValue, nil
}
func (inv *Inventory) AddToWallet(po gamei.MoneyI) error {
	inv.wallet += po.GetValue()
	inv.towerAchieveStat.Add(towerachieve.MoneyIn, float64(po.GetValue()))
	return nil
}
func (inv *Inventory) SubFromWallet(po gamei.MoneyI) error {
	if inv.wallet < po.GetValue() {
		return fmt.Errorf("insufficient money %v %v", inv.wallet, po)
	}
	inv.wallet -= po.GetValue()
	inv.towerAchieveStat.Add(towerachieve.MoneyOut, float64(po.GetValue()))
	return nil
}
func (inv *Inventory) GetWalletValue() float64 {
	return inv.wallet
}

func (inv *Inventory) EquipFromBagByUUID(id string) error {
	inv.mutexBag.Lock()
	defer inv.mutexBag.Unlock()
	o, exist := inv.bag[id]
	if !exist {
		return fmt.Errorf("not in bag %v", id)
	}
	delete(inv.bag, id)
	po, ok := o.(gamei.EquipObjI)
	if !ok {
		return fmt.Errorf("invalid type %v", o)
	}
	i := po.GetEquipType()
	old := inv.equipSlot[i]
	inv.equipSlot[i] = po
	if old != nil {
		inv.bag[old.GetUUID()] = old
	}
	return nil
}

func (inv *Inventory) UnEquipToBagByUUID(id string) (gamei.EquipObjI, error) {
	inv.mutexBag.Lock()
	defer inv.mutexBag.Unlock()
	for eqPos, po := range inv.equipSlot {
		if po != nil && po.GetUUID() == id {
			inv.equipSlot[eqPos] = nil
			inv.bag[po.GetUUID()] = po
			return po, nil
		}
	}
	return nil, fmt.Errorf("not in equipSlot %v", id)
}

func (inv *Inventory) RemoveByUUID(poid string) gamei.CarryingObjectI {
	// if equiped unequip
	inv.UnEquipToBagByUUID(poid)
	inv.mutexBag.Lock()
	po, exist := inv.bag[poid]
	if !exist {
		inv.mutexBag.Unlock()
		return nil
	}
	delete(inv.bag, poid)
	inv.mutexBag.Unlock()
	inv.poTotalWeight -= po.GetWeight()
	inv.poTotalValue -= po.GetValue()
	switch po.(type) {
	default:
		fmt.Printf("unknown obj %v", po)
		return po
	case gamei.EquipObjI:
		inv.towerAchieveStat.Inc(towerachieve.EquipOut)
	case gamei.PotionI:
		inv.towerAchieveStat.Inc(towerachieve.PotionOut)
	case gamei.ScrollI:
		inv.towerAchieveStat.Inc(towerachieve.ScrollOut)
	}
	return po
}
