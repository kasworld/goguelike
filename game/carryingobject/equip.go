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

package carryingobject

import (
	"fmt"

	"github.com/kasworld/g2rand"
	"github.com/kasworld/goguelike/config/gameconst"
	"github.com/kasworld/goguelike/enum/carryingobjecttype"
	"github.com/kasworld/goguelike/enum/equipslottype"
	"github.com/kasworld/goguelike/enum/factiontype"
	"github.com/kasworld/goguelike/game/bias"
	"github.com/kasworld/goguelike/game/gamei"
	"github.com/kasworld/goguelike/lib/idu64str"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
)

var EquipIDMaker = idu64str.New("EquipID")

func (po EquipObj) String() string {
	return fmt.Sprintf("EquipObj[%v %v %v %v]",
		po.uuid, po.equipType, po.Faction, po.BiasLen)
}

type EquipObj struct {
	uuid              string
	remainTurnInFloor int

	equipType equipslottype.EquipSlotType
	name      string

	Faction factiontype.FactionType
	BiasLen float64
}

func NewRandFactionEquipObj(aoname string, ft factiontype.FactionType, rnd *g2rand.G2Rand) gamei.EquipObjI {
	po := EquipObj{
		uuid: EquipIDMaker.New(),
	}
	po.Faction = ft
	po.equipType = equipslottype.EquipSlotType(rnd.Intn(equipslottype.EquipSlotType_Count))

	// materialadj := po.Faction.String()

	materiallist := po.equipType.Materials()
	material := materiallist[rnd.Intn(len(materiallist))]

	namelist := po.equipType.Names()
	name := namelist[rnd.Intn(len(namelist))]
	po.name = fmt.Sprintf("%s's %s %s", aoname, material, name)

	biaslen := po.equipType.BaseLen()
	po.BiasLen = rnd.NormFloat64Range(biaslen, biaslen/2)
	if po.BiasLen < 0 {
		po.BiasLen = -po.BiasLen
	}

	return &po
}

func NewEquipByFactionSlot(aoname string,
	ft factiontype.FactionType,
	eqslot equipslottype.EquipSlotType,
	rnd *g2rand.G2Rand) gamei.EquipObjI {

	po := EquipObj{
		uuid: EquipIDMaker.New(),
	}
	po.Faction = ft
	po.equipType = eqslot

	// materialadj := po.Faction.String()

	materiallist := po.equipType.Materials()
	material := materiallist[rnd.Intn(len(materiallist))]

	namelist := po.equipType.Names()
	name := namelist[rnd.Intn(len(namelist))]
	po.name = fmt.Sprintf("%s's %s %s", aoname, material, name)

	biaslen := po.equipType.BaseLen()
	po.BiasLen = rnd.NormFloat64Range(biaslen, biaslen/2)
	if po.BiasLen < 0 {
		po.BiasLen = -po.BiasLen
	}

	return &po
}

func (po *EquipObj) GetBias() bias.Bias {
	return bias.NewByFaction(po.Faction, po.BiasLen)
}

func (po *EquipObj) ToPacket_CarryObjClientOnFloor(x, y int) *c2t_obj.CarryObjClientOnFloor {
	poc := &c2t_obj.CarryObjClientOnFloor{
		UUID:               po.uuid,
		CarryingObjectType: po.GetCarryingObjectType(),
		EquipType:          po.equipType,
		X:                  x,
		Y:                  y,

		Faction: po.Faction,
	}
	return poc
}

func (po *EquipObj) ToPacket_EquipClient() *c2t_obj.EquipClient {
	poc := &c2t_obj.EquipClient{
		UUID:      po.uuid,
		Name:      po.name,
		EquipType: po.equipType,

		Faction: po.Faction,
		BiasLen: po.BiasLen,
	}
	return poc
}

// IDPosI interface
func (po *EquipObj) GetUUID() string {
	return po.uuid
}
func (po *EquipObj) GetCarryingObjectType() carryingobjecttype.CarryingObjectType {
	return carryingobjecttype.Equip
}

func (po *EquipObj) GetEquipType() equipslottype.EquipSlotType {
	return po.equipType
}

func (po *EquipObj) GetValue() float64 {
	return po.BiasLen * gameconst.EquipABSValue
}

func (po *EquipObj) GetWeight() float64 {
	return po.BiasLen * gameconst.EquipABSGram
}

// life in floor handle

func (po *EquipObj) GetRemainTurnInFloor() int {
	return po.remainTurnInFloor
}
func (po *EquipObj) DecRemainTurnInFloor() int {
	if po.remainTurnInFloor > 0 {
		po.remainTurnInFloor--
	}
	return po.remainTurnInFloor
}
func (po *EquipObj) SetRemainTurnInFloor() {
	po.remainTurnInFloor = gameconst.CarryingObjectLifeTurnInFloor
}
