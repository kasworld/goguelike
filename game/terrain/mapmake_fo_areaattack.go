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
	"math"

	"github.com/kasworld/goguelike/enum/fieldobjacttype"
	"github.com/kasworld/goguelike/enum/fieldobjdisplaytype"
	"github.com/kasworld/goguelike/game/fieldobject"
	"github.com/kasworld/goguelike/game/terrain/roomsort"
	"github.com/kasworld/goguelike/lib/scriptparse"
)

func cmdAddAreaAttack(tr *Terrain, ca *scriptparse.CmdArgs) error {
	var x, y int
	var dispType fieldobjdisplaytype.FieldObjDisplayType
	var acttype fieldobjacttype.FieldObjActType
	var degree, perturn float64
	var message string
	if err := ca.GetArgs(&x, &y, &dispType, &acttype, &degree, &perturn, &message); err != nil {
		return err
	}
	return tr.addAreaAttack(x, y, dispType, acttype, degree/180*math.Pi, perturn/180*math.Pi, message)
}

func cmdAddAreaAttackRand(tr *Terrain, ca *scriptparse.CmdArgs) error {
	var dispType fieldobjdisplaytype.FieldObjDisplayType
	var acttype fieldobjacttype.FieldObjActType
	var degree, perturn float64
	var count int
	var message string
	if err := ca.GetArgs(&dispType, &acttype, &degree, &perturn, &count, &message); err != nil {
		return err
	}
	try := count
	for count > 0 && try > 0 {
		err := tr.addAreaAttackRand(dispType, acttype, degree/180*math.Pi, perturn/180*math.Pi, message)
		if err == nil {
			count--
		} else {
			try--
		}
	}
	if try == 0 {
		tr.log.Warn("AddAreaAttackRand add insufficient")
	}
	return nil
}

func cmdAddAreaAttackRandInRoom(tr *Terrain, ca *scriptparse.CmdArgs) error {
	var dispType fieldobjdisplaytype.FieldObjDisplayType
	var acttype fieldobjacttype.FieldObjActType
	var degree, perturn float64
	var count int
	var message string
	if err := ca.GetArgs(&dispType, &acttype, &degree, &perturn, &count, &message); err != nil {
		return err
	}
	try := count
	for count > 0 && try > 0 {
		err := tr.addAreaAttackRandInRoom(dispType, acttype, degree/180*math.Pi, perturn/180*math.Pi, message)
		if err == nil {
			count--
		} else {
			try--
		}
	}
	if try == 0 {
		tr.log.Warn("AddAreaAttackRand add insufficient")
	}
	return nil
}

func (tr *Terrain) addAreaAttack(
	x, y int, dispType fieldobjdisplaytype.FieldObjDisplayType, acttype fieldobjacttype.FieldObjActType,
	radian, perturnrad float64, message string) error {
	x, y = x%tr.Xlen, y%tr.Ylen
	if !tr.canPlaceFieldObjAt(x, y) {
		return fmt.Errorf("can not add AreaAttack at NonCharPlaceable tile %v %v", x, y)
	}
	po := fieldobject.NewAreaAttack(tr.Name, dispType, message, acttype, radian, perturnrad)
	tr.foPosMan.AddToXY(po, x, y)

	if r := tr.roomManager.GetRoomByPos(x, y); r != nil {
		r.AreaAttackCount++
	}
	return nil
}

func (tr *Terrain) addAreaAttackRand(
	dispType fieldobjdisplaytype.FieldObjDisplayType, acttype fieldobjacttype.FieldObjActType,
	radian, perturnrad float64, message string) error {

	for try := 10; try > 0; try-- {
		x, y := tr.rnd.Intn(tr.Xlen), tr.rnd.Intn(tr.Ylen)
		if !tr.canPlaceFieldObjAt(x, y) {
			continue
		}
		return tr.addAreaAttack(x, y, dispType, acttype, radian, perturnrad, message)
	}
	return fmt.Errorf("fail to addAreaAttackRand at NonCharPlaceable tile")
}

func (tr *Terrain) addAreaAttackRandInRoom(
	dispType fieldobjdisplaytype.FieldObjDisplayType, acttype fieldobjacttype.FieldObjActType,
	radian, perturnrad float64, message string) error {

	if tr.roomManager.GetCount() == 0 {
		return fmt.Errorf("no room to add AreaAttack")
	}
	roomList := tr.roomManager.GetRoomList()
	for try := 100; try > 0; try-- {
		tr.rnd.Shuffle(len(roomList), func(i, j int) {
			roomList[i], roomList[j] = roomList[j], roomList[i]
		})
		rList := roomsort.ByAreaAttackCount(roomList)
		rList.Sort()
		r := rList[0]
		x := tr.rnd.IntRange(r.Area.X, r.Area.X+r.Area.W)
		y := tr.rnd.IntRange(r.Area.Y, r.Area.Y+r.Area.H)
		if !tr.canPlaceFieldObjAt(x, y) {
			continue
		}
		return tr.addAreaAttack(x, y, dispType, acttype, radian, perturnrad, message)
	}
	return fmt.Errorf("cannot find pos in room")
}
