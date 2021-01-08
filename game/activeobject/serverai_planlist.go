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
	"github.com/kasworld/goguelike/enum/aiplan"
	"github.com/kasworld/goguelike/enum/aotype"
)

type planObj struct {
	InitFn func(ao *ActiveObject, sai *ServerAIState) int
	ActFn  func(ao *ActiveObject, sai *ServerAIState) bool
}

var aoType2aiPlan = [...]aiplan.PlanList{
	aotype.System: aiplan.PlanList{
		aiplan.Chat,
		aiplan.StrollAround,
		aiplan.Move2Dest,
		aiplan.Revenge,
		aiplan.UsePortal,
		aiplan.RechargeSafe,
		aiplan.RechargeCan,
		aiplan.PickupCarryObj,
		aiplan.Equip,
		aiplan.Attack,
		aiplan.MoveStraight3,
		aiplan.MoveStraight5,
	},
	aotype.User: aiplan.PlanList{
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
		aiplan.Attack,
		aiplan.MoveStraight3,
		aiplan.MoveStraight5,
	},
}

var allPlanList = []planObj{
	aiplan.None:           {ai_initPlanNone, ai_actPlanNone},
	aiplan.Chat:           {ai_initPlanChat, ai_actPlanChat},
	aiplan.StrollAround:   {ai_initPlanStrollAround, ai_actPlanStrollAround},
	aiplan.Move2Dest:      {ai_initPlanMove2Dest, ai_actPlanMove2Dest},
	aiplan.Revenge:        {ai_initPlanRevenge, ai_actPlanRevenge},
	aiplan.UsePortal:      {ai_initPlanUsePortal, ai_actPlanUsePortal},
	aiplan.MoveToRecycler: {ai_initPlanMoveToRecycler, ai_actPlanMoveToRecycler},
	aiplan.RechargeSafe:   {ai_initPlanRechargeSafe, ai_actPlanRechargeSafe},
	aiplan.RechargeCan:    {ai_initPlanRechargeCan, ai_actPlanRechargeCan},
	aiplan.PickupCarryObj: {ai_initPlanPickupCarryObj, ai_actPlanPickupCarryObj},
	aiplan.Equip:          {ai_initPlanEquip, ai_actPlanEquip},
	aiplan.UsePotion:      {ai_initPlanUsePotion, ai_actPlanUsePotion},
	aiplan.Attack:         {ai_initPlanAttack, ai_actPlanAttack},
	aiplan.MoveStraight3:  {ai_initPlanMoveStraight3, ai_actPlanMoveStraight3},
	aiplan.MoveStraight5:  {ai_initPlanMoveStraight5, ai_actPlanMoveStraight5},
}
