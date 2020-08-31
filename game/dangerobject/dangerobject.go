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
	"github.com/kasworld/goguelike/lib/uuidposman"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
	"github.com/kasworld/uuidstr"
)

type DangerObject struct {
	UUID           string
	Owner          uuidposman.UUIDPosI // id of owner ( ao, floor? , fieldobj?)
	OwnerX, OwnerY int                 // pos at make attack time
	DangerType     dangertype.DangerType
	RemainTurn     int // remain turn to affect

}

// IDPosI interface
func (p *DangerObject) GetUUID() string {
	return p.UUID
}

func NewBasicAttact(attacker uuidposman.UUIDPosI, srcx, srcy int) *DangerObject {
	dt := dangertype.BasicAttack
	return &DangerObject{
		UUID:       uuidstr.New(),
		Owner:      attacker,
		OwnerX:     srcx,
		OwnerY:     srcy,
		DangerType: dt,
		RemainTurn: dt.Turn2Live(),
	}
}

func (p *DangerObject) ToPacket_DangerObjClient(x, y int) *c2t_obj.DangerObjClient {
	return &c2t_obj.DangerObjClient{
		UUID:       p.UUID,
		OwnerID:    p.Owner.GetUUID(),
		DangerType: p.DangerType,
		X:          x,
		Y:          y,
	}
}

// Live1Turn reduce remain turn and return alive
func (p *DangerObject) Live1Turn() bool {
	p.RemainTurn--
	return p.RemainTurn > 0
}
