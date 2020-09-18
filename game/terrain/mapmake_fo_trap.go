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

package terrain

import (
	"fmt"

	"github.com/kasworld/goguelike/enum/fieldobjacttype"
	"github.com/kasworld/goguelike/enum/fieldobjdisplaytype"
	"github.com/kasworld/goguelike/game/fieldobject"
	"github.com/kasworld/goguelike/game/terrain/roomsort"
	"github.com/kasworld/goguelike/lib/scriptparse"
)

func cmdAddTrap(tr *Terrain, ca *scriptparse.CmdArgs) error {
	var x, y int
	var dispType fieldobjdisplaytype.FieldObjDisplayType
	var acttype fieldobjacttype.FieldObjActType
	var message string
	if err := ca.GetArgs(&x, &y, &dispType, &acttype, &message); err != nil {
		return err
	}
	return tr.addTrap(x, y, dispType, acttype, message)
}

func cmdAddTrapRand(tr *Terrain, ca *scriptparse.CmdArgs) error {
	var dispType fieldobjdisplaytype.FieldObjDisplayType
	var acttype fieldobjacttype.FieldObjActType
	var count int
	var message string
	if err := ca.GetArgs(&count, &dispType, &acttype, &message); err != nil {
		return err
	}
	try := count
	for count > 0 && try > 0 {
		err := tr.addTrapRand(dispType, acttype, message)
		if err == nil {
			count--
		} else {
			try--
		}
	}
	if try == 0 {
		tr.log.Warn("AddTrapRand add insufficient")
	}
	return nil
}

func cmdAddTrapRandInRoom(tr *Terrain, ca *scriptparse.CmdArgs) error {
	var dispType fieldobjdisplaytype.FieldObjDisplayType
	var acttype fieldobjacttype.FieldObjActType
	var count int
	var message string
	if err := ca.GetArgs(&count, &dispType, &acttype, &message); err != nil {
		return err
	}
	try := count
	for count > 0 && try > 0 {
		err := tr.addTrapRandInRoom(dispType, acttype, message)
		if err == nil {
			count--
		} else {
			try--
		}
	}
	if try == 0 {
		tr.log.Warn("AddTrapRand add insufficient")
	}
	return nil
}

func (tr *Terrain) addTrap(x, y int, dispType fieldobjdisplaytype.FieldObjDisplayType, actType fieldobjacttype.FieldObjActType, message string) error {
	x, y = x%tr.Xlen, y%tr.Ylen
	if !tr.canPlaceFieldObjAt(x, y) {
		return fmt.Errorf("can not add Recycler at NonCharPlaceable tile %v %v", x, y)
	}
	po := fieldobject.NewTrapNoArg(tr.Name, dispType, message, actType)
	tr.foPosMan.AddToXY(po, x, y)

	if r := tr.roomManager.GetRoomByPos(x, y); r != nil {
		r.TrapCount++
	}
	return nil
}

func (tr *Terrain) addTrapRand(dispType fieldobjdisplaytype.FieldObjDisplayType, actType fieldobjacttype.FieldObjActType, message string) error {

	for try := 10; try > 0; try-- {
		x, y := tr.rnd.Intn(tr.Xlen), tr.rnd.Intn(tr.Ylen)
		if !tr.canPlaceFieldObjAt(x, y) {
			continue
		}
		return tr.addTrap(x, y, dispType, actType, message)
	}
	return fmt.Errorf("fail to addTrapRand at NonCharPlaceable tile")
}

func (tr *Terrain) addTrapRandInRoom(dispType fieldobjdisplaytype.FieldObjDisplayType, actType fieldobjacttype.FieldObjActType, message string) error {

	if tr.roomManager.GetCount() == 0 {
		return fmt.Errorf("no room to add Recycler")
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
		return tr.addTrap(x, y, dispType, actType, message)
	}
	return fmt.Errorf("cannot find pos in room")
}
