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

package floor

import (
	"fmt"

	"github.com/kasworld/goguelike/config/gameconst"
	"github.com/kasworld/goguelike/config/slippperydata"
	"github.com/kasworld/goguelike/config/viewportdata"
	"github.com/kasworld/goguelike/enum/achievetype"
	"github.com/kasworld/goguelike/enum/condition"
	"github.com/kasworld/goguelike/enum/turnresulttype"
	"github.com/kasworld/goguelike/enum/way9type"
	"github.com/kasworld/goguelike/game/activeobject/turnresult"
	"github.com/kasworld/goguelike/game/carryingobject"
	"github.com/kasworld/goguelike/game/gamei"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_error"
)

func (f *Floor) findActiveObjPlacabelNear(sx, sy int) (int, int, error) {
	tr := f.terrain
	if tr == nil {
		return 0, 0, fmt.Errorf("no tiles %v", f)
	}
	tiles := tr.GetTiles()
	if tiles == nil {
		return 0, 0, fmt.Errorf("no tiles %v", f)
	}
	aoPosMan := f.aoPosMan
	if aoPosMan == nil {
		return 0, 0, fmt.Errorf("no aoPosMan %v", f)
	}

	for _, v := range viewportdata.ViewportXYLenList {
		rtnX, rtnY := tr.WrapXY(sx+v.X, sy+v.Y)
		if tiles[rtnX][rtnY].CharPlaceable() && aoPosMan.Get1stObjAt(rtnX, rtnY) == nil {
			return rtnX, rtnY, nil
		}
	}
	return 0, 0, fmt.Errorf("findActiveObjPlacabelNear pos [%v %v] not found", sx, sy)
}

func (f *Floor) SearchRandomActiveObjPosInRoomOrRandPos() (int, int, error) {
	x, y, err := f.searchRandomActiveObjPosInRoom()
	if err == nil {
		return x, y, nil
	}
	x, y, err = f.SearchRandomActiveObjPos()
	return x, y, err
}

func (f *Floor) searchRandomActiveObjPosInRoom() (int, int, error) {
	tr := f.terrain
	if tr == nil {
		return 0, 0, fmt.Errorf("no tiles %v", f)
	}
	tiles := tr.GetTiles()
	if tiles == nil {
		return 0, 0, fmt.Errorf("no tiles %v", f)
	}
	rooms := tr.GetRoomList()
	if len(rooms) == 0 {
		return 0, 0, fmt.Errorf("no rooms in terrain")
	}
	for try := 10; try > 0; try-- {
		rn := f.rnd.Intn(len(rooms))
		if rn < 0 {
			f.log.Fatal("room number %v", rn)
			rn = 0
		}
		r := rooms[rn]
		x, y := r.RndPos(f.rnd)
		if tiles[x][y].CharPlaceable() {
			return x, y, nil
		}
	}
	return 0, 0, fmt.Errorf("fail to find in room pos")
}

func (f *Floor) SearchRandomActiveObjPos() (int, int, error) {
	tr := f.terrain
	if tr == nil {
		return 0, 0, fmt.Errorf("no tiles %v", f)
	}
	tiles := tr.GetTiles()
	if tiles == nil {
		return 0, 0, fmt.Errorf("no tiles %v", f)
	}
	aoPosMan := f.aoPosMan
	if aoPosMan == nil {
		return 0, 0, fmt.Errorf("no aoPosMan %v", f)
	}

	for try := 100; try > 0; try-- {
		x, y := f.rnd.Intn(f.w), f.rnd.Intn(f.h)
		if tiles[x][y].CharPlaceable() && aoPosMan.Get1stObjAt(x, y) == nil {
			return x, y, nil
		}
	}
	return 0, 0, fmt.Errorf("ActiveObj placeable pos not found")
}

func (f *Floor) aoDieDropCarryObj(ao gamei.ActiveObjectI, aox, aoy int, p gamei.CarryingObjectI) error {
	if p == nil {
		return fmt.Errorf("invalid CarryObj %v", p)
	}
	var co2Drop gamei.CarryingObjectI
	switch p.(type) {
	case gamei.MoneyI:
		if err := ao.GetInven().SubFromWallet(p.(gamei.MoneyI)); err != nil {
			return fmt.Errorf("%v %v", f, err)
		}
		co2Drop = p
		ao.AppendTurnResult(turnresult.New(turnresulttype.DropMoney, nil, float64(p.GetValue())))
	default:
		mpo := carryingobject.NewMoney(p.GetValue())
		if err := ao.GetInven().SubFromWallet(mpo); err != nil {
			co2Drop = ao.GetInven().RemoveByUUID(p.GetUUID())
			if co2Drop == nil {
				return fmt.Errorf("CarryObj not ao inven %v %v", ao, p)
			}
			ao.AppendTurnResult(turnresult.New(turnresulttype.DropCarryObj, co2Drop, 0))
		} else {
			co2Drop = mpo
			ao.AppendTurnResult(turnresult.New(turnresulttype.DropMoneyInsteadOfCarryObj, p, float64(p.GetValue())))
		}
	}
	if err := f.placeCarryObj2FloorAt(aox, aoy, co2Drop); err != nil {
		// place fail is NOT drop fail
		f.log.TraceActiveObj("CarryObj place fail po lost, %v %v", f, err)
	}
	ao.GetAchieveStat().Inc(achievetype.DropCarryObj)
	return nil
}

func (f *Floor) ActiveObjDropCarryObjByDie(ao gamei.ActiveObjectI, aox, aoy int) error {
	overloadRate := ao.GetTurnData().LoadRate

	dropRate := gameconst.DropCarryObjonActiveObjDeathRate
	minRate := overloadRate / 10
	reduceRate := 2.0
	if overloadRate > 1 {
		dropRate *= 2
		reduceRate = 1.5
	}
	// greasy increase drop rate
	if ao.GetTurnData().Condition.TestByCondition(condition.Greasy) {
		dropRate *= 2
	}

	for _, v := range ao.GetInven().GetEquipSlot() {
		isdrop := f.rnd.Float64() < dropRate
		if v != nil && isdrop {
			dropRate /= reduceRate
			if dropRate < minRate {
				dropRate = minRate
			}
			if err := f.aoDieDropCarryObj(ao, aox, aoy, v); err != nil {
				f.log.Error("%v %v %v", f, ao, err)
			}
		}
	}
	rtne, rtnp, rtns := ao.GetInven().GetTypeList()
	for _, v := range rtne {
		isdrop := f.rnd.Float64() < dropRate
		if v != nil && isdrop {
			dropRate /= reduceRate
			if dropRate < minRate {
				dropRate = minRate
			}
			if err := f.aoDieDropCarryObj(ao, aox, aoy, v); err != nil {
				f.log.Error("%v %v %v", f, ao, err)
			}
		}
	}

	for _, v := range rtnp {
		isdrop := f.rnd.Float64() < dropRate
		if v != nil && isdrop {
			dropRate /= reduceRate
			if dropRate < minRate {
				dropRate = minRate
			}
			if err := f.aoDieDropCarryObj(ao, aox, aoy, v); err != nil {
				f.log.Error("%v %v %v", f, ao, err)
			}
		}
	}

	for _, v := range rtns {
		isdrop := f.rnd.Float64() < dropRate
		if v != nil && isdrop {
			dropRate /= reduceRate
			if dropRate < minRate {
				dropRate = minRate
			}
			if err := f.aoDieDropCarryObj(ao, aox, aoy, v); err != nil {
				f.log.Error("%v %v %v", f, ao, err)
			}
		}
	}

	if dropRate > 1 {
		dropRate = 1
	}
	if dropRate < minRate {
		dropRate = minRate
	}

	moneyToDrop := ao.GetInven().GetWalletValue() * dropRate
	if moneyToDrop >= 1 {
		if err := f.aoDieDropCarryObj(ao, aox, aoy, carryingobject.NewMoney(moneyToDrop)); err != nil {
			f.log.Error("%v %v %v", f, ao, err)
		}
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////

func (f *Floor) aoTeleportInFloorRandom(ao gamei.ActiveObjectI) error {
	x, y, err := f.SearchRandomActiveObjPos()
	if err != nil {
		return err
	}
	x, y = f.terrain.WrapXY(x, y)
	if err := f.aoPosMan.UpdateToXY(ao, x, y); err != nil {
		return err
	}
	ao.SetNeedTANoti()
	return nil
}

func (f *Floor) aoAct_Move(
	ao gamei.ActiveObjectI,
	trydir way9type.Way9Type,
	aox, aoy int) (
	way9type.Way9Type, c2t_error.ErrorCode) {

	mvdir := trydir
	if mvdir == way9type.Center {
		return mvdir, c2t_error.None
	}
	if ao.GetTurnData().Condition.TestByCondition(condition.Drunken) {
		turnmod := slippperydata.Drunken[f.rnd.Intn(len(slippperydata.Drunken))]
		mvdir = mvdir.TurnDir(turnmod)
	}
	if f.terrain.GetTiles()[aox][aoy].Slippery() &&
		!ao.GetTurnData().Condition.TestByCondition(condition.Float) {
		turnmod := slippperydata.Slippery[f.rnd.Intn(len(slippperydata.Slippery))]
		mvdir = mvdir.TurnDir(turnmod)
	}
	newX, newY, ec := f.canMove2Dir(aox, aoy, mvdir)
	if ec != c2t_error.None {
		return mvdir, ec
	}
	if err := f.aoPosMan.UpdateToXY(ao, newX, newY); err != nil {
		f.log.Fatal("move ao fail %v %v %v", f, ao, err)
		return mvdir, ec
	}
	ao.SetNeedTANoti()
	ao.GetAchieveStat().Inc(achievetype.Move)
	return mvdir, c2t_error.None
}
func (f *Floor) canMove2Dir(aox, aoy int, dir way9type.Way9Type) (int, int, c2t_error.ErrorCode) {
	if !dir.IsValid() || dir == way9type.Center {
		return aox, aoy, c2t_error.None
	}

	newX, newY := aox+dir.Dx(), aoy+dir.Dy()
	newX, newY = f.terrain.WrapXY(newX, newY)
	tl := f.terrain.GetTiles()[newX][newY]
	if !tl.CharPlaceable() {
		return aox, aoy, c2t_error.MoveBlockedByTile
	}
	if !tl.CanPlaceMultiObj() {
		for _, vv := range f.aoPosMan.GetObjListAt(newX, newY) {
			v := vv.(gamei.ActiveObjectI)
			if v.IsAlive() {
				return aox, aoy, c2t_error.MoveBlockedByActiveObj
			}
		}
	}
	return newX, newY, c2t_error.None
}

func (f *Floor) aoDropCarryObj(ao gamei.ActiveObjectI, aox, aoy int, p gamei.CarryingObjectI) error {
	if p == nil {
		return fmt.Errorf("invalid po %v", p)
	}
	var po gamei.CarryingObjectI
	switch p.(type) {
	case gamei.MoneyI:
		if err := ao.GetInven().SubFromWallet(p.(gamei.MoneyI)); err != nil {
			return fmt.Errorf("%v %v", f, err)
		}
		po = p
	default:
		po = ao.GetInven().RemoveByUUID(p.GetUUID())
		if po == nil {
			return fmt.Errorf("po not ao inven %v %v", ao, po)
		}
	}
	if err := f.placeCarryObj2FloorAt(aox, aoy, po); err != nil {
		// place fail is NOT drop fail
		f.log.TraceActiveObj("CarryObj place fail po lost, %v %v", f, err)
	}
	ao.GetAchieveStat().Inc(achievetype.DropCarryObj)
	return nil
}
