// Copyright 2014,2015,2016,2017,2018,2019,2020,2021 SeukWon Kang (kasworld@gmail.com)
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
	"github.com/kasworld/goguelike/enum/dangertype"
	"github.com/kasworld/goguelike/enum/decaytype"
	"github.com/kasworld/goguelike/enum/fieldobjacttype"
	"github.com/kasworld/goguelike/enum/fieldobjdisplaytype"
	"github.com/kasworld/goguelike/game/dangerobject"
)

func NewMine(floorname string, displayType fieldobjdisplaytype.FieldObjDisplayType,
	decay decaytype.DecayType, message string,
) *FieldObject {
	fo := &FieldObject{
		ID:            FOIDMaker.New(),
		FloorName:     floorname,
		ActType:       fieldobjacttype.Mine,
		DisplayType:   displayType,
		Message:       message,
		CurrentRadius: -1, // not triggered
		Decay:         decay,
	}
	fo.makeMineDOs()
	return fo
}

func (fo *FieldObject) makeMineDOs() {
	for r := 0; r < gameconst.ViewPortW; r++ {
		fo.premakeMineXYLDOs[r] = make([]XYlenDO, len(MineData[r]))
		var rr float64
		switch fo.Decay {
		default:
			panic(fmt.Sprintf("invalid decaytype %v", fo))
		case decaytype.Decrease:
			rr = 1 / float64(r+1)
		case decaytype.Even:
			rr = 1 / (gameconst.ViewPortW / 2.0)
		case decaytype.Increase:
			rr = 1 / float64(gameconst.ViewPortW-r)
		}
		for i, v := range MineData[r] {
			xyldo := XYlenDO{
				X:  v.X,
				Y:  v.Y,
				DO: dangerobject.NewFOAttact(fo, dangertype.MineExplode, rr),
			}
			fo.premakeMineXYLDOs[r][i] = xyldo
		}
	}
}

func (fo *FieldObject) GetMineDO() []XYlenDO {
	return fo.premakeMineXYLDOs[fo.CurrentRadius]
}

var MineData [gameconst.ViewPortW]findnear.XYLenList

func init() {
	fullData := findnear.NewXYLenList(
		gameconst.ClientViewPortW, gameconst.ClientViewPortH)[:gameconst.ViewPortWH]

	st := 0
	curRadius := 1.0
	for i, v := range fullData {
		if v.L > curRadius {
			MineData[int(curRadius)-1] = fullData[st:i]
			st = i
			curRadius++
		}
	}
	MineData[int(curRadius)-1] = fullData[st:]
	// fmt.Printf("%v", MineData)
}
