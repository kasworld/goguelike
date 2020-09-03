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
	"github.com/kasworld/goguelike/config/gameconst"
	"github.com/kasworld/goguelike/enum/achievetype"
	"github.com/kasworld/goguelike/enum/turnresulttype"
	"github.com/kasworld/goguelike/game/activeobject/turnresult"
	"github.com/kasworld/goguelike/game/fieldobject"
	"github.com/kasworld/goguelike/game/gamei"
)

// follow floor env bias
func (ao *ActiveObject) Noti_Death(f gamei.FloorI) {
	envBias := f.GetEnvBias()
	newActiveObjBiasLen := float64(gameconst.ActiveObjBaseBiasLen)
	newBias := ao.currentBias.Add(envBias.Idiv(10)).MakeAbsSumTo(newActiveObjBiasLen)
	ao.currentBias = newBias
	ao.battleExp *= gameconst.ActiveObjExp_DieRate
	ao.remainTurn2Rebirth += gameconst.ActiveObjRebirthWaitTurn
	ao.remainTurn2Act = 0
	ao.achieveStat.Inc(achievetype.Death)
	if ao.ai != nil {
		ao.ai.ResetPlan()
	}
}

// Kill other ao, inc exp
func (ao *ActiveObject) Kill(dst gamei.ActiveObjectI) {
	ao.AddBattleExp(dst.GetTurnData().Level * gameconst.ActiveObjExp_KillLevel)
	ao.achieveStat.Inc(achievetype.Kill)
	ao.AppendTurnResult(turnresult.New(turnresulttype.Kill, dst, 0))
	dst.AppendTurnResult(turnresult.New(turnresulttype.KilledBy, ao, 0))
}

func (ao *ActiveObject) ApplyHPSPDecByActOnTile(hp, sp float64) {
	reducedSP := ao.ReduceSP(sp)
	if reducedSP != sp {
		hp += sp - reducedSP
	}
	ao.ReduceHP(hp)
	ao.AppendTurnResult(turnresult.New(turnresulttype.DamagedByTile, nil, hp))
	if !ao.IsAlive() {
		ao.AppendTurnResult(turnresult.New(turnresulttype.DeadByTile, nil, 0))
	}
}

func (ao *ActiveObject) ApplyDamageFromDangerObj() bool {
	before := ao.IsAlive()
	for _, v := range ao.turnResultList {
		if v.GetTurnResultType() == turnresulttype.AttackedFrom {
			ao.ReduceHP(v.GetDamage())
			if before && !ao.IsAlive() { // just killed
				dstObj := v.GetDstObj()
				switch o := dstObj.(type) {
				default:
					ao.log.Fatal("unknown dstao %v", v)
				case *ActiveObject:
					o.Kill(ao)
				case *fieldobject.FieldObject:
				}
			}
		}
	}
	return before && !ao.IsAlive()
}
