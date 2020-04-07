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

	"github.com/kasworld/goguelike/enum/achievetype"
	"github.com/kasworld/goguelike/lib/g2id"
)

type ReqInvalid_data struct {
	Dummy uint8
}
type RspInvalid_data struct {
	Dummy uint8
}

type ReqLogin_data struct {
	SessionG2ID g2id.G2ID
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
	Achieve [achievetype.AchieveType_Count]int64
}

type ReqRebirth_data struct {
	Dummy uint8
}
type RspRebirth_data struct {
	Dummy uint8
}
