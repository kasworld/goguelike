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

package cmd2tower

import (
	"fmt"

	"github.com/kasworld/goguelike/game/fieldobject"
	"github.com/kasworld/goguelike/game/gamei"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_error"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
)

type ActiveObjEnterTower struct {
	ActiveObj gamei.ActiveObjectI
	RspCh     chan<- error
}

func (cet ActiveObjEnterTower) String() string {
	return fmt.Sprintf("ActiveObjEnterTower[%v]",
		cet.ActiveObj,
	)
}

type ActiveObjSuspendFromTower struct {
	ActiveObj gamei.ActiveObjectI
	RspCh     chan<- error
}

func (cet ActiveObjSuspendFromTower) String() string {
	return fmt.Sprintf("ActiveObjSuspendFromTower[%v]",
		cet.ActiveObj,
	)
}

type ActiveObjResumeTower struct {
	ActiveObj gamei.ActiveObjectI
	RspCh     chan<- error
}

func (cet ActiveObjResumeTower) String() string {
	return fmt.Sprintf("ActiveObjResumeTower[%v]",
		cet.ActiveObj,
	)
}

type ActiveObjLeaveTower struct {
	ActiveObjUUID string
	RspCh         chan<- error
}

func (cet ActiveObjLeaveTower) String() string {
	return fmt.Sprintf("ActiveObjLeaveTower[%v]",
		cet.ActiveObjUUID,
	)
}

type AdminFloorMove struct {
	ActiveObj  gamei.ActiveObjectI
	RecvPacket *c2t_obj.ReqAdminFloorMove_data
	RspCh      chan<- c2t_error.ErrorCode
}

func (cet AdminFloorMove) String() string {
	return fmt.Sprintf("AdminFloorMove[%v %v]",
		cet.ActiveObj,
		cet.RecvPacket,
	)
}

type FloorMove struct {
	ActiveObj gamei.ActiveObjectI
	FloorUUID string
	RspCh     chan<- c2t_error.ErrorCode
}

func (cet FloorMove) String() string {
	return fmt.Sprintf("FloorMove[%v %v]",
		cet.ActiveObj,
		cet.FloorUUID,
	)
}

type AdminTowerCmd struct {
	ActiveObj  gamei.ActiveObjectI
	RecvPacket *c2t_obj.ReqAdminTowerCmd_data
	RspCh      chan<- c2t_error.ErrorCode
}

func (cet AdminTowerCmd) String() string {
	return fmt.Sprintf("AdminTowerCmd[%v %v]",
		cet.ActiveObj,
		cet.RecvPacket,
	)
}

type ActiveObjTrapTeleport struct {
	ActiveObj    gamei.ActiveObjectI
	SrcFloor     gamei.FloorI
	DstFloorName string
}

type ActiveObjUsePortal struct {
	ActiveObj gamei.ActiveObjectI
	SrcFloor  gamei.FloorI
	P1, P2    *fieldobject.FieldObject
}

func (pk ActiveObjUsePortal) String() string {
	return fmt.Sprintf(
		"ActiveObjUsePortal[%v %v %v %v]",
		pk.SrcFloor,
		pk.ActiveObj,
		pk.P1,
		pk.P2,
	)
}

type ActiveObjRebirth struct {
	ActiveObj gamei.ActiveObjectI
}

func (pk ActiveObjRebirth) String() string {
	return fmt.Sprintf(
		"ActiveObjRebirth[%v]",
		pk.ActiveObj,
	)
}
