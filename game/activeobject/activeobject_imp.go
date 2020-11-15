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
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/kasworld/goguelike/enum/respawntype"

	"github.com/kasworld/goguelike/enum/achievetype_vector"
	"github.com/kasworld/goguelike/enum/aotype"
	"github.com/kasworld/goguelike/enum/condition_vector"
	"github.com/kasworld/goguelike/enum/fieldobjacttype_vector"
	"github.com/kasworld/goguelike/enum/potiontype_vector"
	"github.com/kasworld/goguelike/enum/scrolltype_vector"
	"github.com/kasworld/goguelike/game/activeobject/activebuff"
	"github.com/kasworld/goguelike/game/activeobject/aoturndata"
	"github.com/kasworld/goguelike/game/activeobject/turnresult"
	"github.com/kasworld/goguelike/game/aoactreqrsp"
	"github.com/kasworld/goguelike/game/bias"
	"github.com/kasworld/goguelike/game/floor4client"
	"github.com/kasworld/goguelike/game/gamei"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_idcmd_stats"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_serveconnbyte"
)

func (ao *ActiveObject) GetUUID() string {
	return ao.uuid
}

func (ao *ActiveObject) GetHomeFloor() gamei.FloorI {
	return ao.homefloor
}

func (ao *ActiveObject) GetCurrentFloor() gamei.FloorI {
	return ao.currentFloor
}

func (ao *ActiveObject) GetFloor4Client(floorname string) *floor4client.Floor4Client {
	r, _ := ao.floor4ClientMan.GetByName(floorname)
	return r
}

func (ao *ActiveObject) GetFloor4ClientList() []*floor4client.Floor4Client {
	return ao.floor4ClientMan.GetList()
}

func (ao *ActiveObject) GetActiveObjType() aotype.ActiveObjType {
	return ao.aoType
}

func (ao *ActiveObject) GetRespawnType() respawntype.RespawnType {
	return ao.respawnType
}

func (ao *ActiveObject) GetChat() string {
	return ao.chat
}
func (ao *ActiveObject) SetChat(c string) {
	ao.chat = c
	ao.chatTime = time.Now()
}

// func (ao *ActiveObject) SetNickName(nickname string) {
// 	ao.nickName = nickname
// }

// func (ao *ActiveObject) GetNickName() string {
// 	return ao.nickName
// }

func (ao *ActiveObject) SetNeedTANoti() {
	ao.needTANoti = true
}
func (ao *ActiveObject) GetAndClearNeedTANoti() bool {
	rtn := ao.needTANoti
	ao.needTANoti = false
	return rtn
}

// clients conn interface
func (ao *ActiveObject) GetClientConn() *c2t_serveconnbyte.ServeConnByte {
	return ao.clientConn
}

// func (ao *ActiveObject) SetActiveObjType(aot aotype.ActiveObjType) {
// 	ao.aoType = aot
// }

func (ao *ActiveObject) GetBias() bias.Bias {
	return ao.currentBias
}

func (ao *ActiveObject) GetInven() gamei.InventoryI {
	return ao.inven
}

////////////////////////////////////////////////////////////////////////////////
// stats, at least used for web

func (ao *ActiveObject) GetAchieveStat() *achievetype_vector.AchieveTypeVector {
	return &ao.achieveStat
}

func (ao *ActiveObject) GetScrollStat() *scrolltype_vector.ScrollTypeVector {
	return &ao.scrollStat
}

func (ao *ActiveObject) GetFieldObjActStat() *fieldobjacttype_vector.FieldObjActTypeVector {
	return &ao.foActStat
}

func (ao *ActiveObject) GetConditionStat() *condition_vector.ConditionVector {
	return &ao.conditionStat
}

func (ao *ActiveObject) GetActStat() *c2t_idcmd_stats.CommandIDStat {
	return &ao.aoActionStat
}

func (ao *ActiveObject) GetPotionStat() *potiontype_vector.PotionTypeVector {
	return &ao.potionStat
}

////////////////////////////////////////////////////////////////////////////////
// battle relate

// aoactreqrsp.Act
func (ao *ActiveObject) SetReq2Handle(req *aoactreqrsp.Act) {
	atomic.StorePointer(&ao.req2Handle, unsafe.Pointer(req))
}

// aoactreqrsp.Act
func (ao *ActiveObject) GetClearReq2Handle() *aoactreqrsp.Act {
	r := atomic.SwapPointer(&ao.req2Handle, nil)
	return (*aoactreqrsp.Act)(r)
}

func (ao *ActiveObject) GetTurnActReqRsp() *aoactreqrsp.ActReqRsp {
	return ao.turnActReqRsp
}

func (ao *ActiveObject) GetRemainTurn2Rebirth() int {
	return ao.remainTurn2Rebirth
}

func (ao *ActiveObject) GetAP() float64 {
	return ao.ap
}

func (ao *ActiveObject) GetTurnResultList() []turnresult.TurnResult {
	return ao.turnResultList
}

func (ao *ActiveObject) GetHP() float64 {
	return ao.hp
}

func (ao *ActiveObject) GetSPRate() float64 {
	return ao.sp / ao.AOTurnData.SPMax
}

func (ao *ActiveObject) GetHPRate() float64 {
	return ao.hp / ao.AOTurnData.HPMax
}

func (ao *ActiveObject) IsAlive() bool {
	return ao.hp > 0
}

func (ao *ActiveObject) NeedCharge(limit float64) bool {
	return ao.GetSPRate() < limit || ao.GetHPRate() < limit
}
func (ao *ActiveObject) Charged(limit float64) bool {
	return !ao.NeedCharge(limit)
}

func (ao *ActiveObject) ReduceHP(hpToReduce float64) float64 {
	oldvalue := ao.hp
	ao.hp -= hpToReduce
	return oldvalue - ao.hp
}
func (ao *ActiveObject) ReduceSP(apToReduce float64) float64 {
	oldvalue := ao.sp
	ao.sp -= apToReduce
	if ao.sp < 0 {
		ao.sp = 0
	}
	return oldvalue - ao.sp
}

func (ao *ActiveObject) GetTurnData() *aoturndata.ActiveObjTurnData {
	return ao.AOTurnData
}

func (ao *ActiveObject) AddBattleExp(v float64) {
	ao.battleExp += v
}

func (ao *ActiveObject) GetBuffManager() *activebuff.BuffManager {
	return ao.buffManager
}
