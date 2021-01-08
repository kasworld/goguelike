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

package cmd2floor

import (
	"github.com/kasworld/goguelike/game/gamei"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_error"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
)

type ReqRebirth2Floor struct {
	ActiveObj gamei.ActiveObjectI
	X         int
	Y         int
}
type ReqEnterFloor struct {
	ActiveObj gamei.ActiveObjectI
	X         int
	Y         int
}
type ReqLeaveFloor struct {
	ActiveObj gamei.ActiveObjectI
}

type APIAdminTeleport2Floor struct {
	ActiveObj gamei.ActiveObjectI
	ReqPk     *c2t_obj.ReqAdminTeleport_data
	RspCh     chan<- c2t_error.ErrorCode
}

type APIAdminCmd2Floor struct {
	ActiveObj gamei.ActiveObjectI
	ReqPk     *c2t_obj.ReqAdminFloorCmd_data
	RspCh     chan<- c2t_error.ErrorCode
}
