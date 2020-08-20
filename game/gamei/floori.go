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
	"context"
	"net/http"

	"github.com/kasworld/actpersec"
	"github.com/kasworld/goguelike/game/bias"
	"github.com/kasworld/goguelike/game/fieldobject"
	"github.com/kasworld/goguelike/game/terraini"
	"github.com/kasworld/goguelike/lib/uuidposman"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
)

type FloorI interface {
	Initialized() bool
	Cleanup()
	GetActTurn() int
	GetName() string
	GetWidth() int
	GetHeight() int
	VisitableCount() int // for visitarea

	GetTower() TowerI

	GetBias() bias.Bias
	GetEnvBias() bias.Bias

	GetTerrain() terraini.TerrainI

	GetStatPacketObjOver() *actpersec.ActPerSec
	GetCmdFloorActStat() *actpersec.ActPerSec

	GetActiveObjPosMan() *uuidposman.UUIDPosMan
	GetCarryObjPosMan() *uuidposman.UUIDPosMan
	GetFieldObjPosMan() *uuidposman.UUIDPosMan

	GetReqCh() chan<- interface{}
	Run(ctx context.Context)

	TotalActiveObjCount() int
	TotalCarryObjCount() int
	SearchRandomActiveObjPos() (int, int, error)
	SearchRandomActiveObjPosInRoomOrRandPos() (int, int, error)

	FindPath(dstx, dsty, srcx, srcy int, limit int) [][2]int

	Web_FloorInfo(w http.ResponseWriter, r *http.Request)
	Web_FloorImageZoom(w http.ResponseWriter, r *http.Request)
	Web_FloorImageAutoZoom(w http.ResponseWriter, r *http.Request)
	Web_TileInfo(w http.ResponseWriter, r *http.Request)

	ToPacket_FloorInfo() *c2t_obj.FloorInfo

	FindUsablePortalPairAt(x, y int) (*fieldobject.FieldObject, *fieldobject.FieldObject, error)
}
