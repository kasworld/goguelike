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

package gamei

import (
	"time"

	"github.com/kasworld/goguelike/config/viewportdata"
	"github.com/kasworld/goguelike/enum/achievetype_vector"
	"github.com/kasworld/goguelike/enum/aotype"
	"github.com/kasworld/goguelike/enum/fieldobjacttype_vector"
	"github.com/kasworld/goguelike/enum/scrolltype_vector"
	"github.com/kasworld/goguelike/game/activeobject/activebuff"
	"github.com/kasworld/goguelike/game/activeobject/aoturndata"
	"github.com/kasworld/goguelike/game/activeobject/turnresult"
	"github.com/kasworld/goguelike/game/aoactreqrsp"
	"github.com/kasworld/goguelike/game/aoscore"
	"github.com/kasworld/goguelike/game/bias"
	"github.com/kasworld/goguelike/game/visitarea"
	"github.com/kasworld/goguelike/lib/scriptparse"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_error"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_serveconnbyte"
)

type ActiveObjectI interface {
	Cleanup()
	GetUUID() string
	// GetNickName() string
	String() string

	GetInven() InventoryI

	GetHP() float64
	// GetSP() float64
	GetHPRate() float64
	GetSPRate() float64

	Suspend()
	Resume(conn *c2t_serveconnbyte.ServeConnByte)

	IsAlive() bool

	Noti_Rebirth()
	Noti_EnterFloor(f FloorI)
	Noti_Death(f FloorI)

	GetBias() bias.Bias
	// GetBornFaction() factiontype.FactionType

	AddBattleExp(v float64)

	// for sort
	GetExpCopy() float64
	UpdateExpCopy()

	NeedCharge(limit float64) bool
	Charged(limit float64) bool
	ReduceSP(apToReduce float64) (apReduced float64)
	ReduceHP(hpToReduce float64) (hpReduced float64)

	GetTurnResultList() []turnresult.TurnResult
	SetTurnActReqRsp(actrsp *aoactreqrsp.ActReqRsp)
	GetTurnActReqRsp() *aoactreqrsp.ActReqRsp
	SetNeedTANoti()
	GetAndClearNeedTANoti() bool
	GetRemainTurn2Act() float64

	SetReq2Handle(req *aoactreqrsp.Act)
	GetClearReq2Handle() *aoactreqrsp.Act

	GetTurnData() *aoturndata.ActiveObjTurnData
	GetBuffManager() *activebuff.BuffManager

	GetClientConn() *c2t_serveconnbyte.ServeConnByte
	GetActiveObjType() aotype.ActiveObjType

	IsAIUse() bool
	SetUseAI(b bool)
	RunAI(turnTime time.Time)

	GetChat() string
	SetChat(c string)

	GetRemainTurn2Rebirth() int
	TryRebirth() error

	GetCurrentFloor() FloorI

	PrepareNewTurn(turnTime time.Time)
	ApplyTurnAct()
	AppendTurnResult(turnResult turnresult.TurnResult)

	ApplyDamageFromActiveObj() bool
	ApplyHPSPDecByActOnTile(hp, sp float64)
	Kill(dst ActiveObjectI)

	DoEquip(poid string) error
	DoUnEquip(poid string) error
	DoUseCarryObj(poid string) error
	DoRecycleCarryObj(poid string) error
	DoAIOnOff(onoff bool) error
	DoPickup(po CarryingObjectI) error

	// for tower
	GetHomeFloor() FloorI

	ToPacket_ActiveObjClient(x, y int) *c2t_obj.ActiveObjClient
	ToPacket_PlayerActiveObjInfo() *c2t_obj.PlayerActiveObjInfo
	To_ActiveObjScore() *aoscore.ActiveObjScore

	GetAchieveStat() *achievetype_vector.AchieveTypeVector
	GetFieldObjActStat() *fieldobjacttype_vector.FieldObjActTypeVector
	// GetPotionStat() *potiontype_vector.PotionTypeVector
	GetScrollStat() *scrolltype_vector.ScrollTypeVector
	// GetActStats() *c2t_idcmd_stats.CommandIDStat

	UpdateBySightMat2(f FloorI, vpCenterX, vpCenterY int,
		sightMat *viewportdata.ViewportSight2, sight float32)

	GetVisitFloor(floorname string) *visitarea.VisitArea
	ForgetFloorByName(floorname string) error
	MakeFloorComplete(f FloorI) error

	DoAdminCmd(cmd, args string) c2t_error.ErrorCode
	DoParsedAdminCmd(ca *scriptparse.CmdArgs) error
}

const (
	ActiveObjectI_HTML_tableheader = `
	<tr>
	<td>ActiveObj</td>
	<td>bias</td>
	<td>level</td>
	<td>exp</td>
	<td>AI</td>
	</tr>	
`
	ActiveObjectI_HTML_row = `
	<tr>
		<td>
		<a href= "/ActiveObj?aoid={{$v.GetUUID}}" >
			{{$v}}
		</a>
		</td>
		<td>
		{{$v.GetBias}} 
		</td>
		<td>
		{{$v.GetTurnData.Level | printf "%4.2f" }} 
		</td>
		<td>
		{{$v.GetTurnData.TotalExp | printf "%4.2f" }} 
		</td>
		<td>
		{{if $v.GetAIObj }}
			{{$v.GetAIObj}}
		{{end}}
		</td>
	</tr>
`
)
