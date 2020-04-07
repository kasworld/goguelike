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

package fieldobject

import (
	"fmt"

	"github.com/kasworld/goguelike/enum/fieldobjacttype"
	"github.com/kasworld/goguelike/enum/fieldobjdisplaytype"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
)

type FieldObject struct {
	FloorName   string
	ID          string // uuid or from script id
	X           int
	Y           int
	DisplayType fieldobjdisplaytype.FieldObjDisplayType
	ActType     fieldobjacttype.FieldObjActType
	Message     string

	// portal
	DstPortalID string // from script not uuid

	// trapteleport
	DstFloorName string
}

func (p FieldObject) String() string {
	return fmt.Sprintf(
		"FieldObject[%v %v [%v %v] %v %v %v %v %v]",
		p.FloorName,
		p.ID,
		p.X, p.Y,
		p.ActType,
		p.DisplayType,
		p.Message,
		p.DstPortalID, p.DstFloorName,
	)
}

// IDPosI interface
func (p *FieldObject) GetUUID() string {
	return p.ID
}
func (p *FieldObject) GetPos() [2]int {
	return [2]int{p.X, p.Y}
}
func (p *FieldObject) GetXY() (int, int) {
	return p.X, p.Y
}

func (p *FieldObject) GetDisplayType() fieldobjdisplaytype.FieldObjDisplayType {
	return p.DisplayType
}

func (p *FieldObject) GetActType() fieldobjacttype.FieldObjActType {
	return p.ActType
}

func (p *FieldObject) ToPacket_FieldObjClient(x, y int) *c2t_obj.FieldObjClient {
	rtn := &c2t_obj.FieldObjClient{
		ID:          p.ID,
		X:           x,
		Y:           y,
		ActType:     p.ActType,
		DisplayType: p.DisplayType,
		Message:     p.Message,
	}
	return rtn
}
