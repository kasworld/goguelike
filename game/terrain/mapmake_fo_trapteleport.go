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

package terrain

import (
	"fmt"

	"github.com/kasworld/goguelike/game/fieldobject"
	"github.com/kasworld/goguelike/game/terrain/roomsort"
	"github.com/kasworld/goguelike/lib/scriptparse"
)

func cmdAddTrapTeleport(tr *Terrain, ca *scriptparse.CmdArgs) error {
	var x, y int
	var DstFloorName string
	var message string
	if err := ca.GetArgs(&x, &y, &DstFloorName, &message); err != nil {
		return err
	}
	return tr.addTrapTeleport(x, y, DstFloorName, message)
}

func cmdAddTrapTeleportRand(tr *Terrain, ca *scriptparse.CmdArgs) error {
	var DstFloorName string
	var count int
	var message string
	if err := ca.GetArgs(&count, &DstFloorName, &message); err != nil {
		return err
	}

	try := count
	for count > 0 && try > 0 {
		err := tr.addTrapTeleportRand(DstFloorName, message)
		if err == nil {
			count--
		} else {
			try--
		}
	}
	if try == 0 {
		tr.log.Warn("AddTrapTeleportRand add insufficient")
	}
	return nil
}

func cmdAddTrapTeleportRandInRoom(tr *Terrain, ca *scriptparse.CmdArgs) error {
	var DstFloorName string
	var count int
	var message string
	if err := ca.GetArgs(&count, &DstFloorName, &message); err != nil {
		return err
	}

	try := count
	for count > 0 && try > 0 {
		err := tr.addTrapTeleportRandInRoom(DstFloorName, message)
		if err == nil {
			count--
		} else {
			try--
		}
	}
	if try == 0 {
		tr.log.Warn("AddTrapTeleportRandInRoom add insufficient")
	}
	return nil
}

func (tr *Terrain) addTrapTeleport(x, y int, dstFloorName string, message string) error {
	x, y = x%tr.Xlen, y%tr.Ylen
	if !tr.canPlaceFieldObjAt(x, y) {
		return fmt.Errorf("can not add Recycler at NonCharPlaceable tile %v %v", x, y)
	}
	po := fieldobject.NewTrapTeleport(tr.Name, message, dstFloorName)
	tr.foPosMan.AddToXY(po, x, y)

	if r := tr.roomManager.GetRoomByPos(x, y); r != nil {
		r.TrapCount++
	}
	return nil
}

func (tr *Terrain) addTrapTeleportRand(dstFloorName string, message string) error {
	for try := 10; try > 0; try-- {
		x, y := tr.rnd.Intn(tr.Xlen), tr.rnd.Intn(tr.Ylen)
		if !tr.canPlaceFieldObjAt(x, y) {
			continue
		}
		if tr.isBlockWay(x, y) {
			continue
		}
		return tr.addTrapTeleport(x, y, dstFloorName, message)
	}
	return fmt.Errorf("fail to addTrapTeleportRand at NonCharPlaceable tile")
}

func (tr *Terrain) addTrapTeleportRandInRoom(dstFloorName string, message string) error {
	if tr.roomManager.GetCount() == 0 {
		return fmt.Errorf("no room to add trap")
	}
	roomList := tr.roomManager.GetRoomList()
	for try := 100; try > 0; try-- {
		tr.rnd.Shuffle(len(roomList), func(i, j int) {
			roomList[i], roomList[j] = roomList[j], roomList[i]
		})
		rList := roomsort.ByTrapCount(roomList)
		rList.Sort()
		r := rList[0]
		x := tr.rnd.IntRange(r.Area.X, r.Area.X+r.Area.W)
		y := tr.rnd.IntRange(r.Area.Y, r.Area.Y+r.Area.H)
		if !tr.canPlaceFieldObjAt(x, y) {
			continue
		}
		if tr.isBlockWay(x, y) {
			continue
		}
		return tr.addTrapTeleport(x, y, dstFloorName, message)
	}
	return fmt.Errorf("cannot find pos in room")
}
