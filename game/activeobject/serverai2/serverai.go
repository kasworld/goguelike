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
	"fmt"
	"sync/atomic"
	"time"

	"github.com/kasworld/g2rand"
	"github.com/kasworld/goguelike/enum/aiplan"
	"github.com/kasworld/goguelike/enum/way9type"
	"github.com/kasworld/goguelike/game/aoactreqrsp"
	"github.com/kasworld/goguelike/game/gamei"
	"github.com/kasworld/goguelike/lib/g2log"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_error"
	"github.com/kasworld/intervalduration"
)

const (
	// actHistorySize = 10
	// objCacheSize    = 1000
	foCacheSize = 100
)

type ServerAI struct {
	rnd *g2rand.G2Rand `prettystring:"hide"`
	log *g2log.LogBase `prettystring:"hide"`

	aouuid string

	ao           gamei.ActiveObjectI
	currentFloor gamei.FloorI
	aox          int
	aoy          int

	isAIRunning int32
	interDur    *intervalduration.IntervalDuration

	runningPlanList planList

	movePath2Dest   [][2]int
	planCarryObj    gamei.CarryingObjectI
	planActiveObj   gamei.ActiveObjectI
	planRemainCount int
	moveDir         way9type.Way9Type

	fieldObjUseTime map[string]time.Time
	turnTime        time.Time
}

func (sai *ServerAI) String() string {
	return fmt.Sprintf("ServerAI2[%v %v]",
		sai.aouuid,
		sai.runningPlanList.getCurrentPlan().String(),
	)
}

func New(ao gamei.ActiveObjectI, l *g2log.LogBase) *ServerAI {
	sai := &ServerAI{
		rnd:    g2rand.New(),
		ao:     ao,
		log:    l,
		aouuid: ao.GetUUID(),
	}
	var err error
	// sai.objCache, err = NewObjCache(objCacheSize)
	if err != nil {
		sai.log.Fatal("ai cache init fail %v %v", ao, err)
		panic("ai cache init fail")
	}
	sai.fieldObjUseTime = make(map[string]time.Time)
	sai.interDur = intervalduration.New("")
	sai.runningPlanList = aoType2aiPlan[sai.ao.GetActiveObjType()].dup()
	sai.rnd.Shuffle(len(sai.runningPlanList), func(i, j int) {
		sai.runningPlanList[i], sai.runningPlanList[j] = sai.runningPlanList[j], sai.runningPlanList[i]
	})
	return sai
}

func (sai *ServerAI) Cleanup() {
}

func (sai *ServerAI) ActTurn(turnTime time.Time) {
	if atomic.CompareAndSwapInt32(&sai.isAIRunning, 0, 1) {
		defer atomic.AddInt32(&sai.isAIRunning, -1)
		sai.actTurn(turnTime)
	} else {
		sai.log.Warn("skip ai Turn %v %v", sai, sai.isAIRunning)
	}
}

func (sai *ServerAI) actTurn(turnTime time.Time) {
	if sai.ao == nil {
		sai.log.Warn("skip ai %v, ao is nil", sai)
		return
	}
	if !sai.ao.IsAlive() {
		return
	}
	if NeedChangePlan(sai.ao.GetTurnActReqRsp()) {
		sai.planRemainCount = 0
	}

	if sai.ao.GetAP() < 0 { // skip
		return
	}

	act := sai.interDur.BeginAct()
	defer func() {
		act.End()
	}()

	sai.turnTime = turnTime

	sai.currentFloor = sai.ao.GetCurrentFloor()
	aox, aoy, exist := sai.currentFloor.GetActiveObjPosMan().GetXYByUUID(sai.ao.GetUUID())
	if !exist {
		sai.log.Error("ao not in currentfloor %v %v", sai.ao, sai.currentFloor)
		return
	}
	sai.aox, sai.aoy = aox, aoy

	// attacked?
	if sai.runningPlanList.getCurrentPlan() != aiplan.Battle &&
		sai.runningPlanList.getCurrentPlan() != aiplan.Revenge &&
		sai.aoAttackLast() != nil {

		sai.runningPlanList.move2Front(aiplan.Revenge)
		sai.selectPlan()
	} else {
		// need select new plan?
		if sai.planRemainCount <= 0 {
			sai.selectPlan()
		}
	}
	if sai.planRemainCount > 0 {
		continuePlan := allPlanList[sai.runningPlanList.getCurrentPlan()].ActFn(sai)
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

func (sai *ServerAI) selectPlan() {
	if sai.runningPlanList.getCurrentPlan() != aiplan.MoveToRecycler && sai.overloadRate() >= 1.0 {
		sai.runningPlanList.move2Front(aiplan.MoveToRecycler)
	}
	if sai.runningPlanList.getCurrentPlan() != aiplan.UsePortal && sai.floorDiscoverRate() >= 1.0 {
		sai.runningPlanList.move2Front(aiplan.UsePortal)
	}

	for tryCount := len(sai.runningPlanList); tryCount > 0; tryCount-- {
		sai.runningPlanList.front2Last()
		sai.planRemainCount = allPlanList[sai.runningPlanList.getCurrentPlan()].InitFn(sai)
		if sai.planRemainCount > 0 {
			break // init success
		}
	}
}
