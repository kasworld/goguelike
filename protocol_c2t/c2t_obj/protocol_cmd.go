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

	"github.com/kasworld/goguelike/enum/achievetype_vector"
	"github.com/kasworld/goguelike/enum/condition_vector"
	"github.com/kasworld/goguelike/enum/fieldobjacttype_vector"
	"github.com/kasworld/goguelike/enum/potiontype_vector"
	"github.com/kasworld/goguelike/enum/scrolltype_vector"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_idcmd_stats"
)

type ReqInvalid_data struct {
	Dummy uint8
}
type RspInvalid_data struct {
	Dummy uint8
}

type ReqLogin_data struct {
	SessionUUID string
	NickName    string
	AuthKey     string
}
type RspLogin_data struct {
	ServiceInfo *ServiceInfo
	AccountInfo *AccountInfo
}

type ReqHeartbeat_data struct {
	Time time.Time `prettystring:"simple"`
}
type RspHeartbeat_data struct {
	Time time.Time `prettystring:"simple"`
}

type ReqChat_data struct {
	Chat string
}
type RspChat_data struct {
	Dummy uint8
}

type ReqAchieveInfo_data struct {
	Dummy uint8
}
type RspAchieveInfo_data struct {
	AchieveStat   achievetype_vector.AchieveTypeVector         `prettystring:"simple"`
	PotionStat    potiontype_vector.PotionTypeVector           `prettystring:"simple"`
	ScrollStat    scrolltype_vector.ScrollTypeVector           `prettystring:"simple"`
	FOActStat     fieldobjacttype_vector.FieldObjActTypeVector `prettystring:"simple"`
	AOActionStat  c2t_idcmd_stats.CommandIDStat                `prettystring:"simple"`
	ConditionStat condition_vector.ConditionVector             `prettystring:"simple"`
}

type ReqRebirth_data struct {
	Dummy uint8
}
type RspRebirth_data struct {
	Dummy uint8
}

type ReqMoveFloor_data struct {
	UUID string
}
type RspMoveFloor_data struct {
	Dummy uint8
}

type ReqAIPlay_data struct {
	On bool
}
type RspAIPlay_data struct {
	Dummy uint8
}

// VisitFloorList floor info of visited
type ReqVisitFloorList_data struct {
	Dummy uint8 // change as you need
}

// VisitFloorList floor info of visited
type RspVisitFloorList_data struct {
	FloorList []*VisitFloorInfo
}
