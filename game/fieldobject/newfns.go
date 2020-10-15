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

	"github.com/kasworld/findnear"
	"github.com/kasworld/goguelike/config/gameconst"
	"github.com/kasworld/goguelike/config/lineattackdata"
	"github.com/kasworld/goguelike/enum/decaytype"
	"github.com/kasworld/goguelike/enum/fieldobjacttype"
	"github.com/kasworld/goguelike/enum/fieldobjdisplaytype"
	"github.com/kasworld/goguelike/lib/idu64str"
)

var FOIDMaker = idu64str.New("FOID")

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
		ID:          FOIDMaker.New(),
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
		ID:           FOIDMaker.New(),
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
		ID:          FOIDMaker.New(),
		FloorName:   floorname,
		ActType:     acttype,
		DisplayType: fieldobjdisplaytype.None,
		Message:     message,
	}
}

// NewRotateLineAttack arg order follow terraincmdenum
func NewRotateLineAttack(floorname string, displayType fieldobjdisplaytype.FieldObjDisplayType,
	winglen, wingcount int, degree, degreeperturn int, decay decaytype.DecayType,
	message string,
) *FieldObject {
	lineattackdata.UpdateCache360Line(winglen)
	return &FieldObject{
		ID:            FOIDMaker.New(),
		FloorName:     floorname,
		ActType:       fieldobjacttype.RotateLineAttack,
		DisplayType:   displayType,
		Message:       message,
		Degree:        degree,
		DegreePerTurn: degreeperturn,
		WingLen:       winglen,
		WingCount:     wingcount,
		Decay:         decay,
	}
}

// GetLineAttack calc dangerobj wingcount * line
func (fo *FieldObject) GetLineAttack() []findnear.XYLenList {
	rtn := make([]findnear.XYLenList, fo.WingCount)
	cache := lineattackdata.GetWingLines(fo.WingLen)
	wingdeg := 360.0 / float64(fo.WingCount)
	for wing := 0; wing < fo.WingCount; wing++ {
		deg := int(float64(wing)*wingdeg + float64(fo.Degree))
		rtn[wing] = cache[wrapInt(deg, 360)]
	}
	return rtn
}

func wrapInt(v, l int) int {
	return (v%l + l) % l
}

func (fo *FieldObject) CalcLineAttackAffectRate(rate float64, i, max int) float64 {
	switch fo.Decay {
	default:
		panic(fmt.Sprintf("invalid decaytype %v", fo))
	case decaytype.Decrease:
		return rate / float64(i+1)
	case decaytype.Even:
		return rate / (float64(max) / 2.0)
	case decaytype.Increase:
		return rate / float64(max-i)
	}
}

func NewMine(floorname string, displayType fieldobjdisplaytype.FieldObjDisplayType,
	decay decaytype.DecayType, message string,
) *FieldObject {
	return &FieldObject{
		ID:          FOIDMaker.New(),
		FloorName:   floorname,
		ActType:     fieldobjacttype.Mine,
		DisplayType: displayType,
		Message:     message,
		Radius:      -1, // not triggered
		Decay:       decay,
	}
}

func (fo *FieldObject) CalcMineAffectRate() float64 {
	switch fo.Decay {
	default:
		panic(fmt.Sprintf("invalid decaytype %v", fo))
	case decaytype.Decrease:
		return 1 / float64(fo.Radius+1)
	case decaytype.Even:
		return 1 / (gameconst.ViewPortW / 2.0)
	case decaytype.Increase:
		return 1 / float64(gameconst.ViewPortW-fo.Radius)
	}
}
