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

	"github.com/kasworld/goguelike/config/gameconst"
	"github.com/kasworld/goguelike/enum/decaytype"
	"github.com/kasworld/goguelike/enum/fieldobjacttype"
	"github.com/kasworld/goguelike/enum/fieldobjdisplaytype"
)

func NewMine(floorname string, displayType fieldobjdisplaytype.FieldObjDisplayType,
	decay decaytype.DecayType, message string,
) *FieldObject {
	fo := &FieldObject{
		ID:          FOIDMaker.New(),
		FloorName:   floorname,
		ActType:     fieldobjacttype.Mine,
		DisplayType: displayType,
		Message:     message,
		Radius:      -1, // not triggered
		Decay:       decay,
	}

	return fo
}

func (fo *FieldObject) makeMineDOs() {

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
