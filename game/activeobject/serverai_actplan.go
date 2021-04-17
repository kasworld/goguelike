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

package activeobject

import (
	"time"

	"github.com/kasworld/findnear"
	"github.com/kasworld/goguelike/config/gamedata"
	"github.com/kasworld/goguelike/config/viewportdata"
	"github.com/kasworld/goguelike/enum/fieldobjacttype"
	"github.com/kasworld/goguelike/enum/tile_flag"
	"github.com/kasworld/goguelike/enum/turnresulttype"
	"github.com/kasworld/goguelike/enum/way9type"
	"github.com/kasworld/goguelike/game/attackcheck"
	"github.com/kasworld/goguelike/game/fieldobject"
	"github.com/kasworld/goguelike/game/gamei"
	"github.com/kasworld/goguelike/lib/g2log"
	"github.com/kasworld/goguelike/lib/uuidposmani"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_idcmd"
)

func ai_initPlanNone(ao *ActiveObject, sai *ServerAIState) int {
	return 0
}
func ai_actPlanNone(ao *ActiveObject, sai *ServerAIState) bool {
	ao.sendActNotiPacket2Floor(sai, c2t_idcmd.Meditate, way9type.Center, "")
	return false
}

func ai_initPlanChat(ao *ActiveObject, sai *ServerAIState) int {
	if ao.rnd.Intn(10) == 0 {
		return 1
	}
	return 0
}
func ai_actPlanChat(ao *ActiveObject, sai *ServerAIState) bool {
	chat := gamedata.ChatData[ao.rnd.Intn(len(gamedata.ChatData))]
	if len(chat) > 32 {
		chat = chat[:32]
	}
	ao.SetChat(chat)
	return true
}

func ai_initPlanMove2Dest(ao *ActiveObject, sai *ServerAIState) int {
	dstx := ao.rnd.Intn(ao.currentFloor.GetWidth())
	dsty := ao.rnd.Intn(ao.currentFloor.GetHeight())
	sai.movePath2Dest = ao.makePath2Dest(sai, dstx, dsty)
	if len(sai.movePath2Dest) == 0 {
		return 0
	}
	return len(sai.movePath2Dest) + 10
}
func ai_actPlanMove2Dest(ao *ActiveObject, sai *ServerAIState) bool {
	moveDir, isContact := ao.followPath2Dest(sai)
	if !isContact {
		// plan fail, change to other
		return false
	}
	if moveDir != way9type.Center {
		ao.sendActNotiPacket2Floor(sai, c2t_idcmd.Move, moveDir, "")
		return true
	}
	// dest arrived
	// do dest action
	// plan change to other
	return false
}

func ai_initPlanStrollAround(ao *ActiveObject, sai *ServerAIState) int {
	n := ao.rnd.IntRange(len(viewportdata.ViewportXYLenList)/2, len(viewportdata.ViewportXYLenList))
	xylen := viewportdata.ViewportXYLenList[n]
	dstx, dsty := sai.aox+xylen.X, sai.aoy+xylen.Y
	sai.movePath2Dest = ao.makePath2Dest(sai, dstx, dsty)
	if len(sai.movePath2Dest) == 0 {
		return 0
	}
	return len(sai.movePath2Dest) + 10
}
func ai_actPlanStrollAround(ao *ActiveObject, sai *ServerAIState) bool {
	moveDir, isContact := ao.followPath2Dest(sai)
	if !isContact {
		// plan fail, change to other
		return false
	}
	if moveDir != way9type.Center {
		ao.sendActNotiPacket2Floor(sai, c2t_idcmd.Move, moveDir, "")
		return true
	}
	// dest arrived
	// do dest action
	// plan change to other
	return false
}

func ai_initPlanMoveStraight3(ao *ActiveObject, sai *ServerAIState) int {
	if sai.moveDir == way9type.Center {
		sai.moveDir = way9type.Way9Type(ao.rnd.IntRange(1, way9type.Way9Type_Count))
	} else {
		olddir := sai.moveDir
		for {
			newdir := way9type.Way9Type(ao.rnd.IntRange(1, way9type.Way9Type_Count))
			if newdir != olddir && newdir.ReverseDir() != olddir {
				sai.moveDir = newdir
				break
			}
		}
	}
	// record plan start pos
	sai.movePath2Dest = [][2]int{
		{sai.aox, sai.aoy},
	}
	w, h := ao.currentFloor.GetWidth(), ao.currentFloor.GetHeight()
	return ao.rnd.NormIntRange((w+h)/3, (w+h)/6)
}
func ai_actPlanMoveStraight3(ao *ActiveObject, sai *ServerAIState) bool {
	// break loop move
	if len(sai.movePath2Dest) > 1 {
		if sai.movePath2Dest[0] == [2]int{sai.aox, sai.aoy} {
			return false
		}
	}
	sai.moveDir = ao.findMovableDir3(sai, sai.aox, sai.aoy, sai.moveDir)
	if sai.moveDir == way9type.Center {
		return false
	}
	sai.movePath2Dest = append(sai.movePath2Dest, [2]int{sai.aox, sai.aoy})
	ao.sendActNotiPacket2Floor(sai, c2t_idcmd.Move, sai.moveDir, "")
	return true
}

func ai_initPlanMoveStraight5(ao *ActiveObject, sai *ServerAIState) int {
	if sai.moveDir == way9type.Center {
		sai.moveDir = way9type.Way9Type(ao.rnd.IntRange(1, way9type.Way9Type_Count))
	} else {
		olddir := sai.moveDir
		for {
			newdir := way9type.Way9Type(ao.rnd.IntRange(1, way9type.Way9Type_Count))
			if newdir != olddir && newdir.ReverseDir() != olddir {
				sai.moveDir = newdir
				break
			}
		}
	}
	// record plan start pos
	sai.movePath2Dest = [][2]int{
		{sai.aox, sai.aoy},
	}
	w, h := ao.currentFloor.GetWidth(), ao.currentFloor.GetHeight()
	return ao.rnd.NormIntRange((w+h)/5, (w+h)/10)
}
func ai_actPlanMoveStraight5(ao *ActiveObject, sai *ServerAIState) bool {
	// break loop move
	if len(sai.movePath2Dest) > 1 {
		if sai.movePath2Dest[0] == [2]int{sai.aox, sai.aoy} {
			return false
		}
	}
	sai.moveDir = ao.findMovableDir5(sai, sai.aox, sai.aoy, sai.moveDir)
	if sai.moveDir == way9type.Center {
		return false
	}
	sai.movePath2Dest = append(sai.movePath2Dest, [2]int{sai.aox, sai.aoy})
	ao.sendActNotiPacket2Floor(sai, c2t_idcmd.Move, sai.moveDir, "")
	return true
}

func ai_initPlanUsePortal(ao *ActiveObject, sai *ServerAIState) int {
	findObj, dstx, dsty := ao.currentFloor.GetFieldObjPosMan().Search1stByXYLenList(
		viewportdata.ViewportXYLenList,
		sai.aox, sai.aoy,
		func(o uuidposmani.UUIDPosI, x, y int, xylen findnear.XYLen) bool {
			p1, _, err := ao.currentFloor.FindUsablePortalPairAt(x, y)
			if err != nil {
				return false
			}
			lastTime := sai.fieldObjUseTime[p1.ID]
			if lastTime.Add(time.Second * 120).After(sai.turnTime) {
				return false
			}
			return true
		},
	)
	if findObj != nil {
		sai.movePath2Dest = ao.makePath2Dest(sai, dstx, dsty)
		if len(sai.movePath2Dest) == 0 {
			return 0
		}
		return len(sai.movePath2Dest) + 10
	}
	return 0
}
func ai_actPlanUsePortal(ao *ActiveObject, sai *ServerAIState) bool {
	moveDir, isContact := ao.followPath2Dest(sai)
	if !isContact {
		// plan fail, change to other
		return false
	}
	if moveDir != way9type.Center {
		ao.sendActNotiPacket2Floor(sai, c2t_idcmd.Move, moveDir, "")
		return true
	}
	// dest arrived
	fl := ao.currentFloor
	inPortal, outPortal, err := fl.FindUsablePortalPairAt(sai.aox, sai.aoy)
	if err != nil {
		if inPortal != nil { // srcPortal found
			g2log.Error("%v %v %v", fl, ao, err)
		} else {
			// igonore no srcPortal
		}
		return false
	}
	sai.fieldObjUseTime[outPortal.ID] = sai.turnTime
	sai.fieldObjUseTime[inPortal.ID] = sai.turnTime
	ao.sendActNotiPacket2Floor(sai, c2t_idcmd.EnterPortal, way9type.Center, "")
	// plan change to other
	return false
}

func ai_initPlanMoveToRecycler(ao *ActiveObject, sai *ServerAIState) int {
	findObj, dstx, dsty := ao.currentFloor.GetFieldObjPosMan().Search1stByXYLenList(
		viewportdata.ViewportXYLenList,
		sai.aox, sai.aoy,
		func(o uuidposmani.UUIDPosI, x, y int, xylen findnear.XYLen) bool {
			if fo, ok := o.(*fieldobject.FieldObject); ok {
				if fo.ActType != fieldobjacttype.RecycleCarryObj {
					return false
				}
				lastTime := sai.fieldObjUseTime[o.GetUUID()]
				if lastTime.Add(time.Second * 120).After(sai.turnTime) {
					return false
				}
				return true
			}
			return false
		},
	)
	if findObj != nil {
		sai.movePath2Dest = ao.makePath2Dest(sai, dstx, dsty)
		if len(sai.movePath2Dest) == 0 {
			return 0
		}
		return len(sai.movePath2Dest) + 10
	}
	return 0
}
func ai_actPlanMoveToRecycler(ao *ActiveObject, sai *ServerAIState) bool {
	moveDir, isContact := ao.followPath2Dest(sai)
	if !isContact {
		// plan fail, change to other
		return false
	}
	if moveDir != way9type.Center {
		ao.sendActNotiPacket2Floor(sai, c2t_idcmd.Move, moveDir, "")
		return true
	}
	// dest arrived
	o := ao.currentFloor.GetFieldObjPosMan().Get1stObjAt(sai.aox, sai.aoy)
	if _, ok := o.(*fieldobject.FieldObject); ok {
		sai.fieldObjUseTime[o.GetUUID()] = time.Now()
		// sai.sendPacket2Floor(&c2t_obj.ReqActiveObjAction_data{
		// 	Act: c2t_idcmd.Recycle,
		// })
	}
	// plan change to other
	return false
}

func ai_initPlanRechargeSafe(ao *ActiveObject, sai *ServerAIState) int {
	if !ao.needRecharge(sai) || ao.overloadRate(sai) > 1 {
		return 0
	}
	// find chargeable place
	dstx, dsty, find := ao.currentFloor.GetTerrain().Search1stByXYLenList(
		viewportdata.ViewportXYLenList,
		sai.aox, sai.aoy,
		func(t *tile_flag.TileFlag) bool {
			return t.Meditateable() && t.NoBattle()
		},
	)
	if !find {
		return 0
	}
	sai.movePath2Dest = ao.makePath2Dest(sai, dstx, dsty)
	if len(sai.movePath2Dest) == 0 {
		return 0
	}
	return len(sai.movePath2Dest) + 100
}
func ai_actPlanRechargeSafe(ao *ActiveObject, sai *ServerAIState) bool {
	if ao.overloadRate(sai) > 1.0 {
		return false
	}
	tl := ao.currentFloor.GetTerrain().GetTiles()[sai.aox][sai.aoy]
	if tl.NoBattle() && tl.Meditateable() {
		if ao.GetHPRate() > 0.9 && ao.GetSPRate() > 0.9 {
			// plan change to other
			return false
		} else {
			return true
		}
	}
	moveDir, isContact := ao.followPath2Dest(sai)
	if !isContact {
		// plan fail, change to other
		return false
	}
	if moveDir != way9type.Center {
		ao.sendActNotiPacket2Floor(sai, c2t_idcmd.Move, moveDir, "")
		return true
	}

	// dest arrived
	if !tl.NoBattle() || !tl.Meditateable() {
		return false
	}

	// change plan if damaged from ao, tile, fo
	for _, v := range ao.GetTurnResultList() {
		switch v.GetTurnResultType() {
		case turnresulttype.AttackedFrom, turnresulttype.DeadByTile,
			turnresulttype.HPDamageFromTrap, turnresulttype.SPDamageFromTrap:
			return false
		}
	}

	if ao.GetHPRate() > 0.9 && ao.GetSPRate() > 0.9 {
		// plan change to other
		return false
	} else {
		return true
	}
}

func ai_initPlanRechargeCan(ao *ActiveObject, sai *ServerAIState) int {
	if !ao.needRecharge(sai) || ao.overloadRate(sai) > 1 {
		return 0
	}
	// find chargeable place
	dstx, dsty, find := ao.currentFloor.GetTerrain().Search1stByXYLenList(
		viewportdata.ViewportXYLenList,
		sai.aox, sai.aoy,
		func(t *tile_flag.TileFlag) bool {
			return t.Meditateable()
		},
	)
	if !find {
		return 0
	}
	sai.movePath2Dest = ao.makePath2Dest(sai, dstx, dsty)
	if len(sai.movePath2Dest) == 0 {
		return 0
	}
	return len(sai.movePath2Dest) + 50
}
func ai_actPlanRechargeCan(ao *ActiveObject, sai *ServerAIState) bool {
	if ao.overloadRate(sai) > 1.0 {
		return false
	}
	tl := ao.currentFloor.GetTerrain().GetTiles()[sai.aox][sai.aoy]
	if tl.Meditateable() {
		if ao.GetHPRate() > 0.7 && ao.GetSPRate() > 0.7 {
			// plan change to other
			return false
		} else {
			return true
		}
	}
	moveDir, isContact := ao.followPath2Dest(sai)
	if !isContact {
		// plan fail, change to other
		return false
	}
	if moveDir != way9type.Center {
		ao.sendActNotiPacket2Floor(sai, c2t_idcmd.Move, moveDir, "")
		return true
	}
	// dest arrived
	if !tl.Meditateable() {
		return false
	}

	// change plan if damaged from ao, tile, fo
	for _, v := range ao.GetTurnResultList() {
		switch v.GetTurnResultType() {
		case turnresulttype.AttackedFrom, turnresulttype.DeadByTile,
			turnresulttype.HPDamageFromTrap, turnresulttype.SPDamageFromTrap:
			return false
		}
	}

	if ao.GetHPRate() > 0.9 && ao.GetSPRate() > 0.9 {
		// plan change to other
		return false
	} else {
		return true
	}
}

func ai_initPlanAttack(ao *ActiveObject, sai *ServerAIState) int {
	// find near ao
	ter := ao.currentFloor.GetTerrain()
	findObj, dstx, dsty := ao.currentFloor.GetActiveObjPosMan().Search1stByXYLenList(
		viewportdata.ViewportXYLenList,
		sai.aox, sai.aoy,
		func(o uuidposmani.UUIDPosI, x, y int, xylen findnear.XYLen) bool {
			if o.GetUUID() != ao.GetUUID() &&
				o.(gamei.ActiveObjectI).IsAlive() &&
				ter.GetTiles()[x][y].CanBattle() {
				return true
			}
			return false
		},
	)
	if findObj == nil {
		return 0
	}
	dstActiveObj := findObj.(gamei.ActiveObjectI)
	sai.planActiveObj = dstActiveObj
	sai.movePath2Dest = ao.makePath2Dest(sai, dstx, dsty)
	if len(sai.movePath2Dest) == 0 {
		return 0
	}
	return len(sai.movePath2Dest) + 10
}
func ai_actPlanAttack(ao *ActiveObject, sai *ServerAIState) bool {
	if sai.planActiveObj == nil || !sai.planActiveObj.IsAlive() {
		return false
	}
	dstx, dsty, exist := ao.currentFloor.GetActiveObjPosMan().GetXYByUUID(sai.planActiveObj.GetUUID())
	if !exist {
		//  ActiveObj not in floor
		return false
	}

	ter := ao.currentFloor.GetTerrain()
	attackdir, canAttack := attackcheck.CanBasicAttackTo(
		ter.GetTiles(), sai.aox, sai.aoy, dstx, dsty)
	if canAttack {
		ao.sendActNotiPacket2Floor(sai, c2t_idcmd.Attack, attackdir, "")
		return true
	}

	attackdir, canAttack = attackcheck.CanLongAttackTo(
		ter.GetTiles(), sai.aox, sai.aoy, dstx, dsty)
	if canAttack {
		ao.sendActNotiPacket2Floor(sai, c2t_idcmd.AttackLong, attackdir, "")
		return true
	}

	moveDir, isContact := ao.followPath2Dest(sai)
	if !isContact {
		// plan fail, change to other
		return false
	}
	if moveDir != way9type.Center {
		ao.sendActNotiPacket2Floor(sai, c2t_idcmd.Move, moveDir, "")
		return true
	}
	return false
}

func ai_initPlanRevenge(ao *ActiveObject, sai *ServerAIState) int {
	dstActiveObj := ao.aoAttackLast(sai)
	if dstActiveObj == nil {
		return 0
	}
	dstx, dsty, exist := ao.currentFloor.GetActiveObjPosMan().GetXYByUUID(dstActiveObj.GetUUID())
	if !exist {
		return 0
	}
	sai.planActiveObj = dstActiveObj
	sai.movePath2Dest = ao.makePath2Dest(sai, dstx, dsty)
	if len(sai.movePath2Dest) == 0 {
		return 0
	}
	return len(sai.movePath2Dest) + 10
}
func ai_actPlanRevenge(ao *ActiveObject, sai *ServerAIState) bool {
	return ai_actPlanAttack(ao, sai)
}

func ai_initPlanPickupCarryObj(ao *ActiveObject, sai *ServerAIState) int {
	// find near po
	findObj, dstx, dsty := ao.currentFloor.GetCarryObjPosMan().Search1stByXYLenList(
		viewportdata.ViewportXYLenList,
		sai.aox, sai.aoy,
		func(o uuidposmani.UUIDPosI, x, y int, xylen findnear.XYLen) bool {
			if _, ok := o.(gamei.CarryingObjectI); ok {
				// str := sai.objCache.Get(o.GetUUID())
				// if str != "drop" {
				return true
				// }
			}
			return false
		},
	)
	if findObj == nil {
		return 0
	}
	dstCarryObj := findObj.(gamei.CarryingObjectI)
	sai.planCarryObj = dstCarryObj
	sai.movePath2Dest = ao.makePath2Dest(sai, dstx, dsty)
	if len(sai.movePath2Dest) == 0 {
		return 0
	}
	return len(sai.movePath2Dest) + 10
}
func ai_actPlanPickupCarryObj(ao *ActiveObject, sai *ServerAIState) bool {
	moveDir, isContact := ao.followPath2Dest(sai)
	if !isContact {
		// plan fail, change to other
		return false
	}
	if moveDir != way9type.Center {
		ao.sendActNotiPacket2Floor(sai, c2t_idcmd.Move, moveDir, "")
		return true
	}
	// dest arrived
	ao.sendActNotiPacket2Floor(sai, c2t_idcmd.Pickup, way9type.Center,
		sai.planCarryObj.GetUUID())
	// plan change to other
	return false
}

func ai_initPlanEquip(ao *ActiveObject, sai *ServerAIState) int {
	if len(ao.GetInven().GetEquipList()) == 0 {
		return 0
	}
	return 1
}
func ai_actPlanEquip(ao *ActiveObject, sai *ServerAIState) bool {
	for _, po := range ao.GetInven().GetEquipSlot() {
		if po != nil && ao.needUnEquipCarryObj(sai, po.GetBias()) {
			ao.sendActNotiPacket2Floor(sai, c2t_idcmd.UnEquip, way9type.Center,
				po.GetUUID())
			return false
		}
	}
	for _, po := range ao.GetInven().GetEquipList() {
		if po != nil && ao.isBetterCarryObj2(sai, po.GetEquipType(), po.GetBias()) {
			ao.sendActNotiPacket2Floor(sai, c2t_idcmd.Equip, way9type.Center,
				po.GetUUID())
			return false
		}
	}
	return false
}

func ai_initPlanUsePotion(ao *ActiveObject, sai *ServerAIState) int {
	if len(ao.GetInven().GetPotionList()) == 0 {
		return 0
	}
	return 1
}
func ai_actPlanUsePotion(ao *ActiveObject, sai *ServerAIState) bool {
	for _, po := range ao.GetInven().GetPotionList() {
		if po == nil {
			continue
		}
		// no use potion
	}
	return false
}
