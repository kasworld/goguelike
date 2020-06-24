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
	"github.com/kasworld/goguelike/enum/achievetype"
	"github.com/kasworld/goguelike/enum/turnresulttype"
	"github.com/kasworld/goguelike/game/aoexpsort"
	"github.com/kasworld/goguelike/game/aoscore"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
)

func (ao *ActiveObject) ToPacket_ActiveObjClient(x, y int) *c2t_obj.ActiveObjClient {
	if ao.aoClientCache != nil {
		return ao.aoClientCache
	}
	aoc := &c2t_obj.ActiveObjClient{
		UUID:       ao.uuid,
		NickName:   ao.nickName,
		Faction:    ao.currentBias.NearFaction(),
		EquippedPo: ao.inven.ToPacket_EquipClient(),
		X:          x,
		Y:          y,
		Alive:      ao.IsAlive(),
		Chat:       ao.chat,
	}

	if stepAct := ao.turnActReqRsp; stepAct != nil && stepAct.Acted() {
		aoc.Act = stepAct.Done.Act
		aoc.Dir = stepAct.Done.Dir
		aoc.Result = stepAct.Error
	}
	var DamageGive float64
	var DamageTake float64

	for _, v := range ao.turnResultList {
		switch v.GetTurnResultType() {
		case turnresulttype.AttackTo:
			DamageGive += v.GetDamage()
		case turnresulttype.AttackedFrom:
			DamageTake += v.GetDamage()
		case turnresulttype.DamagedByTile:
			DamageTake += v.GetDamage()
		}
	}
	aoc.DamageGive = int(DamageGive)
	aoc.DamageTake = int(DamageTake)
	ao.aoClientCache = aoc
	return aoc
}

func (ao *ActiveObject) ToPacket_PlayerActiveObjInfo() *c2t_obj.PlayerActiveObjInfo {
	aoList := ao.GetHomeFloor().GetTower().GetExpRanking()
	rank := aoexpsort.ByExp(aoList).FindPosByKey(ao.AOTurnData.TotalExp)

	rtn := c2t_obj.PlayerActiveObjInfo{
		Bias:       ao.currentBias,
		Conditions: ao.AOTurnData.Condition,

		Exp:     int(ao.AOTurnData.TotalExp),
		Ranking: rank,
		Death:   int(ao.achieveStat.Get(achievetype.Death)),
		Kill:    int(ao.achieveStat.Get(achievetype.Kill)),
		Sight:   ao.AOTurnData.Sight,

		HP:    int(ao.hp),
		HPMax: int(ao.AOTurnData.HPMax),
		SP:    int(ao.sp),
		SPMax: int(ao.AOTurnData.SPMax),

		AIPlan: ao.ai.GetPlan(),

		Act:            ao.turnActReqRsp,
		RemainTurn2Act: ao.remainTurn2Act,
	}
	rtn.Wealth = int(ao.inven.GetTotalValue())
	rtn.EquippedPo, rtn.EquipBag, rtn.PotionBag, rtn.ScrollBag, rtn.Wallet = ao.inven.ToPacket_InvenInfos()
	rtn.TurnResult = make([]c2t_obj.TurnResultClient,
		0, len(ao.turnResultList))
	for _, v := range ao.turnResultList {
		rtn.TurnResult = append(rtn.TurnResult, v.ToPacket_TurnResultClient())
	}

	rtn.ActiveBuff = ao.buffManager.ToPacket_ActiveObjBuffList()

	return &rtn
}

func (ao *ActiveObject) To_ActiveObjScore() *aoscore.ActiveObjScore {
	aos := &aoscore.ActiveObjScore{
		UUID:        ao.uuid,
		NickName:    ao.nickName,
		AchieveStat: ao.achieveStat,
		ActionStat:  ao.aoActionStat,
		Exp:         ao.AOTurnData.TotalExp,
		Wealth:      ao.inven.GetTotalValue(),
		BornFaction: ao.bornFaction,
		CurrentBias: ao.currentBias,
	}
	return aos
}
