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

package dangerobject

import (
	"github.com/kasworld/goguelike/enum/dangertype"
	"github.com/kasworld/goguelike/lib/idu64str"
	"github.com/kasworld/goguelike/lib/uuidposman"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
)

var DOIDMaker = idu64str.New("DOID")

type DangerObject struct {
	UUID           string
	Owner          uuidposman.UUIDPosI // id of owner ( ao, floor? , fieldobj?)
	OwnerX, OwnerY int                 // pos at make attack time
	DangerType     dangertype.DangerType
	RemainTurn     int // remain turn to affect
	AffectRate     float64
}

// IDPosI interface
func (p *DangerObject) GetUUID() string {
	return p.UUID
}

func NewAOAttact(attacker uuidposman.UUIDPosI, dt dangertype.DangerType, srcx, srcy int) *DangerObject {
	return &DangerObject{
		UUID:       DOIDMaker.New(),
		Owner:      attacker,
		OwnerX:     srcx,
		OwnerY:     srcy,
		DangerType: dt,
		RemainTurn: dt.Turn2Live(),
		AffectRate: 1,
	}
}

func NewFOAttact(attacker uuidposman.UUIDPosI, dt dangertype.DangerType, affectRate float64) *DangerObject {
	return &DangerObject{
		UUID:       DOIDMaker.New(),
		Owner:      attacker,
		DangerType: dt,
		RemainTurn: dt.Turn2Live(),
		AffectRate: affectRate,
	}
}

func (p *DangerObject) ToPacket_DangerObjClient(x, y int) *c2t_obj.DangerObjClient {
	return &c2t_obj.DangerObjClient{
		UUID:       p.UUID,
		OwnerID:    p.Owner.GetUUID(),
		DangerType: p.DangerType,
		X:          x,
		Y:          y,
		AffectRate: p.AffectRate,
	}
}

// Live1Turn reduce remain turn and return alive
func (p *DangerObject) Live1Turn() bool {
	p.RemainTurn--
	return p.RemainTurn > 0
}
