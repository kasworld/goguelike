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

package activeobject

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/kasworld/goguelike/enum/aiplan"
	"github.com/kasworld/goguelike/enum/way9type"
	"github.com/kasworld/goguelike/game/aoactreqrsp"
	"github.com/kasworld/goguelike/game/gamei"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_error"
	"github.com/kasworld/intervalduration"
)

type ServerAIState struct {
	aox int
	aoy int

	// export info
	InterDur        *intervalduration.IntervalDuration
	RunningPlanList planList

	turnTime        time.Time
	isAIRunning     int32 // atomic check
	movePath2Dest   [][2]int
	planCarryObj    gamei.CarryingObjectI
	planActiveObj   gamei.ActiveObjectI
	planRemainCount int
	moveDir         way9type.Way9Type

	fieldObjUseTime map[string]time.Time

	allPlanList []planObj
}

func (sai *ServerAIState) String() string {
	return fmt.Sprintf("ServerAIState[%v]", sai.RunningPlanList.GetCurrentPlan())
}

func (ao *ActiveObject) NewServerAI() *ServerAIState {
	sai := &ServerAIState{}
	sai.fieldObjUseTime = make(map[string]time.Time)
	sai.InterDur = intervalduration.New("")
	sai.RunningPlanList = aoType2aiPlan[ao.GetActiveObjType()].dup()
	ao.rnd.Shuffle(len(sai.RunningPlanList), func(i, j int) {
		sai.RunningPlanList[i], sai.RunningPlanList[j] = sai.RunningPlanList[j], sai.RunningPlanList[i]
	})

	sai.allPlanList = []planObj{
		aiplan.None:           {ao.initPlanNone, ao.actPlanNone},
		aiplan.Chat:           {ao.initPlanChat, ao.actPlanChat},
		aiplan.StrollAround:   {ao.initPlanStrollAround, ao.actPlanStrollAround},
		aiplan.Move2Dest:      {ao.initPlanMove2Dest, ao.actPlanMove2Dest},
		aiplan.Revenge:        {ao.initPlanRevenge, ao.actPlanRevenge},
		aiplan.UsePortal:      {ao.initPlanUsePortal, ao.actPlanUsePortal},
		aiplan.MoveToRecycler: {ao.initPlanMoveToRecycler, ao.actPlanMoveToRecycler},
		aiplan.RechargeSafe:   {ao.initPlanRechargeSafe, ao.actPlanRechargeSafe},
		aiplan.RechargeCan:    {ao.initPlanRechargeCan, ao.actPlanRechargeCan},
		aiplan.PickupCarryObj: {ao.initPlanPickupCarryObj, ao.actPlanPickupCarryObj},
		aiplan.Equip:          {ao.initPlanEquip, ao.actPlanEquip},
		aiplan.UsePotion:      {ao.initPlanUsePotion, ao.actPlanUsePotion},
		aiplan.Attack:         {ao.initPlanAttack, ao.actPlanAttack},
		aiplan.MoveStraight3:  {ao.initPlanMoveStraight3, ao.actPlanMoveStraight3},
		aiplan.MoveStraight5:  {ao.initPlanMoveStraight5, ao.actPlanMoveStraight5},
	}

	return sai
}

func (ao *ActiveObject) ActTurn(sai *ServerAIState, turnTime time.Time) {
	if atomic.CompareAndSwapInt32(&sai.isAIRunning, 0, 1) {
		defer atomic.AddInt32(&sai.isAIRunning, -1)
		ao.actTurn(sai, turnTime)
	} else {
		ao.log.Warn("skip ai Turn %v %v", ao, sai.isAIRunning)
	}
}

func (ao *ActiveObject) actTurn(sai *ServerAIState, turnTime time.Time) {
	if !ao.IsAlive() {
		return
	}
	if NeedChangePlan(ao.GetTurnActReqRsp()) {
		sai.planRemainCount = 0
	}

	if ao.GetAP() < 0 { // skip
		return
	}

	act := sai.InterDur.BeginAct()
	defer func() {
		act.End()
	}()

	sai.turnTime = turnTime

	aox, aoy, exist := ao.currentFloor.GetActiveObjPosMan().GetXYByUUID(ao.GetUUID())
	if !exist {
		ao.log.Error("ao not in currentfloor %v %v", ao, ao.currentFloor)
		return
	}
	sai.aox, sai.aoy = aox, aoy

	// attacked?
	if sai.RunningPlanList.GetCurrentPlan() != aiplan.Attack &&
		sai.RunningPlanList.GetCurrentPlan() != aiplan.Revenge &&
		ao.aoAttackLast(sai) != nil {

		sai.RunningPlanList.move2Front(aiplan.Revenge)
		ao.selectPlan(sai)
	} else {
		// need select new plan?
		if sai.planRemainCount <= 0 {
			ao.selectPlan(sai)
		}
	}
	if sai.planRemainCount > 0 {
		continuePlan := sai.allPlanList[sai.RunningPlanList.GetCurrentPlan()].ActFn(sai)
		if continuePlan {
			sai.planRemainCount--
		} else {
			sai.planRemainCount = 0
		}
	}
}

func NeedChangePlan(actresult *aoactreqrsp.ActReqRsp) bool {
	if actresult == nil {
		return false
	}
	if !actresult.Acted() {
		return false
	}
	return actresult.Error != c2t_error.None
}

func (ao *ActiveObject) selectPlan(sai *ServerAIState) {
	if sai.RunningPlanList.GetCurrentPlan() != aiplan.MoveToRecycler && ao.overloadRate(sai) >= 1.0 {
		sai.RunningPlanList.move2Front(aiplan.MoveToRecycler)
	}
	if sai.RunningPlanList.GetCurrentPlan() != aiplan.UsePortal && ao.floorDiscoverRate(sai) >= 1.0 {
		sai.RunningPlanList.move2Front(aiplan.UsePortal)
	}

	for tryCount := len(sai.RunningPlanList); tryCount > 0; tryCount-- {
		sai.RunningPlanList.front2Last()
		sai.planRemainCount = sai.allPlanList[sai.RunningPlanList.GetCurrentPlan()].InitFn(sai)
		if sai.planRemainCount > 0 {
			break // init success
		}
	}
}
