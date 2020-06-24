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

package carryingobject

import (
	"fmt"

	"github.com/kasworld/goguelike/config/gameconst"
	"github.com/kasworld/goguelike/enum/carryingobjecttype"
	"github.com/kasworld/goguelike/enum/potiontype"
	"github.com/kasworld/goguelike/game/gamei"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
	"github.com/kasworld/uuidstr"
)

type Potion struct {
	uuid              string
	remainTurnInFloor int

	potionType potiontype.PotionType
}

func (p Potion) String() string {
	return fmt.Sprintf("Potion[%v %v]",
		p.uuid, p.potionType.String())
}

func NewPotion(pt potiontype.PotionType) gamei.PotionI {
	rtn := &Potion{
		uuid:       uuidstr.New(),
		potionType: pt,
	}
	return rtn
}

func NewPotionByMakeRate(n int) gamei.PotionI {
	for i := 0; i < potiontype.PotionType_Count; i++ {
		pot := potiontype.PotionType(i)
		n -= pot.MakeRate()
		if n < 0 {
			return NewPotion(pot)
		}
	}
	return NewPotion(potiontype.Empty)
}

func (po *Potion) ToPacket_CarryObjClientOnFloor(x, y int) *c2t_obj.CarryObjClientOnFloor {
	poc := &c2t_obj.CarryObjClientOnFloor{
		UUID:               po.uuid,
		CarryingObjectType: po.GetCarryingObjectType(),
		X:                  x,
		Y:                  y,

		PotionType: po.potionType,
	}
	return poc
}

func (po *Potion) ToPacket_PotionClient() *c2t_obj.PotionClient {
	poc := &c2t_obj.PotionClient{
		UUID:       po.uuid,
		PotionType: po.potionType,
	}
	return poc
}

// IDPosI interface
func (po *Potion) GetUUID() string {
	return po.uuid
}

func (po *Potion) GetName() string {
	return po.potionType.String()
}

func (po *Potion) GetPotionType() potiontype.PotionType {
	return po.potionType
}

func (po *Potion) GetCarryingObjectType() carryingobjecttype.CarryingObjectType {
	return carryingobjecttype.Potion
}

func (po *Potion) GetValue() float64 {
	return gameconst.PotionValue
}

func (po *Potion) GetWeight() float64 {
	return gameconst.PotionGram
}

// life in floor handle

func (po *Potion) GetRemainTurnInFloor() int {
	return po.remainTurnInFloor
}
func (po *Potion) DecRemainTurnInFloor() int {
	if po.remainTurnInFloor > 0 {
		po.remainTurnInFloor--
	}
	return po.remainTurnInFloor
}
func (po *Potion) SetRemainTurnInFloor() {
	po.remainTurnInFloor = gameconst.CarryingObjectLifeTurnInFloor
}
