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

package floormanager

import (
	"fmt"

	"github.com/kasworld/goguelike/config/gameconst"
	"github.com/kasworld/goguelike/config/leveldata"
	"github.com/kasworld/goguelike/game/fieldobject"
	"github.com/kasworld/goguelike/game/gamei"
)

func (fm *FloorManager) GetFloorByIndex(i int) gamei.FloorI {
	return fm.floorList[i]
}

func wrapInt(v, l int) int {
	return (v%l + l) % l
}

func (fm *FloorManager) GetFloorByIndexWrap(i int) gamei.FloorI {
	return fm.floorList[wrapInt(i, len(fm.floorList))]
}

func (fm *FloorManager) GetStartFloor() gamei.FloorI {
	return fm.startFloor
}

func (fm *FloorManager) GetFloorCount() int {
	return len(fm.floorList)
}

func (fm *FloorManager) GetFloorList() []gamei.FloorI {
	return fm.floorList
}

func (fm *FloorManager) GetSendBufferSize() int {
	return fm.sendBufferSize
}

func (fm *FloorManager) GetFloorIndexByName(id string) (int, error) {
	for i, v := range fm.floorList {
		if v == nil {
			continue
		}
		if v.GetName() == id {
			return i, nil
		}
	}
	return 0, fmt.Errorf("Floor not found %v", id)
}

func (fm *FloorManager) GetFloorByName(name string) gamei.FloorI {
	if ff := fm.floorName2Floor; ff != nil {
		if f, exist := ff[name]; exist {
			return f
		}
	}
	return nil
}

func (fm *FloorManager) FindPortalByID(id string) *fieldobject.FieldObject {
	if pp := fm.portalID2Portal; pp != nil {
		if pt, exist := pp[id]; exist {
			return pt
		}
	}
	return nil
}

func (fm *FloorManager) CalcTiles2Discover() int {
	sum := 0
	for _, v := range fm.floorList {
		sum += v.GetTerrain().GetTile2Discover()
	}
	return sum
}

func (fm *FloorManager) CalcFullDiscoverExp() int {
	totalTile := fm.CalcTiles2Discover()
	return totalTile*gameconst.ActiveObjExp_DiscoverTile +
		len(fm.floorList)*gameconst.ActiveObjExp_VisitFloor +
		totalTile*gameconst.ActiveObjExp_CompleteFloorPerTile
}

func (fm *FloorManager) CalcFullDiscoverLevel() float64 {
	return leveldata.CalcLevelFromExp(float64(fm.CalcFullDiscoverExp()))
}
