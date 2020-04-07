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
	"github.com/kasworld/goguelike/game/gamei"
	"github.com/kasworld/goguelike/lib/g2id"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
	"github.com/kasworld/uuidstr"
)

type Money struct {
	uuid              string
	remainTurnInFloor int

	value float64
}

func (m Money) String() string {
	return fmt.Sprintf("Money[%v $%v]", m.uuid, m.GetValue())
}

func NewMoney(v float64) gamei.MoneyI {
	if v < 0 {
		v = 0
	}
	rtn := Money{
		uuid:  uuidstr.New(),
		value: v,
	}
	return &rtn
}

func (po Money) GetValue() float64 {
	return po.value
}

func (po Money) Add(po2 gamei.MoneyI) gamei.MoneyI {
	v := po.GetValue() + po2.GetValue()
	return NewMoney(v)
}
func (po Money) Sub(po2 gamei.MoneyI) gamei.MoneyI {
	v := po.GetValue() - po2.GetValue()
	return NewMoney(v)
}

func (po *Money) ToPacket_CarryObjClientOnFloor(x, y int) *c2t_obj.CarryObjClientOnFloor {
	poc := &c2t_obj.CarryObjClientOnFloor{
		G2ID:               g2id.NewFromString(po.uuid),
		CarryingObjectType: po.GetCarryingObjectType(),
		X:                  x,
		Y:                  y,

		Value: int(po.value),
	}
	return poc
}

// IDPosI interface
func (po *Money) GetUUID() string {
	return po.uuid
}

func (po *Money) GetCarryingObjectType() carryingobjecttype.CarryingObjectType {
	return carryingobjecttype.Money
}
func (po *Money) GetWeight() float64 {
	return float64(po.value) * gameconst.MoneyGram
}

// life in floor handle

func (po *Money) GetRemainTurnInFloor() int {
	return po.remainTurnInFloor
}
func (po *Money) DecRemainTurnInFloor() int {
	if po.remainTurnInFloor > 0 {
		po.remainTurnInFloor--
	}
	return po.remainTurnInFloor
}
func (po *Money) SetRemainTurnInFloor() {
	po.remainTurnInFloor = gameconst.CarryingObjectLifeTurnInFloor
}
