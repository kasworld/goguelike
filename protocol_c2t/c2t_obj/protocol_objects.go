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

package c2t_obj

import (
	"time"

	"github.com/kasworld/goguelike/enum/aiplan"
	"github.com/kasworld/goguelike/enum/dangertype"
	"github.com/kasworld/goguelike/enum/way9type"

	"github.com/kasworld/goguelike/enum/carryingobjecttype"
	"github.com/kasworld/goguelike/enum/condition_flag"
	"github.com/kasworld/goguelike/enum/equipslottype"
	"github.com/kasworld/goguelike/enum/factiontype"
	"github.com/kasworld/goguelike/enum/fieldobjacttype"
	"github.com/kasworld/goguelike/enum/fieldobjdisplaytype"
	"github.com/kasworld/goguelike/enum/potiontype"
	"github.com/kasworld/goguelike/enum/scrolltype"
	"github.com/kasworld/goguelike/enum/turnresulttype"
	"github.com/kasworld/goguelike/game/aoactreqrsp"
	"github.com/kasworld/goguelike/game/bias"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_error"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_idcmd"
	"github.com/kasworld/prettystring"
)

type ServiceInfo struct {
	Version         string
	ProtocolVersion string
	DataVersion     string
}

func (info *ServiceInfo) StringForm() string {
	return prettystring.PrettyString(info, 4)
}

type AccountInfo struct {
	SessionUUID   string
	ActiveObjUUID string
	NickName      string
	CmdList       [c2t_idcmd.CommandID_Count]bool
}

type TowerInfo struct {
	UUID          string
	Name          string
	Factor        [3]int64 `prettystring:"simple"`
	TotalFloorNum int
	StartTime     time.Time `prettystring:"simple"`
	TurnPerSec    float64
}

func (info *TowerInfo) StringForm() string {
	return prettystring.PrettyString(info, 4)
}

type FloorInfo struct {
	Name       string
	W          int
	H          int
	Tiles      int
	Bias       bias.Bias
	TurnPerSec float64
}

func (fi FloorInfo) GetName() string {
	return fi.Name
}
func (fi FloorInfo) GetWidth() int {
	return fi.W
}
func (fi FloorInfo) GetHeight() int {
	return fi.H
}
func (fi FloorInfo) VisitableCount() int {
	return fi.Tiles
}

type CarryObjClientOnFloor struct {
	UUID               string
	X                  int
	Y                  int
	CarryingObjectType carryingobjecttype.CarryingObjectType

	// for equip
	EquipType equipslottype.EquipSlotType
	Faction   factiontype.FactionType

	// for potion
	PotionType potiontype.PotionType

	// for scroll
	ScrollType scrolltype.ScrollType

	// for money
	Value int
}

type FieldObjClient struct {
	ID          string
	X           int
	Y           int
	ActType     fieldobjacttype.FieldObjActType
	DisplayType fieldobjdisplaytype.FieldObjDisplayType
	Message     string
}

func (p *FieldObjClient) GetUUID() string {
	return p.ID
}

type DangerObjClient struct {
	UUID       string
	OwnerID    string
	DangerType dangertype.DangerType
	X          int
	Y          int
	AffectRate float64
}

type ActiveObjClient struct {
	UUID       string
	NickName   string
	Faction    factiontype.FactionType
	EquippedPo []*EquipClient
	Conditions condition_flag.ConditionFlag // not all condition
	X          int
	Y          int
	Alive      bool
	Chat       string

	// turn result
	Act        c2t_idcmd.CommandID
	Dir        way9type.Way9Type
	Result     c2t_error.ErrorCode
	DamageGive int
	DamageTake int
}

type PlayerActiveObjInfo struct {
	Bias       bias.Bias
	Conditions condition_flag.ConditionFlag
	Exp        int
	TotalAO    int // total activeobject count
	Ranking    int // ranking / totalao
	Death      int
	Kill       int
	Sight      float64
	HP         int
	HPMax      int
	SP         int
	SPMax      int
	AIPlan     aiplan.AIPlan
	EquippedPo []*EquipClient
	EquipBag   []*EquipClient
	PotionBag  []*PotionClient
	ScrollBag  []*ScrollClient
	Wallet     int
	Wealth     int
	ActiveBuff []*ActiveObjBuff
	AP         float64

	Act        *aoactreqrsp.ActReqRsp
	TurnResult []TurnResultClient
}

type TurnResultClient struct {
	ResultType turnresulttype.TurnResultType
	DstUUID    string
	Arg        float64
}

type EquipClient struct {
	UUID      string
	Name      string
	EquipType equipslottype.EquipSlotType
	Faction   factiontype.FactionType
	BiasLen   float64
}

type PotionClient struct {
	UUID       string
	PotionType potiontype.PotionType
}
type ScrollClient struct {
	UUID       string
	ScrollType scrolltype.ScrollType
}

type ActiveObjBuff struct {
	Name        string
	RemainCount int
}
