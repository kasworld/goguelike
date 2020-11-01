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

package aoid2floor

import (
	"fmt"
	"sync"

	"github.com/kasworld/goguelike/game/cmd2floor"
	"github.com/kasworld/goguelike/game/gamei"
	"github.com/kasworld/goguelike/lib/g2log"
)

func (taom *ActiveObjID2Floor) String() string {
	return fmt.Sprintf("ActiveObjID2Floor[%v]", taom.GetCount())
}

type ActiveObjID2Floor struct {
	mutex      sync.RWMutex   `prettystring:"hide"`
	log        *g2log.LogBase `prettystring:"hide"`
	tower      gamei.TowerI
	aoID2Floor map[string]gamei.FloorI
}

func New(tw gamei.TowerI) *ActiveObjID2Floor {
	toam := &ActiveObjID2Floor{
		aoID2Floor: make(map[string]gamei.FloorI),
		tower:      tw,
		log:        tw.Log(),
	}
	return toam
}

func (toam *ActiveObjID2Floor) Cleanup() {
}

func (toam *ActiveObjID2Floor) ActiveObjEnterTower(
	dstFloor gamei.FloorI,
	ao gamei.ActiveObjectI) error {

	toam.mutex.Lock()
	defer toam.mutex.Unlock()

	x, y, err := dstFloor.SearchRandomActiveObjPosInRoomOrRandPos()
	if err != nil {
		return err
	}

	oldfloor := toam.aoID2Floor[ao.GetUUID()]
	if oldfloor != nil {
		toam.log.Fatal("ao in other floor %v %v %v", toam, dstFloor, ao)
		toam.aoLeaveFloorNolock(ao, oldfloor)
	}
	toam.aoEnterFloorNolock(dstFloor, ao, x, y)
	return nil
}

func (toam *ActiveObjID2Floor) ActiveObjResumeToFloor(
	dstFloor gamei.FloorI, ao gamei.ActiveObjectI) error {

	toam.mutex.Lock()
	defer toam.mutex.Unlock()

	x, y, err := dstFloor.SearchRandomActiveObjPosInRoomOrRandPos()
	if err != nil {
		return err
	}

	oldfloor := toam.aoID2Floor[ao.GetUUID()]
	if oldfloor != nil && oldfloor != dstFloor {
		toam.aoLeaveFloorNolock(ao, oldfloor)
	}
	toam.aoEnterFloorNolock(dstFloor, ao, x, y)
	return nil
}

func (toam *ActiveObjID2Floor) ActiveObjMoveToFloor(
	dstFloor gamei.FloorI, ao gamei.ActiveObjectI, x, y int) error {

	toam.mutex.Lock()
	defer toam.mutex.Unlock()

	oldfloor := toam.aoID2Floor[ao.GetUUID()]
	if oldfloor != nil && oldfloor != dstFloor {
		toam.aoLeaveFloorNolock(ao, oldfloor)
	}
	toam.aoEnterFloorNolock(dstFloor, ao, x, y)
	return nil
}

func (toam *ActiveObjID2Floor) ActiveObjRebirthToFloor(
	dstFloor gamei.FloorI, ao gamei.ActiveObjectI) error {

	toam.mutex.Lock()
	defer toam.mutex.Unlock()

	oldfloor := toam.aoID2Floor[ao.GetUUID()]
	if oldfloor != nil && oldfloor != dstFloor {
		toam.aoLeaveFloorNolock(ao, oldfloor)
	}
	x, y, err := dstFloor.SearchRandomActiveObjPosInRoomOrRandPos()
	if err != nil {
		toam.log.Fatal("%v %v %v", err, dstFloor, ao)
	}
	if oldfloor != nil && oldfloor != dstFloor {
		toam.aoEnterFloorNolock(dstFloor, ao, x, y)
	}
	toam.aoRebirthToFloorNolock(dstFloor, ao, x, y)
	return nil
}

func (toam *ActiveObjID2Floor) ActiveObjLeaveFloor(ao gamei.ActiveObjectI) {
	toam.mutex.Lock()
	defer toam.mutex.Unlock()
	oldfloor := toam.aoID2Floor[ao.GetUUID()]
	if oldfloor != nil {
		toam.aoLeaveFloorNolock(ao, oldfloor)
	}
}

func (toam *ActiveObjID2Floor) aoRebirthToFloorNolock(dstFloor gamei.FloorI,
	ao gamei.ActiveObjectI, dstX, dstY int) {

	toam.aoID2Floor[ao.GetUUID()] = dstFloor
	dstFloor.GetReqCh() <- &cmd2floor.ReqRebirth2Floor{
		ActiveObj: ao,
		X:         dstX,
		Y:         dstY,
	}
}

func (toam *ActiveObjID2Floor) aoEnterFloorNolock(dstFloor gamei.FloorI,
	ao gamei.ActiveObjectI, dstX, dstY int) {

	toam.aoID2Floor[ao.GetUUID()] = dstFloor
	dstFloor.GetReqCh() <- &cmd2floor.ReqEnterFloor{
		ActiveObj: ao,
		X:         dstX,
		Y:         dstY,
	}
}

func (toam *ActiveObjID2Floor) aoLeaveFloorNolock(ao gamei.ActiveObjectI, oldfloor gamei.FloorI) {
	delete(toam.aoID2Floor, ao.GetUUID())
	oldfloor.GetReqCh() <- &cmd2floor.ReqLeaveFloor{
		ActiveObj: ao,
	}
}

func (toam *ActiveObjID2Floor) GetFloorByActiveObjID(aoid string) gamei.FloorI {
	toam.mutex.RLock()
	defer toam.mutex.RUnlock()
	return toam.aoID2Floor[aoid]
}

func (toam *ActiveObjID2Floor) GetCount() int {
	return len(toam.aoID2Floor)
}
