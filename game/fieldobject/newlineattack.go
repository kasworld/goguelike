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
	"math"

	"github.com/kasworld/goguelike/enum/dangertype"
	"github.com/kasworld/goguelike/enum/decaytype"
	"github.com/kasworld/goguelike/enum/fieldobjacttype"
	"github.com/kasworld/goguelike/enum/fieldobjdisplaytype"
	"github.com/kasworld/goguelike/game/dangerobject"
	"github.com/kasworld/goguelike/lib/lineofsight"
)

type XYlenDO struct {
	X  int
	Y  int
	DO *dangerobject.DangerObject
}

// NewRotateLineAttack arg order follow terraincmdenum
func NewRotateLineAttack(floorname string, displayType fieldobjdisplaytype.FieldObjDisplayType,
	winglen, wingcount int, degree, degreeperturn int, decay decaytype.DecayType,
	message string,
) *FieldObject {
	fo := &FieldObject{
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
	for deg := 0; deg < 360; deg++ {
		fo.premakeWingsXYLDOs[deg] = fo.makeWingDOs(deg)
	}
	return fo
}

func (fo *FieldObject) makeWingDOs(deg int) []XYlenDO {
	rad := float64(deg) / 180.0 * math.Pi
	dx := float64(fo.WingLen) * math.Cos(rad)
	dy := float64(fo.WingLen) * math.Sin(rad)
	xyll := lineofsight.MakePosLenList(0.5, 0.5, dx+0.5, dy+0.5).ToCellLenList()
	rtn := make([]XYlenDO, len(xyll))
	for i, v := range xyll {
		rr := fo.calcLineAttackAffectRate(v.L, i, len(xyll))
		rtn[i] = XYlenDO{
			X:  v.X,
			Y:  v.Y,
			DO: dangerobject.NewFOAttact(fo, dangertype.RotateLineAttack, rr),
		}
	}
	return rtn
}

func (fo *FieldObject) GetWingByNum(wing int) []XYlenDO {
	wingdeg := 360.0 / float64(fo.WingCount)
	deg := int(float64(wing)*wingdeg + float64(fo.Degree))
	return fo.premakeWingsXYLDOs[wrapInt(deg, 360)]
}

func wrapInt(v, l int) int {
	return (v%l + l) % l
}

func (fo *FieldObject) calcLineAttackAffectRate(rate float64, i, max int) float64 {
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
