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
	"github.com/kasworld/goguelike/enum/fieldobjacttype"
	"github.com/kasworld/goguelike/enum/fieldobjdisplaytype"
	"github.com/kasworld/uuidstr"
)

func NewPortal(floorname string, displayType fieldobjdisplaytype.FieldObjDisplayType, message string,
	acttype fieldobjacttype.FieldObjActType,
	portalID string,
	dstPortalID string,
) *FieldObject {
	return &FieldObject{
		FloorName:   floorname,
		ID:          portalID,
		ActType:     acttype,
		DisplayType: displayType,
		Message:     message,
		DstPortalID: dstPortalID,
	}
}

func NewRecycler(floorname string, displayType fieldobjdisplaytype.FieldObjDisplayType, message string,
) *FieldObject {
	return &FieldObject{
		ID:          uuidstr.New(),
		FloorName:   floorname,
		ActType:     fieldobjacttype.RecycleCarryObj,
		DisplayType: displayType,
		Message:     message,
	}
}

func NewTrapTeleport(floorname string, message string,
	dstFloorName string,
) *FieldObject {
	return &FieldObject{
		ID:           uuidstr.New(),
		FloorName:    floorname,
		ActType:      fieldobjacttype.Teleport,
		DisplayType:  fieldobjdisplaytype.None,
		Message:      message,
		DstFloorName: dstFloorName,
	}
}

func NewTrapNoArg(floorname string, displayType fieldobjdisplaytype.FieldObjDisplayType, message string,
	acttype fieldobjacttype.FieldObjActType,
) *FieldObject {
	return &FieldObject{
		ID:          uuidstr.New(),
		FloorName:   floorname,
		ActType:     acttype,
		DisplayType: fieldobjdisplaytype.None,
		Message:     message,
	}
}

func NewAreaAttack(floorname string, displayType fieldobjdisplaytype.FieldObjDisplayType, message string,
	acttype fieldobjacttype.FieldObjActType, radian, radperturn float64,
) *FieldObject {
	return &FieldObject{
		ID:          uuidstr.New(),
		FloorName:   floorname,
		ActType:     acttype,
		DisplayType: displayType,
		Message:     message,
		Radian:      radian,
		RadPerTurn:  radperturn,
	}
}

func NewMine(floorname string, displayType fieldobjdisplaytype.FieldObjDisplayType, message string,
) *FieldObject {
	return &FieldObject{
		ID:          uuidstr.New(),
		FloorName:   floorname,
		ActType:     fieldobjacttype.Mine,
		DisplayType: displayType,
		Message:     message,
		Radius:      -1, // not triggered
	}
}
