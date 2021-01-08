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

package t2g_obj

import (
	"github.com/kasworld/goguelike/config/towerconfig"
	"github.com/kasworld/goguelike/game/aoscore"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
)

type ReqInvalid_data struct {
	Dummy uint8
}
type RspInvalid_data struct {
	Dummy uint8
}

type ReqRegister_data struct {
	TI *c2t_obj.TowerInfo
	SI *c2t_obj.ServiceInfo
	TC *towerconfig.TowerConfig
}
type RspRegister_data struct {
	Dummy uint8
}

type ReqHeartbeat_data struct {
	TowerUUID  string
	StatusInfo []string
}
type RspHeartbeat_data struct {
	Dummy uint8
}

type ReqHighScore_data struct {
	ActiveObjScore *aoscore.ActiveObjScore
}
type RspHighScore_data struct {
	Dummy uint8
}
