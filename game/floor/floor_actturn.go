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
	"time"

	"github.com/kasworld/goguelike/config/contagionarea"
	"github.com/kasworld/goguelike/config/slippperydata"
	"github.com/kasworld/goguelike/enum/equipslottype"
	"github.com/kasworld/goguelike/enum/way9type"

	"github.com/kasworld/goguelike/enum/condition"

	"github.com/kasworld/goguelike/enum/achievetype"
	"github.com/kasworld/goguelike/enum/fieldobjacttype"
	"github.com/kasworld/goguelike/enum/fieldobjdisplaytype"
	"github.com/kasworld/goguelike/enum/scrolltype"
	"github.com/kasworld/goguelike/enum/turnresulttype"
	"github.com/kasworld/goguelike/game/activeobject/turnresult"
	"github.com/kasworld/goguelike/game/aoactreqrsp"
	"github.com/kasworld/goguelike/game/cmd2tower"
	"github.com/kasworld/goguelike/game/fieldobject"
	"github.com/kasworld/goguelike/game/gamei"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_error"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_idcmd"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_idnoti"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
)

func (f *Floor) processTurn(turnTime time.Time) error {
	act := f.interDur.BeginAct()
	f.log.Monitor("Start Turn %v %v", f, f.interDur)
	defer func() {
		act.End()
		f.log.Monitor("End Turn %v %v", f, f.interDur)
	}()

	// wait ai run last turn
	f.aiWG.Wait()

	// prepare to process ao
	ao2ActReqRsp := make(map[gamei.ActiveObjectI]*aoactreqrsp.ActReqRsp, f.aoPosMan.Count())
	aoMapSkipTurn := make(map[string]bool) // skip noti
	aoListToProcessInTurn := make([]gamei.ActiveObjectI, 0, f.aoPosMan.Count())
	aoAliveInFloorAtStart := make([]gamei.ActiveObjectI, 0, f.aoPosMan.Count())
	for _, v := range f.aoPosMan.GetAllList() {
		ao := v.(gamei.ActiveObjectI)
		aoListToProcessInTurn = append(aoListToProcessInTurn, ao)
		if ao.IsAlive() {
			aoAliveInFloorAtStart = append(aoAliveInFloorAtStart, ao)
			if ao.GetRemainTurn2Act() <= 0 {
				req := ao.GetClearReq2Handle()
				if req != nil {
					ao2ActReqRsp[ao] = &aoactreqrsp.ActReqRsp{
						Req: *req,
					}
				}
			}
		}
	}
	f.floorCmdActStat.Add(len(ao2ActReqRsp))

	f.log.Monitor("%v ActiveObj:%v Alive:%v Acted:%v",
		f,
		len(aoListToProcessInTurn),
		len(aoAliveInFloorAtStart),
		len(ao2ActReqRsp),
	)

	for _, ao := range aoListToProcessInTurn {
		ao.PrepareNewTurn(turnTime)
	}

	// process auto portal and trapteleport steped ao
	for _, ao := range aoListToProcessInTurn {
		if !ao.IsAlive() {
			continue
		}
		if ao.GetTurnData().Condition.TestByCondition(condition.Float) {
			continue
		}
		aox, aoy, exist := f.aoPosMan.GetXYByUUID(ao.GetUUID())
		if !exist {
			f.log.Error("ao not in currentfloor %v %v", f, ao)
			continue
		}

		p, ok := f.foPosMan.Get1stObjAt(aox, aoy).(*fieldobject.FieldObject)
		if !ok {
			continue
		}
		triggered := false
		if p.ActType.AutoTrigger() && p.ActType.TriggerRate() > f.rnd.Float64() {
			triggered = true
		}
		if aoconn := ao.GetClientConn(); aoconn != nil {
			if p.DisplayType == fieldobjdisplaytype.None {
				if err := aoconn.SendNotiPacket(c2t_idnoti.FoundFieldObj,
					&c2t_obj.NotiFoundFieldObj_data{
						FloorName: f.GetName(),
						FieldObj:  p.ToPacket_FieldObjClient(aox, aoy),
					},
				); err != nil {
					f.log.Error("%v %v %v", f, ao, err)
				}
			}
			if p.ActType.TrapNoti() {
				if err := aoconn.SendNotiPacket(c2t_idnoti.ActivateTrap,
					&c2t_obj.NotiActivateTrap_data{
						FieldObjAct: p.ActType,
						Triggered:   triggered,
					},
				); err != nil {
					f.log.Error("%v %v %v", f, ao, err)
				}
			}
		}
		if !triggered {
			continue
		}

		switch p.ActType {
		default:
			fob := fieldobjacttype.GetBuffByFieldObjActType(p.ActType)
			if fob != nil {
				replaced := ao.GetBuffManager().Add(p.ActType.String(), true, true, fob)
				if replaced {
					// need noti?
				}
			}
		case fieldobjacttype.PortalAutoIn:
			p1, p2, err := f.FindUsablePortalPairAt(aox, aoy)
			if err != nil {
				f.log.Error("fail to use portal %v %v %v %v", f, p, ao, err)
				continue
			}
			ao.GetAchieveStat().Inc(achievetype.EnterPortal)
			ao.GetFieldObjActStat().Inc(p2.ActType)
			f.tower.GetReqCh() <- &cmd2tower.ActiveObjUsePortal{
				SrcFloor:  f,
				ActiveObj: ao,
				P1:        p1,
				P2:        p2,
			}

		case fieldobjacttype.Teleport:
			f.tower.GetReqCh() <- &cmd2tower.ActiveObjTrapTeleport{
				SrcFloor:     f,
				ActiveObj:    ao,
				DstFloorName: p.DstFloorName,
			}
		}
		if p.ActType.SkipAOAct() {
			aoMapSkipTurn[ao.GetUUID()] = true
			if arr, exist := ao2ActReqRsp[ao]; exist {
				arr.SetDone(
					aoactreqrsp.Act{Act: c2t_idcmd.Meditate},
					c2t_error.ActionCanceled)
			}
		}
		if p.ActType.NeedTANoti() {
			ao.SetNeedTANoti()
		}
		ao.GetFieldObjActStat().Inc(p.ActType)
		ao.GetAchieveStat().Inc(achievetype.UseFieldObj)
	}

	// handle sleep condition
	for ao, arr := range ao2ActReqRsp {
		if arr.Acted() {
			continue
		}
		if arr.Req.Act.SleepCancel() &&
			ao.GetTurnData().Condition.TestByCondition(condition.Sleep) {
			aoMapSkipTurn[ao.GetUUID()] = true
			arr.SetDone(
				aoactreqrsp.Act{Act: c2t_idcmd.Meditate},
				c2t_error.ActionCanceled)
		}
	}
	// handle remain turn2act ao
	for ao, arr := range ao2ActReqRsp {
		if arr.Acted() {
			continue
		}
		if ao.GetRemainTurn2Act() > 0 {
			arr.SetDone(
				aoactreqrsp.Act{Act: c2t_idcmd.Meditate},
				c2t_error.ActionCanceled)
		}
	}
	// handle attack
	for ao, arr := range ao2ActReqRsp {
		if arr.Acted() || arr.Req.Act != c2t_idcmd.Attack || !ao.IsAlive() {
			continue
		}
		aox, aoy, exist := f.aoPosMan.GetXYByUUID(ao.GetUUID())
		if !exist {
			f.log.Error("ao not in currentfloor %v %v", f, ao)
			arr.SetDone(
				aoactreqrsp.Act{Act: c2t_idcmd.Meditate},
				c2t_error.ActionProhibited)
			continue
		}
		atkdir := arr.Req.Dir
		if ao.GetTurnData().Condition.TestByCondition(condition.Drunken) {
			turnmod := slippperydata.Drunken[f.rnd.Intn(len(slippperydata.Drunken))]
			atkdir = atkdir.TurnDir(turnmod)
		}
		dstao, srcTile, dstTile, ec := f.canActiveObjAttack2Dir(aox, aoy, atkdir)
		if ec == c2t_error.None {
			f.aoAttackActiveObj(ao, dstao, srcTile, dstTile)
			arr.SetDone(
				aoactreqrsp.Act{Act: c2t_idcmd.Attack, Dir: atkdir},
				c2t_error.None)
			continue
		}
		arr.SetDone(
			aoactreqrsp.Act{Act: c2t_idcmd.Attack, Dir: atkdir},
			ec)
	}

	for _, ao := range aoListToProcessInTurn {
		if ao.ApplyDamageFromActiveObj() { // just killed
			// do nothing here
		}
	}

	// handle ao action except attack
	for ao, arr := range ao2ActReqRsp {
		if arr.Acted() || arr.Req.Act == c2t_idcmd.Attack || !ao.IsAlive() {
			continue
		}
		aox, aoy, exist := f.aoPosMan.GetXYByUUID(ao.GetUUID())
		if !exist {
			f.log.Error("ao not in currentfloor %v %v", f, ao)
			arr.SetDone(
				aoactreqrsp.Act{Act: c2t_idcmd.Meditate},
				c2t_error.ActionProhibited)
			continue
		}

		switch arr.Req.Act {
		default:
			f.log.Error("unknown aoact %v %v", f, arr)

		case c2t_idcmd.Meditate:
			arr.SetDone(
				aoactreqrsp.Act{Act: c2t_idcmd.Meditate},
				c2t_error.None)

		case c2t_idcmd.Move:
			mvdir, ec := f.aoAct_Move(ao, arr.Req.Dir, aox, aoy)
			arr.SetDone(
				aoactreqrsp.Act{Act: c2t_idcmd.Move, Dir: mvdir},
				ec)

		case c2t_idcmd.Pickup:
			if ao.GetTurnData().Condition.TestByCondition(condition.Float) {
				arr.SetDone(
					aoactreqrsp.Act{Act: c2t_idcmd.Pickup, UUID: arr.Req.UUID},
					c2t_error.ActionProhibited)
				continue
			}
			obj, err := f.poPosMan.GetByXYAndUUID(arr.Req.UUID, aox, aoy)
			if err != nil {
				f.log.Debug("Pickup obj not found %v %v %v", f, ao, err)
				arr.SetDone(
					aoactreqrsp.Act{Act: c2t_idcmd.Pickup, UUID: arr.Req.UUID},
					c2t_error.ObjectNotFound)
				continue
			}
			po, ok := obj.(gamei.CarryingObjectI)
			if !ok {
				f.log.Fatal("obj not carryingobject %v", po)
				arr.SetDone(
					aoactreqrsp.Act{Act: c2t_idcmd.Pickup, UUID: arr.Req.UUID},
					c2t_error.ObjectNotFound)
				continue
			}
			if err := f.poPosMan.Del(po); err != nil {
				f.log.Fatal("remove po fail %v %v %v", f, po, err)
				arr.SetDone(
					aoactreqrsp.Act{Act: c2t_idcmd.Pickup, UUID: arr.Req.UUID},
					c2t_error.ObjectNotFound)
				continue
			}
			if err := ao.DoPickup(po); err != nil {
				f.log.Error("%v %v %v", f, po, err)
				arr.SetDone(
					aoactreqrsp.Act{Act: c2t_idcmd.Pickup, UUID: arr.Req.UUID},
					c2t_error.ObjectNotFound)
				continue
			}
			arr.SetDone(
				aoactreqrsp.Act{Act: c2t_idcmd.Pickup, UUID: arr.Req.UUID},
				c2t_error.None)

		case c2t_idcmd.Drop:
			po := ao.GetInven().GetByUUID(arr.Req.UUID)
			if err := f.aoDropCarryObj(ao, aox, aoy, po); err != nil {
				f.log.Error("%v %v %v", f, ao, err)
				arr.SetDone(
					aoactreqrsp.Act{Act: c2t_idcmd.Drop, UUID: arr.Req.UUID},
					c2t_error.ObjectNotFound)
				continue
			}
			arr.SetDone(
				aoactreqrsp.Act{Act: c2t_idcmd.Drop, UUID: arr.Req.UUID},
				c2t_error.None)

		case c2t_idcmd.Equip:
			if err := ao.DoEquip(arr.Req.UUID); err != nil {
				f.log.Error("%v %v %v", f, ao, err)
				arr.SetDone(
					aoactreqrsp.Act{Act: c2t_idcmd.Equip, UUID: arr.Req.UUID},
					c2t_error.ActionProhibited)
				continue
			}
			arr.SetDone(
				aoactreqrsp.Act{Act: c2t_idcmd.Equip, UUID: arr.Req.UUID},
				c2t_error.None)

		case c2t_idcmd.UnEquip:
			if err := ao.DoUnEquip(arr.Req.UUID); err != nil {
				f.log.Error("%v %v %v", f, ao, err)
				arr.SetDone(
					aoactreqrsp.Act{Act: c2t_idcmd.UnEquip, UUID: arr.Req.UUID},
					c2t_error.ObjectNotFound)
				continue
			}
			arr.SetDone(
				aoactreqrsp.Act{Act: c2t_idcmd.UnEquip, UUID: arr.Req.UUID},
				c2t_error.None)

		case c2t_idcmd.DrinkPotion:
			po := ao.GetInven().GetByUUID(arr.Req.UUID)
			if po == nil {
				f.log.Error("po not in inventory %v %v", ao, arr.Req.UUID)
				arr.SetDone(
					aoactreqrsp.Act{Act: c2t_idcmd.DrinkPotion, UUID: arr.Req.UUID},
					c2t_error.ObjectNotFound)
				continue
			}
			if err := ao.DoUseCarryObj(arr.Req.UUID); err != nil {
				f.log.Error("%v %v %v", f, ao, err)
				arr.SetDone(
					aoactreqrsp.Act{Act: c2t_idcmd.DrinkPotion, UUID: arr.Req.UUID},
					c2t_error.ObjectNotFound)
				continue
			}
			ao.SetNeedTANoti()
			arr.SetDone(
				aoactreqrsp.Act{Act: c2t_idcmd.DrinkPotion, UUID: arr.Req.UUID},
				c2t_error.None)

		case c2t_idcmd.ReadScroll:
			po := ao.GetInven().GetByUUID(arr.Req.UUID)
			if po == nil {
				f.log.Error("po not in inventory %v %v", ao, arr.Req.UUID)
				arr.SetDone(
					aoactreqrsp.Act{Act: c2t_idcmd.ReadScroll, UUID: arr.Req.UUID},
					c2t_error.ObjectNotFound)
				continue
			}
			if so, ok := po.(gamei.ScrollI); ok && so.GetScrollType() == scrolltype.Teleport {
				err := f.aoTeleportInFloorRandom(ao)
				if err != nil {
					arr.SetDone(
						aoactreqrsp.Act{Act: c2t_idcmd.ReadScroll, UUID: arr.Req.UUID},
						c2t_error.ActionCanceled)
					f.log.Fatal("fail to teleport %v %v %v", f, ao, err)
					continue
				}
				arr.SetDone(
					aoactreqrsp.Act{Act: c2t_idcmd.ReadScroll, UUID: arr.Req.UUID},
					c2t_error.None)
				ao.GetInven().RemoveByUUID(arr.Req.UUID)
				ao.GetAchieveStat().Inc(achievetype.UseCarryObj)
				ao.GetScrollStat().Inc(scrolltype.Teleport)
			} else {
				if err := ao.DoUseCarryObj(arr.Req.UUID); err != nil {
					f.log.Error("%v %v %v", f, ao, err)
					arr.SetDone(
						aoactreqrsp.Act{Act: c2t_idcmd.ReadScroll, UUID: arr.Req.UUID},
						c2t_error.ObjectNotFound)
					continue
				}
				ao.SetNeedTANoti()
				arr.SetDone(
					aoactreqrsp.Act{Act: c2t_idcmd.ReadScroll, UUID: arr.Req.UUID},
					c2t_error.None)
			}

		case c2t_idcmd.Recycle:
			if ao.GetTurnData().Condition.TestByCondition(condition.Float) {
				arr.SetDone(
					aoactreqrsp.Act{Act: c2t_idcmd.Recycle, UUID: arr.Req.UUID},
					c2t_error.ActionProhibited)
				continue
			}
			_, ok := f.foPosMan.Get1stObjAt(aox, aoy).(*fieldobject.FieldObject)
			if !ok {
				f.log.Error("not at Recycler FieldObj %v %v", f, ao)
				arr.SetDone(
					aoactreqrsp.Act{Act: c2t_idcmd.Recycle, UUID: arr.Req.UUID},
					c2t_error.ActionProhibited)
				continue
			}
			if err := ao.DoRecycleCarryObj(arr.Req.UUID); err != nil {
				f.log.Error("%v %v %v", f, ao, err)
				arr.SetDone(
					aoactreqrsp.Act{Act: c2t_idcmd.Recycle, UUID: arr.Req.UUID},
					c2t_error.ObjectNotFound)
				continue
			}
			arr.SetDone(
				aoactreqrsp.Act{Act: c2t_idcmd.Recycle, UUID: arr.Req.UUID},
				c2t_error.None)

		case c2t_idcmd.EnterPortal:
			if ao.GetTurnData().Condition.TestByCondition(condition.Float) {
				arr.SetDone(
					aoactreqrsp.Act{Act: c2t_idcmd.EnterPortal},
					c2t_error.ActionProhibited)
				continue
			}
			p1, p2, err := f.FindUsablePortalPairAt(aox, aoy)
			if err != nil {
				arr.SetDone(
					aoactreqrsp.Act{Act: c2t_idcmd.EnterPortal},
					c2t_error.ActionProhibited)
				continue
			}
			ao.GetAchieveStat().Inc(achievetype.EnterPortal)
			ao.GetFieldObjActStat().Inc(p1.ActType)
			ao.GetFieldObjActStat().Inc(p2.ActType)
			ao.SetNeedTANoti()
			f.tower.GetReqCh() <- &cmd2tower.ActiveObjUsePortal{
				SrcFloor:  f,
				ActiveObj: ao,
				P1:        p1,
				P2:        p2,
			}
			arr.SetDone(
				aoactreqrsp.Act{Act: c2t_idcmd.EnterPortal},
				c2t_error.None)
			aoMapSkipTurn[ao.GetUUID()] = true
			f.log.Debug("manual in portal %v %v", f, ao)

		case c2t_idcmd.ActTeleport:
			if !ao.GetVisitFloor(f.GetName()).IsComplete() {
				arr.SetDone(
					aoactreqrsp.Act{Act: c2t_idcmd.ActTeleport},
					c2t_error.ActionProhibited)
				continue
			}
			err := f.aoTeleportInFloorRandom(ao)
			if err != nil {
				arr.SetDone(
					aoactreqrsp.Act{Act: c2t_idcmd.ActTeleport},
					c2t_error.ActionCanceled)
				f.log.Fatal("fail to teleport %v %v %v", f, ao, err)
			} else {
				arr.SetDone(
					aoactreqrsp.Act{Act: c2t_idcmd.ActTeleport},
					c2t_error.None)
			}
		case c2t_idcmd.KillSelf:
			ao.ReduceHP(ao.GetHP())
			arr.SetDone(
				aoactreqrsp.Act{Act: c2t_idcmd.Meditate},
				c2t_error.None)
			aoMapSkipTurn[ao.GetUUID()] = true
		}
	}

	// handle condition greasy Contagion
	for _, ao := range aoListToProcessInTurn {
		aox, aoy, exist := f.aoPosMan.GetXYByUUID(ao.GetUUID())
		if !exist {
			continue
		}

		// contagion can infected from dead ao
		bufname := fieldobjacttype.Contagion.String()
		if ao.GetBuffManager().Exist(bufname) {
			// infact other near
			aoList := f.aoPosMan.GetVPIXYObjByXYLenList(contagionarea.ContagionArea, aox, aoy, 10)
			for _, v := range aoList {
				dstAo := v.O.(gamei.ActiveObjectI)
				if !dstAo.IsAlive() { // skip dead dst
					continue
				}
				if ao.GetUUID() == dstAo.GetUUID() { // skip self
					continue
				}
				if dstAo.GetBuffManager().Exist(bufname) { // skip infected dst
					continue
				}
				if fieldobjacttype.Contagion.TriggerRate() > f.rnd.Float64() {
					fob := fieldobjacttype.GetBuffByFieldObjActType(fieldobjacttype.Contagion)
					dstAo.GetBuffManager().Add(bufname, true, true, fob)

					ao.AppendTurnResult(turnresult.New(turnresulttype.ContagionTo, dstAo, 0))
					dstAo.AppendTurnResult(turnresult.New(turnresulttype.ContagionFrom, ao, 0))

					// fmt.Printf("%v %v to %v\n", bufname, ao, dstAo)
				} else {
					ao.AppendTurnResult(turnresult.New(turnresulttype.ContagionToFail, dstAo, 0))
					dstAo.AppendTurnResult(turnresult.New(turnresulttype.ContagionFromFail, ao, 0))
				}
			}
		}

		if !ao.IsAlive() {
			continue
		}
		if ao.GetTurnData().Condition.TestByCondition(condition.Greasy) {
			eqi := f.rnd.Intn(equipslottype.EquipSlotType_Count)
			co2drop := ao.GetInven().GetEquipSlot()[eqi]
			if co2drop == nil {
				continue
			}
			if err := f.aoDropCarryObj(ao, aox, aoy, co2drop); err != nil {
				f.log.Error("%v %v %v", f, ao, err)
			}
			ao.AppendTurnResult(turnresult.New(turnresulttype.DropCarryObj, co2drop, 0))
		}

	}

	// set ao act result
	for ao, arr := range ao2ActReqRsp {
		ao.SetTurnActReqRsp(arr)
	}

	// apply terrain damage to ao by ao action
	for _, ao := range aoListToProcessInTurn {
		if !ao.IsAlive() {
			continue
		}
		aox, aoy, exist := f.aoPosMan.GetXYByUUID(ao.GetUUID())
		if !exist {
			f.log.Error("ao not in currentfloor %v %v", f, ao)
			continue
		}
		act := c2t_idcmd.Meditate
		dir := way9type.Center
		if arr, exist := ao2ActReqRsp[ao]; exist {
			act = arr.Done.Act
			dir = arr.Done.Dir
		}
		hp, sp := f.terrain.GetTiles()[aox][aoy].ActHPSPCalced(act, dir)
		if hp == 0 && sp == 0 {
			continue
		}
		ao.ApplyHPSPDecByActOnTile(hp, sp)
	}

	f.processCarryObj2floor()

	// apply act result
	for _, ao := range aoListToProcessInTurn {
		ao.ApplyTurnAct() // can die in fn
	}

	// handle ao died in Turn
	for _, ao := range aoAliveInFloorAtStart {
		if ao.IsAlive() {
			continue
		}
		aox, aoy, exist := f.aoPosMan.GetXYByUUID(ao.GetUUID())
		if !exist {
			f.log.Fatal("ao not in currentfloor %v %v", f, ao)
		}
		if err := f.ActiveObjDropCarryObjByDie(ao, aox, aoy); err != nil {
			f.log.Error("%v %v %v", f, ao, err)
		}
		ao.Noti_Death(f) // set rebirth count
		if aoconn := ao.GetClientConn(); aoconn != nil {
			if err := aoconn.SendNotiPacket(
				c2t_idnoti.Death,
				&c2t_obj.NotiDeath_data{},
			); err != nil {
				f.log.Error("%v %v %v", f, ao, err)
			}
		}
	}

	// for next turn
	// request next turn act for user
	f.sendViewportNoti(turnTime, aoListToProcessInTurn, aoMapSkipTurn)

	// requext next turn act for ai
	for _, ao := range aoListToProcessInTurn {
		if _, exist := aoMapSkipTurn[ao.GetUUID()]; exist {
			// skip leaved ao
			continue
		}
		f.aiWG.Add(1)
		go func(ao gamei.ActiveObjectI) {
			ao.RunAI(turnTime)
			f.aiWG.Done()
		}(ao)
	}

	return nil
}

// ensure po count placed on floor to map script
func (f *Floor) processCarryObj2floor() {
	for _, v := range f.poPosMan.GetAllList() {
		po, ok := v.(gamei.CarryingObjectI)
		if !ok {
			f.log.Fatal("invalid po in poPosMan %v", v)
		}
		if po.DecRemainTurnInFloor() == 0 {
			if err := f.poPosMan.Del(po); err != nil {
				f.log.Warn("remove po fail %v %v %v, maybe already removed",
					f, po, err)
			}
		}
	}

	poNeed := f.terrain.GetCarryObjCount() - f.poPosMan.Count()
	poFailCount := 0
	for i := 0; i < poNeed; i++ {
		if err := f.addNewRandCarryObj2Floor(); err != nil {
			poFailCount++
		}
	}
	if poFailCount > 0 {
		f.log.Monitor("addNewRandCarryObj2Floor fail %v %v/%v", f, poFailCount, poNeed)
	}
}

// send VPTiles, ObjectList noti when need
// request next turn act
func (f *Floor) sendViewportNoti(
	turnTime time.Time,
	aoListToProcessInTurn []gamei.ActiveObjectI,
	aoMapSkipTurn map[string]bool) {

	vpixyolistcache := f.NewCacheVPIXYOList()
	for _, ao := range aoListToProcessInTurn {
		if _, exist := aoMapSkipTurn[ao.GetUUID()]; exist {
			// skip leaved ao
			continue
		}
		aox, aoy, exist := f.aoPosMan.GetXYByUUID(ao.GetUUID())
		if !exist {
			f.log.Warn("ao not in currentfloor %v %v, skip tile, obj noti", f, ao)
			continue
		}
		if ao.GetAndClearNeedTANoti() {
			sight := ao.GetTurnData().Sight
			sightMat := f.terrain.GetViewportCache().GetByCache(aox, aoy)
			ao.UpdateBySightMat2(f, aox, aoy, sightMat,
				float32(sight))
			if aoconn := ao.GetClientConn(); aoconn != nil {
				notiTA := f.ToPacket_NotiTileArea(aox, aoy, sight)
				if err := aoconn.SendNotiPacket(
					c2t_idnoti.VPTiles,
					notiTA,
				); err != nil {
					f.log.Error("%v %v %v", f, ao, err)
				}
			}
		}
		if aoconn := ao.GetClientConn(); aoconn != nil {
			notiOL := f.ToPacket_NotiObjectList(
				turnTime,
				vpixyolistcache,
				aox, aoy, ao.GetTurnData().Sight)
			aoContidion := ao.GetTurnData().Condition
			if aoContidion.TestByCondition(condition.Blind) ||
				aoContidion.TestByCondition(condition.Invisible) {
				// if blind, invisible add self
				notiOL.ActiveObjList = append(notiOL.ActiveObjList,
					ao.ToPacket_ActiveObjClient(aox, aoy))
			}
			notiOL.ActiveObj = ao.ToPacket_PlayerActiveObjInfo()
			if err := aoconn.SendNotiPacket(
				c2t_idnoti.ObjectList,
				notiOL,
			); err != nil {
				f.log.Error("%v %v %v", f, ao, err)
			}
		}
	}
	// fmt.Printf("%v\n", vpixyolistcache)
}
