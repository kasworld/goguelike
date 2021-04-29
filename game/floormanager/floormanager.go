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

package floormanager

import (
	"fmt"
	"sync"

	"github.com/kasworld/g2rand"
	"github.com/kasworld/goguelike/config/gameconst"
	"github.com/kasworld/goguelike/enum/fieldobjacttype"
	"github.com/kasworld/goguelike/game/fieldobject"
	"github.com/kasworld/goguelike/game/floor"
	"github.com/kasworld/goguelike/game/gamei"
	"github.com/kasworld/goguelike/game/towerscript"
	"github.com/kasworld/goguelike/lib/g2log"
)

func (fm *FloorManager) String() string {
	return fmt.Sprintf("FloorManager[Tile2Discover:%v Total:%v]",
		fm.CalcTiles2Discover(),
		len(fm.floorList),
	)
}

type FloorManager struct {
	mutex         sync.RWMutex   `prettystring:"hide"`
	log           *g2log.LogBase `prettystring:"hide"`
	tower         gamei.TowerI
	terrainScript towerscript.TowerScript

	startFloor gamei.FloorI
	floorList  []gamei.FloorI

	floorName2Floor map[string]gamei.FloorI
	portalID2Portal map[string]*fieldobject.FieldObject

	sendBufferSize int
}

func New(terrainScript towerscript.TowerScript, tw gamei.TowerI) *FloorManager {
	fm := &FloorManager{
		log:             tw.Log(),
		tower:           tw,
		terrainScript:   terrainScript,
		floorList:       make([]gamei.FloorI, 0, len(terrainScript)),
		floorName2Floor: make(map[string]gamei.FloorI),
		portalID2Portal: make(map[string]*fieldobject.FieldObject),
	}
	return fm
}

func (fm *FloorManager) Cleanup() {
}

func (fm *FloorManager) Init(rnd *g2rand.G2Rand) error {
	// make floor list order to terrainscript
	tmpFloorList := make([]gamei.FloorI, len(fm.terrainScript))
	var wg sync.WaitGroup
	for i, v := range fm.terrainScript {
		wg.Add(1)
		seed := rnd.Int63()
		go func(i int, v []string, seed int64) {
			defer wg.Done()
			f := floor.New(seed, v, fm.tower)
			if err := f.Init(); err != nil {
				fm.log.Fatal("floor init fail, %v", err)
			}
			tmpFloorList[i] = f
		}(i, v, seed)
	}
	wg.Wait()

	for i, f := range tmpFloorList {
		if !f.Initialized() {
			if len(fm.terrainScript[i]) > 0 {
				fm.log.Warn("skip not initialized floor %v", fm.terrainScript[i][0])
			} else {
				fm.log.Warn("skip not initialized floor No.%v", i)
			}
			continue
		}
		if fm.startFloor == nil {
			fm.startFloor = f
		}
		fm.log.TraceService("Floor generated %v", f)

		if oldFloor, exist := fm.floorName2Floor[f.GetName()]; exist {
			return fmt.Errorf("floor name duplicate %v %v", oldFloor, f)
		} else {
			fm.floorName2Floor[f.GetName()] = f
		}
		fm.floorList = append(fm.floorList, f)
		for _, o := range f.GetTerrain().GetFieldObjPosMan().GetAllList() {
			pt, ok := o.(*fieldobject.FieldObject)
			if !ok {
				fm.log.Fatal("not *fieldobject.FieldObject %v", o)
				continue
			}
			if pt.ActType == fieldobjacttype.PortalInOut ||
				pt.ActType == fieldobjacttype.PortalIn ||
				pt.ActType == fieldobjacttype.PortalOut ||
				pt.ActType == fieldobjacttype.PortalAutoIn {
				if oldPortal, exist := fm.portalID2Portal[pt.ID]; exist {
					return fmt.Errorf("Portal id duplicate %v %v", oldPortal, pt)
				} else {
					fm.portalID2Portal[pt.ID] = pt
				}
			}
		}
	}
	// verify trapteleport dst floor name
	for _, f := range fm.floorList {
		for _, o := range f.GetTerrain().GetFieldObjPosMan().GetAllList() {
			pt, ok := o.(*fieldobject.FieldObject)
			if !ok {
				return fmt.Errorf("not *fieldobject.FieldObject %v", o)
			}
			if pt.ActType == fieldobjacttype.Teleport {
				_, exist := fm.floorName2Floor[pt.DstFloorName]
				if !exist {
					return fmt.Errorf("dstfloor name not exist %v", pt.DstFloorName)
				}
			}
		}
	}
	for _, srcPortal := range fm.portalID2Portal {
		if dstPortal, exist := fm.portalID2Portal[srcPortal.DstPortalID]; !exist {
			return fmt.Errorf("portal dest not found %v", srcPortal)
		} else {
			switch srcPortal.ActType {
			case fieldobjacttype.PortalInOut:
				if dstPortal.ActType == fieldobjacttype.PortalInOut {
					continue
				}
			case fieldobjacttype.PortalIn, fieldobjacttype.PortalAutoIn:
				if dstPortal.ActType == fieldobjacttype.PortalOut {
					continue
				}
			case fieldobjacttype.PortalOut:
				if dstPortal.ActType == fieldobjacttype.PortalIn || dstPortal.ActType == fieldobjacttype.PortalAutoIn {
					continue
				}
			}
			return fmt.Errorf("portal not match %v %v", srcPortal, dstPortal)
		}
	}
	fm.sendBufferSize = fm.calcSendBufferSize()
	return nil
}

// calcSendBufferSize find max split floor value
func (fm *FloorManager) calcSendBufferSize() int {
	rtn := 0
	for _, f := range fm.floorList {
		cv := f.GetWidth()*f.GetHeight()/gameconst.TileAreaSplitSize + 1
		if cv > rtn {
			rtn = cv
		}
	}
	return rtn + 10
}
