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

package gamei

import (
	"github.com/kasworld/goguelike/enum/carryingobjecttype"
	"github.com/kasworld/goguelike/enum/equipslottype"
	"github.com/kasworld/goguelike/enum/potiontype"
	"github.com/kasworld/goguelike/enum/scrolltype"
	"github.com/kasworld/goguelike/game/bias"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
)

type CarryingObjectI interface {
	ToPacket_CarryObjClientOnFloor(x, y int) *c2t_obj.CarryObjClientOnFloor
	GetRemainTurnInFloor() int
	DecRemainTurnInFloor() int
	SetRemainTurnInFloor()

	GetUUID() string
	GetCarryingObjectType() carryingobjecttype.CarryingObjectType
	GetWeight() float64
	GetValue() float64
}

type EquipObjI interface {
	CarryingObjectI
	ToPacket_EquipClient() *c2t_obj.EquipClient

	// bias, faction
	GetEquipType() equipslottype.EquipSlotType
	GetBias() bias.Bias
}

type PotionI interface {
	CarryingObjectI
	ToPacket_PotionClient() *c2t_obj.PotionClient
	GetPotionType() potiontype.PotionType
}

type MoneyI interface {
	CarryingObjectI
	Add(po2 MoneyI) MoneyI
	Sub(po2 MoneyI) MoneyI
}

type ScrollI interface {
	CarryingObjectI
	GetScrollType() scrolltype.ScrollType
	ToPacket_ScrollClient() *c2t_obj.ScrollClient
}
