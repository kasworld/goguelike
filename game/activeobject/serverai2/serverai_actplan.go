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

package serverai2

import (
	"time"

	"github.com/kasworld/findnear"
	"github.com/kasworld/goguelike/config/gamedata"
	"github.com/kasworld/goguelike/config/viewportdata"
	"github.com/kasworld/goguelike/enum/fieldobjacttype"
	"github.com/kasworld/goguelike/enum/tile_flag"
	"github.com/kasworld/goguelike/enum/way9type"
	"github.com/kasworld/goguelike/game/fieldobject"
	"github.com/kasworld/goguelike/game/gamei"
	"github.com/kasworld/goguelike/lib/uuidposman"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_idcmd"
)

func initPlanNone(sai *ServerAI) int {
	return 0
}
func actPlanNone(sai *ServerAI) bool {
	sai.sendActNotiPacket2Floor(c2t_idcmd.Meditate, way9type.Center, "")
	return false
}

func initPlanChat(sai *ServerAI) int {
	if sai.rnd.Intn(10) == 0 {
		return 1
	}
	return 0
}
func actPlanChat(sai *ServerAI) bool {
	chat := gamedata.ChatData[sai.rnd.Intn(len(gamedata.ChatData))]
	if len(chat) > 32 {
		chat = chat[:32]
	}
	sai.ao.SetChat(chat)
	return true
}

func initPlanMove2Dest(sai *ServerAI) int {
	dstx := sai.rnd.Intn(sai.currentFloor.GetWidth())
	dsty := sai.rnd.Intn(sai.currentFloor.GetHeight())
	sai.movePath2Dest = sai.makePath2Dest(dstx, dsty)
	if len(sai.movePath2Dest) == 0 {
		return 0
	}
	return len(sai.movePath2Dest) + 10
}
func actPlanMove2Dest(sai *ServerAI) bool {
	moveDir, isContact := sai.followPath2Dest()
	if !isContact {
		// plan fail, change to other
		return false
	}
	if moveDir != way9type.Center {
		sai.sendActNotiPacket2Floor(c2t_idcmd.Move, moveDir, "")
		return true
	}
	// dest arrived
	// do dest action
	// plan change to other
	return false
}

func initPlanStrollAround(sai *ServerAI) int {
	n := sai.rnd.IntRange(len(viewportdata.ViewportXYLenList)/2, len(viewportdata.ViewportXYLenList))
	xylen := viewportdata.ViewportXYLenList[n]
	dstx, dsty := sai.aox+xylen.X, sai.aoy+xylen.Y
	sai.movePath2Dest = sai.makePath2Dest(dstx, dsty)
	if len(sai.movePath2Dest) == 0 {
		return 0
	}
	return len(sai.movePath2Dest) + 10
}
func actPlanStrollAround(sai *ServerAI) bool {
	moveDir, isContact := sai.followPath2Dest()
	if !isContact {
		// plan fail, change to other
		return false
	}
	if moveDir != way9type.Center {
		sai.sendActNotiPacket2Floor(c2t_idcmd.Move, moveDir, "")
		return true
	}
	// dest arrived
	// do dest action
	// plan change to other
	return false
}

func initPlanMoveStraight3(sai *ServerAI) int {
	if sai.moveDir == way9type.Center {
		sai.moveDir = way9type.Way9Type(sai.rnd.IntRange(1, way9type.Way9Type_Count))
	} else {
		olddir := sai.moveDir
		for {
			newdir := way9type.Way9Type(sai.rnd.IntRange(1, way9type.Way9Type_Count))
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
	w, h := sai.currentFloor.GetWidth(), sai.currentFloor.GetHeight()
	return sai.rnd.NormIntRange((w+h)/3, (w+h)/6)
}
func actPlanMoveStraight3(sai *ServerAI) bool {
	// break loop move
	if len(sai.movePath2Dest) > 1 {
		if sai.movePath2Dest[0] == [2]int{sai.aox, sai.aoy} {
			return false
		}
	}
	sai.moveDir = sai.findMovableDir3(sai.aox, sai.aoy, sai.moveDir)
	if sai.moveDir == way9type.Center {
		return false
	}
	sai.movePath2Dest = append(sai.movePath2Dest, [2]int{sai.aox, sai.aoy})
	sai.sendActNotiPacket2Floor(c2t_idcmd.Move, sai.moveDir, "")
	return true
}

func initPlanMoveStraight5(sai *ServerAI) int {
	if sai.moveDir == way9type.Center {
		sai.moveDir = way9type.Way9Type(sai.rnd.IntRange(1, way9type.Way9Type_Count))
	} else {
		olddir := sai.moveDir
		for {
			newdir := way9type.Way9Type(sai.rnd.IntRange(1, way9type.Way9Type_Count))
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
	w, h := sai.currentFloor.GetWidth(), sai.currentFloor.GetHeight()
	return sai.rnd.NormIntRange((w+h)/5, (w+h)/10)
}
func actPlanMoveStraight5(sai *ServerAI) bool {
	// break loop move
	if len(sai.movePath2Dest) > 1 {
		if sai.movePath2Dest[0] == [2]int{sai.aox, sai.aoy} {
			return false
		}
	}
	sai.moveDir = sai.findMovableDir5(sai.aox, sai.aoy, sai.moveDir)
	if sai.moveDir == way9type.Center {
		return false
	}
	sai.movePath2Dest = append(sai.movePath2Dest, [2]int{sai.aox, sai.aoy})
	sai.sendActNotiPacket2Floor(c2t_idcmd.Move, sai.moveDir, "")
	return true
}

func initPlanUsePortal(sai *ServerAI) int {
	findObj, dstx, dsty := sai.currentFloor.GetFieldObjPosMan().Search1stByXYLenList(
		viewportdata.ViewportXYLenList,
		sai.aox, sai.aoy,
		func(o uuidposman.UUIDPosI, x, y int, xylen findnear.XYLen) bool {
			p1, _, err := sai.currentFloor.FindUsablePortalPairAt(x, y)
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
		sai.movePath2Dest = sai.makePath2Dest(dstx, dsty)
		if len(sai.movePath2Dest) == 0 {
			return 0
		}
		return len(sai.movePath2Dest) + 10
	}
	return 0
}
func actPlanUsePortal(sai *ServerAI) bool {
	moveDir, isContact := sai.followPath2Dest()
	if !isContact {
		// plan fail, change to other
		return false
	}
	if moveDir != way9type.Center {
		sai.sendActNotiPacket2Floor(c2t_idcmd.Move, moveDir, "")
		return true
	}
	// dest arrived
	fl := sai.currentFloor
	inPortal, outPortal, err := fl.FindUsablePortalPairAt(sai.aox, sai.aoy)
	if err != nil {
		sai.log.Error("%v %v %v", fl, sai, err)
		return false
	}
	sai.fieldObjUseTime[outPortal.ID] = sai.turnTime
	sai.fieldObjUseTime[inPortal.ID] = sai.turnTime
	sai.sendActNotiPacket2Floor(c2t_idcmd.EnterPortal, way9type.Center, "")
	// plan change to other
	return false
}

func initPlanMoveToRecycler(sai *ServerAI) int {
	findObj, dstx, dsty := sai.currentFloor.GetFieldObjPosMan().Search1stByXYLenList(
		viewportdata.ViewportXYLenList,
		sai.aox, sai.aoy,
		func(o uuidposman.UUIDPosI, x, y int, xylen findnear.XYLen) bool {
			if iao, ok := o.(*fieldobject.FieldObject); ok {
				if iao.ActType != fieldobjacttype.RecycleCarryObj {
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
		sai.movePath2Dest = sai.makePath2Dest(dstx, dsty)
		if len(sai.movePath2Dest) == 0 {
			return 0
		}
		return len(sai.movePath2Dest) + 10
	}
	return 0
}
func actPlanMoveToRecycler(sai *ServerAI) bool {
	moveDir, isContact := sai.followPath2Dest()
	if !isContact {
		// plan fail, change to other
		return false
	}
	if moveDir != way9type.Center {
		sai.sendActNotiPacket2Floor(c2t_idcmd.Move, moveDir, "")
		return true
	}
	// dest arrived
	o := sai.currentFloor.GetFieldObjPosMan().Get1stObjAt(sai.aox, sai.aoy)
	if _, ok := o.(*fieldobject.FieldObject); ok {
		sai.fieldObjUseTime[o.GetUUID()] = time.Now()
		// sai.sendPacket2Floor(&c2t_obj.ReqActiveObjAction_data{
		// 	Act: c2t_idcmd.Recycle,
		// })
	}
	// plan change to other
	return false
}

func initPlanRechargeSafe(sai *ServerAI) int {
	if !sai.needRecharge() || sai.overloadRate() > 1 {
		return 0
	}
	// find chargeable place
	dstx, dsty, find := sai.currentFloor.GetTerrain().Search1stByXYLenList(
		viewportdata.ViewportXYLenList,
		sai.aox, sai.aoy,
		func(t *tile_flag.TileFlag) bool {
			return t.Meditateable() && t.NoBattle()
		},
	)
	if !find {
		return 0
	}
	sai.movePath2Dest = sai.makePath2Dest(dstx, dsty)
	if len(sai.movePath2Dest) == 0 {
		return 0
	}
	return len(sai.movePath2Dest) + 100
}
func actPlanRechargeSafe(sai *ServerAI) bool {
	if sai.overloadRate() > 1.0 {
		return false
	}
	tl := sai.currentFloor.GetTerrain().GetTiles()[sai.aox][sai.aoy]
	if tl.NoBattle() && tl.Meditateable() {
		if sai.ao.GetHPRate() > 0.9 && sai.ao.GetSPRate() > 0.9 {
			// plan change to other
			return false
		} else {
			return true
		}
	}
	moveDir, isContact := sai.followPath2Dest()
	if !isContact {
		// plan fail, change to other
		return false
	}
	if moveDir != way9type.Center {
		sai.sendActNotiPacket2Floor(c2t_idcmd.Move, moveDir, "")
		return true
	}

	// dest arrived
	if !tl.NoBattle() || !tl.Meditateable() {
		return false
	}

	if sai.ao.GetHPRate() > 0.9 && sai.ao.GetSPRate() > 0.9 {
		// plan change to other
		return false
	} else {
		return true
	}
}

func initPlanRechargeCan(sai *ServerAI) int {
	if !sai.needRecharge() || sai.overloadRate() > 1 {
		return 0
	}
	// find chargeable place
	dstx, dsty, find := sai.currentFloor.GetTerrain().Search1stByXYLenList(
		viewportdata.ViewportXYLenList,
		sai.aox, sai.aoy,
		func(t *tile_flag.TileFlag) bool {
			return t.Meditateable()
		},
	)
	if !find {
		return 0
	}
	sai.movePath2Dest = sai.makePath2Dest(dstx, dsty)
	if len(sai.movePath2Dest) == 0 {
		return 0
	}
	return len(sai.movePath2Dest) + 50
}
func actPlanRechargeCan(sai *ServerAI) bool {
	if sai.overloadRate() > 1.0 {
		return false
	}
	tl := sai.currentFloor.GetTerrain().GetTiles()[sai.aox][sai.aoy]
	if tl.Meditateable() {
		if sai.ao.GetHPRate() > 0.7 && sai.ao.GetSPRate() > 0.7 {
			// plan change to other
			return false
		} else {
			return true
		}
	}
	moveDir, isContact := sai.followPath2Dest()
	if !isContact {
		// plan fail, change to other
		return false
	}
	if moveDir != way9type.Center {
		sai.sendActNotiPacket2Floor(c2t_idcmd.Move, moveDir, "")
		return true
	}
	// dest arrived
	if !tl.Meditateable() {
		return false
	}

	if sai.ao.GetHPRate() > 0.9 && sai.ao.GetSPRate() > 0.9 {
		// plan change to other
		return false
	} else {
		return true
	}
}

func initPlanBattle(sai *ServerAI) int {
	// find near ao
	ter := sai.currentFloor.GetTerrain()
	findObj, dstx, dsty := sai.currentFloor.GetActiveObjPosMan().Search1stByXYLenList(
		viewportdata.ViewportXYLenList,
		sai.aox, sai.aoy,
		func(o uuidposman.UUIDPosI, x, y int, xylen findnear.XYLen) bool {
			if o.GetUUID() != sai.ao.GetUUID() &&
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
	sai.movePath2Dest = sai.makePath2Dest(dstx, dsty)
	if len(sai.movePath2Dest) == 0 {
		return 0
	}
	return len(sai.movePath2Dest) + 10
}
func actPlanBattle(sai *ServerAI) bool {
	if sai.planActiveObj == nil || !sai.planActiveObj.IsAlive() {
		return false
	}
	dstx, dsty, exist := sai.currentFloor.GetActiveObjPosMan().GetXYByUUID(sai.planActiveObj.GetUUID())
	if !exist {
		//  ActiveObj not in floor
		return false
	}

	attackdir, canAttack := sai.canAttackTo(sai.aox, sai.aoy, dstx, dsty)
	if canAttack {
		sai.sendActNotiPacket2Floor(c2t_idcmd.Attack, attackdir, "")
		return true
	}

	moveDir, isContact := sai.followPath2Dest()
	if !isContact {
		// plan fail, change to other
		return false
	}
	if moveDir != way9type.Center {
		sai.sendActNotiPacket2Floor(c2t_idcmd.Move, moveDir, "")
		return true
	}
	return false
}

func initPlanRevenge(sai *ServerAI) int {
	dstActiveObj := sai.aoAttackLast()
	if dstActiveObj == nil {
		return 0
	}
	dstx, dsty, exist := sai.currentFloor.GetActiveObjPosMan().GetXYByUUID(dstActiveObj.GetUUID())
	if !exist {
		return 0
	}
	sai.planActiveObj = dstActiveObj
	sai.movePath2Dest = sai.makePath2Dest(dstx, dsty)
	if len(sai.movePath2Dest) == 0 {
		return 0
	}
	return len(sai.movePath2Dest) + 10
}
func actPlanRevenge(sai *ServerAI) bool {
	return actPlanBattle(sai)
}

func initPlanPickupCarryObj(sai *ServerAI) int {
	// find near po
	findObj, dstx, dsty := sai.currentFloor.GetCarryObjPosMan().Search1stByXYLenList(
		viewportdata.ViewportXYLenList,
		sai.aox, sai.aoy,
		func(o uuidposman.UUIDPosI, x, y int, xylen findnear.XYLen) bool {
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
	sai.movePath2Dest = sai.makePath2Dest(dstx, dsty)
	if len(sai.movePath2Dest) == 0 {
		return 0
	}
	return len(sai.movePath2Dest) + 10
}
func actPlanPickupCarryObj(sai *ServerAI) bool {
	moveDir, isContact := sai.followPath2Dest()
	if !isContact {
		// plan fail, change to other
		return false
	}
	if moveDir != way9type.Center {
		sai.sendActNotiPacket2Floor(c2t_idcmd.Move, moveDir, "")
		return true
	}
	// dest arrived
	sai.sendActNotiPacket2Floor(c2t_idcmd.Pickup, way9type.Center,
		sai.planCarryObj.GetUUID())
	// plan change to other
	return false
}

func initPlanEquip(sai *ServerAI) int {
	if len(sai.ao.GetInven().GetEquipList()) == 0 {
		return 0
	}
	return 1
}
func actPlanEquip(sai *ServerAI) bool {
	for _, po := range sai.ao.GetInven().GetEquipSlot() {
		if po != nil && sai.needUnEquipCarryObj(po.GetBias()) {
			sai.sendActNotiPacket2Floor(c2t_idcmd.UnEquip, way9type.Center,
				po.GetUUID())
			return false
		}
	}
	for _, po := range sai.ao.GetInven().GetEquipList() {
		if po != nil && sai.isBetterCarryObj2(po.GetEquipType(), po.GetBias()) {
			sai.sendActNotiPacket2Floor(c2t_idcmd.Equip, way9type.Center,
				po.GetUUID())
			return false
		}
	}
	return false
}

func initPlanUsePotion(sai *ServerAI) int {
	if len(sai.ao.GetInven().GetPotionList()) == 0 {
		return 0
	}
	return 1
}
func actPlanUsePotion(sai *ServerAI) bool {
	for _, po := range sai.ao.GetInven().GetPotionList() {
		if po == nil {
			continue
		}
		// no use potion
	}
	return false
}
