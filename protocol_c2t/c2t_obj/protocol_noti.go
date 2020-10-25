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

	"github.com/kasworld/goguelike/config/viewportdata"
	"github.com/kasworld/goguelike/enum/fieldobjacttype"
	"github.com/kasworld/goguelike/game/tilearea"
)

type NotiInvalid_data struct {
	Dummy uint8
}

type NotiEnterTower_data struct {
	TowerInfo *TowerInfo
}
type NotiLeaveTower_data struct {
	TowerInfo *TowerInfo
}

type NotiEnterFloor_data struct {
	FI *FloorInfo
}
type NotiLeaveFloor_data struct {
	FI *FloorInfo
}

type NotiAgeing_data struct {
	FloorName string
}

type NotiDeath_data struct {
	Dummy uint8
}

type NotiReadyToRebirth_data struct {
	Dummy uint8
}
type NotiRebirthed_data struct {
	Dummy uint8
}

type NotiBroadcast_data struct {
	Msg string
}

type NotiObjectList_data struct {
	Time          time.Time `prettystring:"simple"`
	FloorName     string
	ActiveObj     *PlayerActiveObjInfo
	ActiveObjList []*ActiveObjClient
	CarryObjList  []*CarryObjClientOnFloor
	FieldObjList  []*FieldObjClient
	DangerObjList []*DangerObjClient
}

// NotiVPTiles_data contains tile info center from pos
type NotiVPTiles_data struct {
	FloorName string
	VPX       int // viewport center X
	VPY       int // viewport center Y
	VPTiles   *viewportdata.ViewportTileArea2
}

// NotiFloorTiles_data used for floor map, reconnect client
type NotiFloorTiles_data struct {
	FI    *FloorInfo
	X     int // X start position, not center
	Y     int // Y start position, not center
	Tiles tilearea.TileArea
}

type NotiFoundFieldObj_data struct {
	FloorName string
	FieldObj  *FieldObjClient
}

type NotiForgetFloor_data struct {
	FloorName string
}

type NotiActivateTrap_data struct {
	FieldObjAct fieldobjacttype.FieldObjActType
	Triggered   bool
}
