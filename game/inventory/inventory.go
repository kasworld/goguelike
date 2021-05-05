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
	"sync"

	"github.com/kasworld/goguelike/enum/equipslottype"
	"github.com/kasworld/goguelike/enum/towerachieve_vector_float64"
	"github.com/kasworld/goguelike/game/gamei"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
)

// equip, bag is exclusive
type Inventory struct {
	towerAchieveStat *towerachieve_vector_float64.TowerAchieveVector_float64 `prettystring:"simple"`
	equipSlot        [equipslottype.EquipSlotType_Count]gamei.EquipObjI
	mutexBag         sync.RWMutex `prettystring:"hide"`
	bag              map[string]gamei.CarryingObjectI
	wallet           float64
	poTotalWeight    float64
	poTotalValue     float64
}

func New(towerAchieveStat *towerachieve_vector_float64.TowerAchieveVector_float64) *Inventory {
	return &Inventory{
		bag:              make(map[string]gamei.CarryingObjectI),
		wallet:           0,
		towerAchieveStat: towerAchieveStat,
	}
}

func (inv *Inventory) String() string {
	return fmt.Sprintf(
		"Inventory[equip:%v bag:%v wallet:%v weight:%v]",
		inv.GetEquipedCount(), len(inv.bag), inv.wallet, inv.GetTotalWeight())
}

func (inv *Inventory) TotalCarryObjCount() int {
	return inv.GetEquipedCount() + len(inv.bag)
}

func (inv *Inventory) getFromEquipByUUID(id string) (gamei.EquipObjI, error) {
	for _, v := range inv.equipSlot {
		if v != nil && v.GetUUID() == id {
			return v, nil
		}
	}
	return nil, fmt.Errorf("not in equipSlot %v", id)
}

func (inv *Inventory) ToPacket_EquipClient() []*c2t_obj.EquipClient {
	var EquippedPo []*c2t_obj.EquipClient
	for _, v := range inv.equipSlot {
		if v == nil {
			continue
		}
		EquippedPo = append(EquippedPo, v.ToPacket_EquipClient())
	}
	return EquippedPo
}

func (inv *Inventory) ToPacket_InvenInfos() (
	[]*c2t_obj.EquipClient,
	[]*c2t_obj.EquipClient,
	[]*c2t_obj.PotionClient,
	[]*c2t_obj.ScrollClient,
	int,
) {
	var EquippedPo []*c2t_obj.EquipClient
	var equipBag []*c2t_obj.EquipClient
	var potionBag []*c2t_obj.PotionClient
	var scrollBag []*c2t_obj.ScrollClient
	for _, v := range inv.equipSlot {
		if v == nil {
			continue
		}
		EquippedPo = append(EquippedPo, v.ToPacket_EquipClient())
	}
	inv.mutexBag.RLock()
	for _, v := range inv.bag {
		switch o := v.(type) {
		default:
		case gamei.EquipObjI:
			equipBag = append(equipBag, o.ToPacket_EquipClient())
		case gamei.PotionI:
			potionBag = append(potionBag, o.ToPacket_PotionClient())
		case gamei.ScrollI:
			scrollBag = append(scrollBag, o.ToPacket_ScrollClient())
		}
	}
	inv.mutexBag.RUnlock()

	return EquippedPo, equipBag, potionBag, scrollBag, int(inv.wallet)
}
