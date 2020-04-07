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
	"bytes"
	"fmt"

	"github.com/kasworld/goguelike/enum/aiplan"
	"github.com/kasworld/goguelike/enum/aotype"
)

type planObj struct {
	Name   string
	InitFn func(sai *ServerAI) int
	ActFn  func(sai *ServerAI) bool
}

var allPlanList = []planObj{
	aiplan.None:           {"None", initPlanNone, actPlanNone},
	aiplan.Chat:           {"Chat", initPlanChat, actPlanChat},
	aiplan.StrollAround:   {"StrollAround", initPlanStrollAround, actPlanStrollAround},
	aiplan.Move2Dest:      {"Move2Dest", initPlanMove2Dest, actPlanMove2Dest},
	aiplan.Revenge:        {"Revenge", initPlanRevenge, actPlanRevenge},
	aiplan.UsePortal:      {"UsePortal", initPlanUsePortal, actPlanUsePortal},
	aiplan.MoveToRecycler: {"MoveToRecycler", initPlanMoveToRecycler, actPlanMoveToRecycler},
	aiplan.RechargeSafe:   {"RechargeSafe", initPlanRechargeSafe, actPlanRechargeSafe},
	aiplan.RechargeCan:    {"RechargeCan", initPlanRechargeCan, actPlanRechargeCan},
	aiplan.PickupCarryObj: {"PickupCarryObj", initPlanPickupCarryObj, actPlanPickupCarryObj},
	aiplan.Equip:          {"Equip", initPlanEquip, actPlanEquip},
	aiplan.UsePotion:      {"UsePotion", initPlanUsePotion, actPlanUsePotion},
	aiplan.Battle:         {"Battle", initPlanBattle, actPlanBattle},
	aiplan.MoveStraight3:  {"MoveStraight3", initPlanMoveStraight3, actPlanMoveStraight3},
	aiplan.MoveStraight5:  {"MoveStraight5", initPlanMoveStraight5, actPlanMoveStraight5},
}

var aoType2aiPlan = [...]planList{
	aotype.System: planList{
		aiplan.Chat,
		aiplan.StrollAround,
		aiplan.Move2Dest,
		aiplan.Revenge,
		aiplan.UsePortal,
		aiplan.RechargeSafe,
		aiplan.RechargeCan,
		aiplan.PickupCarryObj,
		aiplan.Equip,
		aiplan.Battle,
		aiplan.MoveStraight3,
		aiplan.MoveStraight5,
	},
	aotype.User: planList{
		aiplan.StrollAround,
		aiplan.Move2Dest,
		aiplan.Revenge,
		aiplan.UsePortal,
		aiplan.MoveToRecycler,
		aiplan.RechargeSafe,
		aiplan.RechargeCan,
		aiplan.PickupCarryObj,
		aiplan.Equip,
		aiplan.UsePotion,
		aiplan.Battle,
		aiplan.MoveStraight3,
		aiplan.MoveStraight5,
	},
}

type planList []aiplan.AIPlan

func (pl planList) String() string {
	var buf bytes.Buffer
	for _, v := range pl {
		fmt.Fprintf(&buf, "%s ", allPlanList[v].Name)
	}
	return buf.String()
}

// get front plan
func (pl planList) getCurrentPlan() aiplan.AIPlan {
	return pl.getLast()
}

func (pl planList) getFront() aiplan.AIPlan {
	return pl[0]
}
func (pl planList) getLast() aiplan.AIPlan {
	return pl[len(pl)-1]
}

func (pl planList) findPos(p aiplan.AIPlan) int {
	pos := -1
	for i, v := range pl {
		if v == p {
			pos = i
			break
		}
	}
	return pos
}

// move p to front
func (pl planList) move2Front(p aiplan.AIPlan) bool {
	pos := pl.findPos(p)
	if pos != -1 {
		copy(pl[1:], pl[:pos])
		pl[0] = p
		return true
	}
	return false
}

// pop current and push back
func (pl planList) front2Last() {
	p := pl[0]
	copy(pl, pl[1:])
	pl[len(pl)-1] = p
}

func (pl planList) dup() planList {
	rtn := make(planList, len(pl))
	copy(rtn, pl)
	return rtn
}
